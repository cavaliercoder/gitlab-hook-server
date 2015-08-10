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
// See: https://gitlab.com/gitlab-org/gitlab-ce/blob/master/doc/web_hooks/web_hooks.md
package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
)

const EVENT_HEADER = "X-Gitlab-Event"

const (
	PUSH_EVENT  = "Push Hook"
	TAG_EVENT   = "Tag Push Hook"
	ISSUE_EVENT = "Issue Hook"
)

var (
	Debug           = true
	ListenInterface = ""
	ListenPort      = 8000
	HookfilePath    = ""
)

var hooks []Hook

func main() {
	// parse command line args
	flag.StringVar(&HookfilePath, "file", "", "load hookfile")
	flag.BoolVar(&Debug, "debug", false, "print program debug messages")
	flag.StringVar(&ListenInterface, "interface", "0.0.0.0", "listen interface")
	flag.IntVar(&ListenPort, "port", 8000, "listen TCP port")
	flag.Parse()

	if HookfilePath == "" {
		eprintf("No configuration file specified\n")
		os.Exit(1)
	}

	// load config
	hookfile, err := LoadHookfile(HookfilePath)
	if err != nil {
		panic(err)
	}

	hooks = hookfile.Hooks

	// configure routes
	http.HandleFunc("/", HandleHookRequest)

	// listen
	addr := fmt.Sprintf("%s:%d", ListenInterface, ListenPort)
	printf("starting server on %s\n", addr)
	err = http.ListenAndServe(addr, nil)
	if err != nil {
		panic(err)
	}
}

func HandleHookRequest(w http.ResponseWriter, r *http.Request) {
	if r.Body != nil {
		defer r.Body.Close()
	}

	// decode request body
	hookRequest, err := NewHookRequest(r)
	if err != nil {
		eprintf("%s\n", err)
		return
	}

	printf("new request: %s %s\n", hookRequest.RequestID, hookRequest)

	for _, hook := range hooks {
		if hook.Eval(hookRequest) {
			go func(hook *Hook) {
				if err := hook.Exec(hookRequest); err != nil {
					eprintf("action failed with: %s\n", err.Error())
				}
			}(&hook)
		}
	}
}
