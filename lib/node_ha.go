package lib

import (
	"context"
	"goimpulse/conf"
	"time"

	"github.com/coreos/etcd/client"
	"github.com/labstack/gommon/log"
)

const MasterNode = "/goimpulse/master"

var Running chan bool = make(chan bool, 1)

func RegisterSelf() {
	ec := GetEtcd()
	kapi := client.NewKeysAPI(ec)

	opts := &client.SetOptions{TTL: 2 * time.Second, PrevExist: client.PrevNoExist}
	_, err := kapi.Set(context.Background(), MasterNode, conf.Cfg.App.Host, opts)
	if err != nil {
		go func() {
			log.Info("standby....")
			watcher := kapi.Watcher(MasterNode, nil)
			for {
				resp, _ := watcher.Next(context.Background())
				if resp.Action == "expire" {
					RegisterSelf()
				}
				time.Sleep(100 * time.Millisecond)
			}
		}()
	} else {
		log.Info("running....")
		Running <- true
		go reportHealth()
	}
}

func reportHealth() {
	ticker := time.NewTicker(300 * time.Millisecond)
LOOP:
	for {
		select {
		case <-ticker.C:
			ec := GetEtcd()
			kapi := client.NewKeysAPI(ec)
			opts := &client.SetOptions{TTL: 2 * time.Second, PrevExist: client.PrevExist, Refresh: true}
			resp, err := kapi.Set(context.Background(), MasterNode, "", opts)
			if err != nil || resp.PrevNode.Value != conf.Cfg.App.Host {
				Running <- false
				log.Info("standby....")
				break LOOP
			}
		}
	}
}
