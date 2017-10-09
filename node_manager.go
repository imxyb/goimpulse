package main

import (
	"context"
	"goimpulse/conf"
	"net/http"

	"goimpulse/lib"

	"goimpulse/sender"
	"time"

	"fmt"
	"io/ioutil"

	"github.com/coreos/etcd/client"
	"github.com/facebookgo/grace/gracehttp"
	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
)

var masterHost string

func main() {
	masterHost, _ = sender.GetMaster()
	go watchMasterHost()
	lib.OnReload()

	e := echo.New()
	e.GET("/getid", func(c echo.Context) error {
		result := map[string]interface{}{
			"id":   -1,
			"code": 0,
			"msg":  "success",
		}

		typeName := c.QueryParam("type")

		httpClient := &http.Client{}
		req, _ := http.NewRequest("GET", fmt.Sprintf("http://%s/getid?type=%s", masterHost, typeName), nil)

		user, pass, _ := c.Request().BasicAuth()
		req.SetBasicAuth(user, pass)

		if !lib.CheckAuth(req) {
			result["code"] = -1
			result["msg"] = "auth failed"
			return c.JSON(http.StatusForbidden, result)
		}

		res, err := httpClient.Do(req)

		if err != nil {
			pass := false
			for i := 0; i < 5; i++ {
				res, err = httpClient.Do(req)
				if err == nil {
					pass = true
					break
				}
				time.Sleep(800 * time.Millisecond)
			}
			if !pass {
				result["code"] = -1
				result["msg"] = "error"
				return c.JSON(http.StatusInternalServerError, result)
			}
		}

		data, _ := ioutil.ReadAll(res.Body)
		res.Body.Close()
		return c.JSONBlob(http.StatusOK, data)
	})

	log.Info("node_manager running...")
	e.Server.Addr = conf.Cfg.NodeManager.Host
	e.Server.SetKeepAlivesEnabled(false)
	e.Logger.Fatal(gracehttp.Serve(e.Server))
}

func watchMasterHost() {
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
}
