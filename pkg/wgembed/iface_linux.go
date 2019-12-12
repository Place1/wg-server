// +build linux

package wgembed

import (
	"fmt"

	"github.com/coreos/go-iptables/iptables"
	"github.com/pkg/errors"
	"github.com/vishvananda/netlink"
)

func (wg *WireGuardInterface) Up() error {
	link, err := netlink.LinkByName(wg.Name())
	if err != nil {
		return errors.Wrap(err, "failed to find wireguard interface")
	}

	if err := netlink.LinkSetUp(link); err != nil {
		return errors.Wrap(err, "failed to bring wireguard interface up")
	}

	return nil
}

func (wg *WireGuardInterface) SetIP(ip string) error {
	link, err := netlink.LinkByName(wg.Name())
	if err != nil {
		return errors.Wrap(err, "failed to find wireguard interface")
	}

	linkaddr, err := netlink.ParseAddr(fmt.Sprintf("%s/32", ip))
	if err != nil {
		return errors.Wrap(err, "failed to parse wireguard interface ip address")
	}

	if err := netlink.AddrAdd(link, linkaddr); err != nil {
		return errors.Wrap(err, "failed to set ip address of wireguard interface")
	}

	return nil
}

func (wg *WireGuardInterface) ConfigureForwarding(gatewayIface string) error {
	ipt, err := iptables.New()
	if err != nil {
		return errors.Wrap(err, "failed to init iptables")
	}
	if err := ipt.AppendUnique("filter", "FORWARD", "-i", gatewayIface, "-o", wg.Name(), "-j", "ACCEPT"); err != nil {
		return errors.Wrap(err, "failed to set ip tables rule")
	}
	if err := ipt.AppendUnique("filter", "FORWARD", "-i", wg.Name(), "-o", gatewayIface, "-j", "ACCEPT"); err != nil {
		return errors.Wrap(err, "failed to set ip tables rule")
	}
	if err := ipt.AppendUnique("nat", "POSTROUTING", "-o", gatewayIface, "-j", "MASQUERADE"); err != nil {
		return errors.Wrap(err, "failed to set ip tables rule")
	}
	return nil
}
