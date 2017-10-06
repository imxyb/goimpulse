package sender

import (
	"context"
	"goimpulse/lib"

	"fmt"

	"github.com/coreos/etcd/client"
	log "github.com/sirupsen/logrus"
)

type Sender interface {
	GetId() interface{}
}

func GetLastId(name, defaultVal string) string {
	ec := lib.GetEtcd()
	api := client.NewKeysAPI(ec)
	resp, err := api.Get(context.Background(), GetLastIdKey(name), nil)

	if err != nil {
		return defaultVal
	}

	return resp.Node.Value
}

func SaveLastId(typeName, lastId string) {
	ec := lib.GetEtcd()
	api := client.NewKeysAPI(ec)

	lastIdKey := GetLastIdKey(typeName)
	_, err := api.Set(context.Background(), lastIdKey, lastId, nil)
	if err != nil {
		log.Warn("fail to save lastId")
		// 重试三次
		for i := 0; i < 3; i++ {
			_, err := api.Set(context.Background(), lastIdKey, lastId, nil)
			if err == nil {
				break
			}
		}
	}
}

func GetLastIdKey(typeName string) string {
	return fmt.Sprintf("/goimpulse/%s/lastid", typeName)
}

func GetMaster() (string, error) {
	ec := lib.GetEtcd()
	kapi := client.NewKeysAPI(ec)
	res, err := kapi.Get(context.Background(), lib.MasterNode, nil)

	if err != nil {
		return "", err
	}

	return res.Node.Value, nil
}
