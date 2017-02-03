// Copyright (c) - Damien Fontaine <damien.fontaine@lineolia.net>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>

package security

//Application to OAuth2
type Application struct {
	ClientID    string
	Secret      string
	Name        string
	URL         string
	Description string
	Callback    string
	User        User
}

// ApplicationManager to manage Application
type ApplicationManager interface {
	GetByClientID(clientID string) (application Application, err error)
}
