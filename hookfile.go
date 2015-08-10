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
	"encoding/json"
	"io"
	"os"
)

type Hookfile struct {
	Path  string `json:"-"`
	Hooks []Hook `json:"hooks"`
}

func NewHookfile() *Hookfile {
	return &Hookfile{}
}

func LoadHookfile(path string) (*Hookfile, error) {
	printf("loading configuration file: %s\n", path)

	// open file
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// parse
	hookfile, err := ParseHookfile(f)
	if err != nil {
		return nil, err
	}

	hookfile.Path = path

	printf("loaded %d hooks\n", len(hookfile.Hooks))

	return hookfile, nil
}

func ParseHookfile(r io.Reader) (*Hookfile, error) {
	hookfile := NewHookfile()

	decoder := json.NewDecoder(r)
	err := decoder.Decode(hookfile)
	if err != nil {
		return nil, err
	}

	return hookfile, nil
}
