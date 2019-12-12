package main

import (
	"net"
	"os"
	"os/signal"
	"syscall"

	"golang.zx2c4.com/wireguard/wgctrl"
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/pkg/errors"
	"github.com/place1/wg-server/pkg/wgembed"
	"github.com/sirupsen/logrus"
)

var (
	app     = kingpin.New("wg-server", "a cli app to run a simple userspace wireguard server")
	config  = app.Arg("config", "a wireguard configuration file, missing then config will be read from stdin").String()
	debug   = app.Flag("debug", "enable verbose logging").Bool()
	iface   = app.Flag("iface", "wirguard interface name").Default("wg0").String()
	gateway = app.Flag("gateway", "the gateway nic").Default("eth0").String()
	ip      = app.Flag("iface-ip", "the ip address of the wireguard interface").Default("10.50.0.1").String()
)

func main() {
	kingpin.MustParse(app.Parse(os.Args[1:]))

	if *debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	if *config == "" {
		*config = "/dev/stdin"
	}

	opts, err := wgembed.ReadConfig(*config)
	if err != nil {
		logrus.Fatal(errors.Wrap(err, "failed to read wireguard config file"))
	}

	logrus.Info("starting interface")
	wg0, err := wgembed.New(*iface)
	if err != nil {
		logrus.Fatal(errors.Wrap(err, "failed to create wg interface"))
	}

	client, err := wgctrl.New()
	if err != nil {
		logrus.Fatal(errors.Wrap(err, "failed to create wg client"))
	}

	err = client.ConfigureDevice(*iface, opts.Config())
	if err != nil {
		logrus.Fatal(errors.Wrap(err, "failed to configure wireguard"))
	}

	logrus.Info("interface configured")

	if err := wg0.Up(); err != nil {
		logrus.Fatal(errors.Wrap(err, "failed to bring interface up"))
	}

	if err := wg0.SetIP(*ip); err != nil {
		logrus.Fatal(errors.Wrap(err, "failed to set interface ip"))
	}

	if err := wg0.ConfigureForwarding(*gateway); err != nil {
		logrus.Fatal(errors.Wrap(err, "failed to set interface forwarding"))
	}

	term := make(chan os.Signal)
	signal.Notify(term, syscall.SIGTERM)
	signal.Notify(term, os.Interrupt)

	select {
	case <-term:
	}

	logrus.Info("shutting down")

	wg0.Close()
}

func parseCIDR(cidr string) net.IPNet {
	_, peerAllowedIPs, err := net.ParseCIDR(cidr)
	if err != nil {
		panic(err)
	}
	return *peerAllowedIPs
}
