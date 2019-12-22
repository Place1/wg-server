package wgserver

import (
	"github.com/pkg/errors"
	"github.com/place1/wg-embed/pkg/wgembed"
	"github.com/sirupsen/logrus"
)

type WireGuardServer struct {
	iface *wgembed.WireGuardInterface
}

func Start(configFile string, gateway string) (*WireGuardServer, error) {
	srv := &WireGuardServer{}

	wg0, err := wgembed.New("wg0")
	if err != nil {
		logrus.Fatal(errors.Wrap(err, "failed to create wg interface"))
	}
	srv.iface = wg0

	if err := wg0.LoadConfigFile(configFile); err != nil {
		logrus.Fatal(errors.Wrap(err, "failed to configure wireguard"))
	}

	if err := ConfigureForwarding(wg0, gateway); err != nil {
		logrus.Fatal(errors.Wrap(err, "failed to set interface forwarding"))
	}

	return srv, nil
}

func (srv *WireGuardServer) Close() error {
	return srv.iface.Close()
}
