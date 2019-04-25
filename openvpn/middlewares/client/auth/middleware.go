/*
 * Copyright (C) 2018 The "MysteriumNetwork/go-openvpn" Authors.
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

package auth

import (
	"regexp"

	log "github.com/cihub/seelog"
	"github.com/thanhps42/go-openvpn/openvpn"
	"github.com/thanhps42/go-openvpn/openvpn/management"
)

// CredentialsProvider returns client's current auth primitives (i.e. customer identity signature / node's sessionId)
type CredentialsProvider func() (username string, password string, err error)

type middleware struct {
	fetchCredentials CredentialsProvider
	commandWriter    management.CommandWriter
	lastUsername     string
	lastPassword     string
	state            openvpn.State
}

var rule = regexp.MustCompile("^>PASSWORD:Need 'Auth' username/password$")

// NewMiddleware creates client user_auth challenge authentication middleware
func NewMiddleware(credentials CredentialsProvider) *middleware {
	return &middleware{
		fetchCredentials: credentials,
		commandWriter:    nil,
	}
}

func (m *middleware) Start(commandWriter management.CommandWriter) error {
	m.commandWriter = commandWriter
	log.Info("starting client user-pass provider middleware")
	return nil
}

func (m *middleware) Stop(connection management.CommandWriter) error {
	return nil
}

func (m *middleware) ConsumeLine(line string) (consumed bool, err error) {
	match := rule.FindStringSubmatch(line)
	if len(match) == 0 {
		return false, nil
	}

	username, password, err := m.fetchCredentials()
	if err != nil {
		return false, err
	}

	log.Info("authenticating user ", username)

	_, err = m.commandWriter.SingleLineCommand("password 'Auth' %s", password)
	if err != nil {
		return true, err
	}

	_, err = m.commandWriter.SingleLineCommand("username 'Auth' %s", username)
	if err != nil {
		return true, err
	}
	return true, nil
}
