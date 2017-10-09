package main

import (
	"goimpulse/conf"
	"goimpulse/lib"
	"goimpulse/sender"
	"net/http"

	"github.com/facebookgo/grace/gracehttp"
	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
)

var business map[string]*sender.Sequence
var service bool

func main() {
	go serviceInit()
	log.SetLevel(log.InfoLevel)
	lib.RegisterSelf()

	lib.OnReload()

	result := map[string]interface{}{
		"id":   -1,
		"code": 0,
		"msg":  "success",
	}

	e := echo.New()
	e.GET("/getid", func(c echo.Context) error {
		if !service {
			result["code"] = -1
			result["msg"] = "service offline"
			return c.JSON(http.StatusInternalServerError, result)
		}

		if !lib.CheckAuth(c.Request()) {
			result["id"] = -1
			result["code"] = -1
			result["msg"] = "auth failed"
			return c.JSON(http.StatusForbidden, result)
		}

		var seq *sender.Sequence
		var ok bool

		typeName := c.QueryParam("type")

		if typeName == "" {
			seq, ok = business["default"]
		} else {
			seq, ok = business[typeName]
		}

		if !ok {
			result["code"] = -1
			result["msg"] = "not found this type"
			return c.JSON(http.StatusNotFound, result)
		}

		id, ok := seq.GetId().(int64)
		if !ok {
			result["code"] = -1
			result["msg"] = "etcd error"
			return c.JSON(http.StatusInternalServerError, result)
		}

		result["code"] = 0
		result["msg"] = "success"
		result["id"] = id
		return c.JSON(http.StatusOK, result)
	})

	lib.LogPid(conf.Cfg.App.Host)
	e.Server.Addr = conf.Cfg.App.Host
	e.Server.SetKeepAlivesEnabled(false)
	e.Logger.Fatal(gracehttp.Serve(e.Server))
}

func serviceInit() {
	for {
		select {
		case run := <-lib.Running:
			if run {
				business = make(map[string]*sender.Sequence)
				business["default"] = sender.NewSeq("default")

				for _, typeName := range conf.Cfg.Type {
					business[typeName] = sender.NewSeq(typeName)
				}
				service = true
			} else {
				service = false
			}
		}
	}
}
