package main

import (
	"goimpulse/conf"
	"goimpulse/lib"
	"goimpulse/sender"
	"net/http"

	"io"

	"encoding/json"

	log "github.com/sirupsen/logrus"
)

var business map[string]*sender.Sequence
var service bool

func main() {
	go serviceInit()
	log.SetLevel(log.InfoLevel)
	lib.RegisterSelf()

	http.HandleFunc("/getid", func(w http.ResponseWriter, req *http.Request) {
		if !service {
			return
		}

		if !lib.CheckAuth(w, req) {
			return
		}

		w.Header().Add("Content-type", "application/json")

		result := map[string]interface{}{
			"id":   -1,
			"code": 0,
			"msg":  "success",
		}

		query := req.URL.Query()

		var seq *sender.Sequence
		var ok bool
		typeName := query.Get("type")
		if typeName == "" {
			seq, ok = business["default"]
		} else {
			seq, ok = business[query.Get("type")]
		}

		if !ok {
			result["code"] = -1
			result["msg"] = "not found this type"
			data, _ := json.Marshal(result)
			http.Error(w, string(data), http.StatusNotFound)
			return
		}

		id, ok := seq.GetId().(int64)
		if !ok {
			result["code"] = -1
			result["msg"] = "etcd error"
			data, _ := json.Marshal(result)
			http.Error(w, string(data), http.StatusNotFound)
		}

		result["id"] = id
		data, _ := json.Marshal(result)
		io.WriteString(w, string(data))
	})

	err := http.ListenAndServe(conf.Cfg.App.Host, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
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
