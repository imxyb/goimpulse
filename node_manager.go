package main

import (
	"context"
	"goimpulse/conf"
	"net/http"

	"goimpulse/lib"

	"io"

	"goimpulse/sender"
	"time"

	"fmt"
	"io/ioutil"

	"github.com/coreos/etcd/client"
	log "github.com/sirupsen/logrus"
)

var masterHost string

func main() {
	masterHost, _ = sender.GetMaster()

	go func() {
		ec := lib.GetEtcd()
		kapi := client.NewKeysAPI(ec)
		watcher := kapi.Watcher(lib.MasterNode, nil)
		for {
			resp, _ := watcher.Next(context.Background())
			if resp.Action == "expire" || resp.Action == "create" {
				for i := 0; i < 3; i++ {
					host, err := sender.GetMaster()
					if err == nil {
						masterHost = host
						break
					}
					time.Sleep(100 * time.Millisecond)
				}
			}

			time.Sleep(100 * time.Millisecond)
		}
	}()

	http.HandleFunc("/getid", func(w http.ResponseWriter, req *http.Request) {
		query := req.URL.Query()
		typeName := query.Get("type")

		var res *http.Response
		var err error
		res, err = http.Get(fmt.Sprintf("http://%s/getid?type=%s", masterHost, typeName))

		if err != nil {
			pass := false
			for i := 0; i < 5; i++ {
				res, err = http.Get(fmt.Sprintf("http://%s/getid?type=%s", masterHost, typeName))
				if err == nil {
					pass = true
					break
				}
				time.Sleep(800 * time.Millisecond)
			}
			if !pass {
				http.Error(w, "error", http.StatusInternalServerError)
				return
			}
		}

		data, _ := ioutil.ReadAll(res.Body)
		res.Body.Close()
		w.Header().Add("Content-type", "application/json")
		io.WriteString(w, string(data))
	})

	log.Info("node_manager running")
	err := http.ListenAndServe(conf.Cfg.NodeManager.Host, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
