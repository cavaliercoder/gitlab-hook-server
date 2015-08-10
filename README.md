# GitLab Hook Server

This self-contained web server listens for web hook `POST` requests from an
on-premise [GitLab](https://about.gitlab.com/) server. Hooks can be filtered
and configured with custom actions to be performed if a filter is satisfied.

Some examples s'il vous plaît:

 * Execute tests when a merge request is created
 * Deploy an application when a tag is created
 * Deploy Puppet modules with `r10k` when a branch is updated
 * Send an email whenever an issue is created

The project is a __work-in-progress__ and remains unreleased.

## License

GitLab Hook Server (C) 2015  Ryan Armstrong <ryan@cavaliercoder.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see http://www.gnu.org/licenses/.
