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

package openvpn

import (
	"github.com/thanhps42/go-openvpn/openvpn/config"
	"github.com/thanhps42/go-openvpn/openvpn/management"
	"github.com/thanhps42/go-openvpn/openvpn/tunnel"
	"os/exec"
	"syscall"
)

// CreateNewProcess creates new openvpn process with given config params
func CreateNewProcess(openvpnBinary string, config *config.GenericConfig, middlewares ...management.Middleware) *OpenvpnProcess {
	tunnelSetup := tunnel.NewTunnelSetup()
	execCommand := func(arg ...string) *exec.Cmd {
		cmd := exec.Command(openvpnBinary, arg...)
		cmd.SysProcAttr = &syscall.SysProcAttr{
			HideWindow: true,
		}
		return cmd
	}
	return newProcess(tunnelSetup, config, execCommand, middlewares...)
}
