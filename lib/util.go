package lib

import (
	"goimpulse/conf"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"strconv"

	"context"

	"github.com/coreos/etcd/client"
)

func GetEtcd() client.Client {
	cfg := client.Config{
		Endpoints: conf.Cfg.Etcd.Host,
	}
	c, _ := client.New(cfg)
	return c
}

func CheckAuth(req *http.Request) bool {
	if !conf.Cfg.Auth.Enable {
		return true
	}

	username, password, _ := req.BasicAuth()

	if username != conf.Cfg.Auth.User || password != conf.Cfg.Auth.Pass {
		return false
	}

	return true
}

func OnReload() {
	var sig os.Signal
	signalChan := make(chan os.Signal)
	signal.Notify(
		signalChan,
		syscall.SIGUSR1,
	)

	go func() {
		for {
			sig = <-signalChan
			switch sig {
			case syscall.SIGUSR1:
				conf.LoadConfig()
			default:
			}
		}
	}()
}

const PidKey = "/goimpulse/pid/"

func LogPid(key string) {
	ec := GetEtcd()
	api := client.NewKeysAPI(ec)
	pidStr := strconv.Itoa(os.Getpid())
	api.Set(context.Background(), PidKey+key, pidStr, nil)
}

func GetPid(filename string) (int, error) {
	f, err := os.Open(filename)
	if err != nil {
		return 0, err
	}
	var data []byte
	f.Read(data)

	pid, _ := strconv.Atoi(string(data))
	return pid, nil
}
