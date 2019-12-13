// +build linux

package wgserver

import (
	"github.com/coreos/go-iptables/iptables"
	"github.com/pkg/errors"
	"github.com/place1/wg-embed/pkg/wgembed"
)

func forwarding(iface *wgembed.WireGuardInterface, gatewayIface string) error {
	ipt, err := iptables.New()
	if err != nil {
		return errors.Wrap(err, "failed to init iptables")
	}
	if err := ipt.AppendUnique("filter", "FORWARD", "-i", gatewayIface, "-o", iface.Name(), "-j", "ACCEPT"); err != nil {
		return errors.Wrap(err, "failed to set ip tables rule")
	}
	if err := ipt.AppendUnique("filter", "FORWARD", "-i", iface.Name(), "-o", gatewayIface, "-j", "ACCEPT"); err != nil {
		return errors.Wrap(err, "failed to set ip tables rule")
	}
	if err := ipt.AppendUnique("nat", "POSTROUTING", "-o", gatewayIface, "-j", "MASQUERADE"); err != nil {
		return errors.Wrap(err, "failed to set ip tables rule")
	}
	return nil
}
