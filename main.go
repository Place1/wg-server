package main

import (
	"net"
	"os"
	"os/signal"
	"syscall"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/pkg/errors"
	"github.com/place1/wg-server/pkg/wgserver"
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

	srv, err := wgserver.Start(*config, *gateway)
	if err != nil {
		logrus.Fatal(errors.Wrap(err, "failed to start server"))
	}

	logrus.Info("wireguard server started")

	term := make(chan os.Signal)
	signal.Notify(term, syscall.SIGTERM)
	signal.Notify(term, os.Interrupt)

	select {
	case <-term:
	}

	logrus.Info("shutting down")

	srv.Close()
}

func parseCIDR(cidr string) net.IPNet {
	_, peerAllowedIPs, err := net.ParseCIDR(cidr)
	if err != nil {
		panic(err)
	}
	return *peerAllowedIPs
}
