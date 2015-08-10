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
	"fmt"
	"os"
)

func printf(format string, a ...interface{}) {
	fmt.Printf(format, a...)
}

func eprintf(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, "error: %s", fmt.Sprintf(format, a...))
}

func dprintf(format string, a ...interface{}) {
	if Debug {
		fmt.Fprintf(os.Stdout, "debug: %s", fmt.Sprintf(format, a...))
	}
}
