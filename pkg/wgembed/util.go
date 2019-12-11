package wgembed

import (
	"github.com/coreos/go-iptables/iptables"
	"github.com/pkg/errors"
	"github.com/vishvananda/netlink"
)

func Up(name string) error {
	link, err := netlink.LinkByName(name)
	if err != nil {
		return errors.Wrap(err, "failed to find wireguard interface")
	}

	if err := netlink.LinkSetUp(link); err != nil {
		return errors.Wrap(err, "failed to bring wireguard interface up")
	}

	return nil
}

func SetIP(name string, cidr string) error {
	link, err := netlink.LinkByName(name)
	if err != nil {
		return errors.Wrap(err, "failed to find wireguard interface")
	}

	linkaddr, err := netlink.ParseAddr(cidr)
	if err != nil {
		return errors.Wrap(err, "failed to parse wireguard interface ip address")
	}

	if err := netlink.AddrAdd(link, linkaddr); err != nil {
		return errors.Wrap(err, "failed to set ip address of wireguard interface")
	}

	return nil
}

func ConfigureForwarding(wgIface string, gatewayIface string) error {
	// Networking configuration (iptables) configuration
	// to ensure that traffic from clients the wireguard interface
	// is sent to the provided network interface
	ipt, err := iptables.New()
	if err != nil {
		return errors.Wrap(err, "failed to init iptables")
	}
	if err := ipt.AppendUnique("filter", "FORWARD", "-i", gatewayIface, "-o", wgIface, "-j", "ACCEPT"); err != nil {
		return errors.Wrap(err, "failed to set ip tables rule")
	}
	if err := ipt.AppendUnique("filter", "FORWARD", "-i", wgIface, "-o", gatewayIface, "-j", "ACCEPT"); err != nil {
		return errors.Wrap(err, "failed to set ip tables rule")
	}
	if err := ipt.AppendUnique("nat", "POSTROUTING", "-o", gatewayIface, "-j", "MASQUERADE"); err != nil {
		return errors.Wrap(err, "failed to set ip tables rule")
	}
	return nil
}
