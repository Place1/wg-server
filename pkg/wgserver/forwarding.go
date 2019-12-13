package wgserver

import (
	"github.com/place1/wg-embed/pkg/wgembed"
)

func ConfigureForwarding(iface *wgembed.WireGuardInterface, gatewayIface string) error {
	return forwarding(iface, gatewayIface)
}
