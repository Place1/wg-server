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
	app    = kingpin.New("wg-server", "a cli app to run a simple userspace wireguard server")
	config = app.Arg("config", "a wireguard configuration file").Required().File()
	debug  = app.Flag("debug", "enable verbose logging").Bool()
)

func main() {
	kingpin.MustParse(app.Parse(os.Args[1:]))

	if *debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	opts, err := wgembed.ReadConfig((*config).Name())
	if err != nil {
		logrus.Fatal(errors.Wrap(err, "failed to read wireguard config file"))
	}

	logrus.Info("starting wg0")
	wg0, err := wgembed.StartInterface("wg0")
	if err != nil {
		logrus.Fatal(errors.Wrap(err, "failed to create wg interface"))
	}

	client, err := wgctrl.New()
	if err != nil {
		logrus.Fatal(errors.Wrap(err, "failed to create wg client"))
	}

	err = client.ConfigureDevice("wg0", opts.Config())
	if err != nil {
		logrus.Fatal(errors.Wrap(err, "failed to configure wireguard"))
	}

	logrus.Info("wg0 configured")

	if err := wgembed.Up("wg0"); err != nil {
		logrus.Fatal(errors.Wrap(err, "failed to bring wg0 up"))
	}

	if err := wgembed.SetIP("wg0", "10.50.0.1/32"); err != nil {
		logrus.Fatal(errors.Wrap(err, "failed to set wg0 ip"))
	}

	if err := wgembed.ConfigureForwarding("wg0", "eth0"); err != nil {
		logrus.Fatal(errors.Wrap(err, "failed to set wg0 forwarding"))
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
