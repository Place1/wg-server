// +build !windows

/* SPDX-License-Identifier: MIT
 *
 * Copyright (C) 2017-2019 WireGuard LLC. All Rights Reserved.
 */

// modified from https://git.zx2c4.com/wireguard-go

package wgembed

import (
	"net"

	"github.com/pkg/errors"
	"golang.zx2c4.com/wireguard/device"
	"golang.zx2c4.com/wireguard/ipc"
	"golang.zx2c4.com/wireguard/tun"
)

type WireGuardInterface struct {
	device *device.Device
	uapi   net.Listener
}

func StartInterface(interfaceName string) (*WireGuardInterface, error) {
	wg := &WireGuardInterface{}

	tun, err := tun.CreateTUN(interfaceName, device.DefaultMTU)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create TUN device")
	}

	// open UAPI file (or use supplied fd)
	fileUAPI, err := ipc.UAPIOpen(interfaceName)
	if err != nil {
		return nil, errors.Wrap(err, "UAPI listen error")
	}

	wg.device = device.NewDevice(tun, device.NewLogger(device.LogLevelError, interfaceName))

	errs := make(chan error)

	uapi, err := ipc.UAPIListen(interfaceName, fileUAPI)
	if err != nil {
		return nil, errors.Wrap(err, "failed to listen on uapi socket")
	}
	wg.uapi = uapi

	go func() {
		for {
			conn, err := uapi.Accept()
			if err != nil {
				errs <- err
				return
			}
			go wg.device.IpcHandle(conn)
		}
	}()

	return wg, nil
}

func (wg *WireGuardInterface) Close() {
	wg.uapi.Close()
	wg.device.Close()
}
