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

type Config struct {
	Path  string `json:"-"`
	Rules []Rule `json:"rules"`
}

func NewConfig() *Config {
	return &Config{}
}

func LoadConfig(path string) (*Config, error) {
	printf("loading configuration file: %s\n", path)

	// open file
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// parse
	config, err := ParseConfig(f)
	if err != nil {
		return nil, err
	}

	config.Path = path

	printf("loaded %d hook rules\n", len(config.Rules))

	return config, nil
}

func ParseConfig(r io.Reader) (*Config, error) {
	config := NewConfig()

	decoder := json.NewDecoder(r)
	err := decoder.Decode(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
