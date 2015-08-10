/*
 * GitLab Hook Server (C) 2015  Ryan Armstrong <ryan@cavaliercoder.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */
package main

import (
	"bufio"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type Hook struct {
	Name     string    `json:"name"`
	Action   []string  `json:"action"`
	Triggers []Trigger `json:"triggers"`
}

type Trigger struct {
	Event      string `json:"event"`
	Repository string `json:"repository"`
	User       string `json:"user"`
}

var macroPattern = regexp.MustCompile(`\$\{.*?\}`)

// Eval processes a hook request, evaluates triggers and executes actions
// if any triggers match.
func (c *Hook) Eval(request *HookRequest) bool {
	// check triggers
	doAction := false
	for _, trigger := range c.Triggers {
		// check if trigger is for this request kind
		if trigger.Event == request.ObjectKind {
			// check if repository matches
			if trigger.Repository != "" {
				switch trigger.Repository {
				case
					strconv.Itoa(request.ProjectID),
					request.Repository.Name,
					request.Repository.URL,
					request.Repository.GitHTTPURL,
					request.Repository.GitSSHURL:

					doAction = true
					break
				}
			}

			if trigger.User != "" {
				switch trigger.User {
				case
					strconv.Itoa(request.UserID),
					request.UserName,
					request.UserEmail:

					doAction = true
					break
				}
			}
		}

		if doAction {
			break
		}
	}

	return doAction
}

// ExpandAction expands macros in a Hook Actions definition with values from
// the Hook Request.
func (c *Hook) ExpandAction(request *HookRequest) []string {
	res := make([]string, len(c.Action))

	// parse each arg in the action
	for i, arg := range c.Action {
		// find macros
		macros := macroPattern.FindAllString(arg, -1)
		for _, macro := range macros {
			val := ""

			switch macro {
			case "${REF}":
				val = request.Ref
			case "${REFNAME}":
				val = filepath.Base(request.Ref)
			}

			// replace macro with value or blank
			arg = strings.Replace(arg, macro, val, -1)
		}

		// append to result
		res[i] = arg
	}

	dprintf("debug: expanded %s to %s\n", c.Action, res)
	return res
}

// Exec executes a Hook Action, and captures the stdout and stderr.
func (c *Hook) Exec(request *HookRequest) error {
	printf("hook: triggered: %s\n", c.Name)

	action := c.ExpandAction(request)

	path := action[0]
	args := action[1:]

	cmd := exec.Command(path, args...)

	// parse stdout async
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			printf("action: ==> %s\n", scanner.Text())
		}
	}()

	// attach to stderr
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			printf("action: ==> %s\n", scanner.Text())
		}
	}()

	// execute
	printf("action: executing: %s %s\n", path, strings.Join(args, " "))
	err = cmd.Start()
	if err != nil {
		return err
	}
	printf("action: started with pid: %d\n", cmd.Process.Pid)

	// wait for process to finish
	err = cmd.Wait()
	if err != nil {
		return err
	}
	printf("action: finished\n")

	return nil
}
