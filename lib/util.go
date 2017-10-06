package lib

import (
	"goimpulse/conf"
	"net/http"

	"github.com/coreos/etcd/client"
)

func GetEtcd() client.Client {
	cfg := client.Config{
		Endpoints: []string{conf.Cfg.Etcd.Host},
	}
	c, _ := client.New(cfg)
	return c
}

func CheckAuth(w http.ResponseWriter, req *http.Request) bool {
	if !conf.Cfg.Auth.Enable {
		return true
	}

	username, password, _ := req.BasicAuth()

	if username != conf.Cfg.Auth.User || password != conf.Cfg.Auth.Pass {
		http.Error(w, "auth failed", http.StatusForbidden)
		return false
	}

	return true
}
