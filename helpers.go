package main

import (
	"net"
	"net/http"

	"github.com/Luzifer/ots/pkg/metrics"
	"github.com/Luzifer/ots/pkg/storage"
	"github.com/sirupsen/logrus"
)

func requestInSubnetList(r *http.Request, subnets []string) bool {
	if len(subnets) == 0 {
		// No subnets specififed: None allowed (without doing the parsing)
		return false
	}

	remote, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		logrus.WithError(err).Error("parsing remote address")
		return false
	}

	remoteIP := net.ParseIP(remote)
	if remoteIP == nil {
		logrus.WithError(err).Error("parsing remote address")
		return false
	}

	for _, sn := range subnets {
		_, netw, err := net.ParseCIDR(sn)
		if err != nil {
			logrus.WithError(err).WithField("subnet", sn).Warn("invalid subnet specified")
			continue
		}

		if netw.Contains(remoteIP) {
			return true
		}
	}

	return false
}

func updateStoredSecretsCount(store storage.Storage, collector *metrics.Collector) {
	n, err := store.Count()
	if err != nil {
		logrus.WithError(err).Error("counting stored secrets")
		return
	}
	collector.UpdateSecretsCount(n)
}
