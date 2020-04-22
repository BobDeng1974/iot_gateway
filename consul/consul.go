package consul

import (
	"iot_gateway/config"
	consulapi "github.com/hashicorp/consul/api"
	"sync"
	log "github.com/sirupsen/logrus"
	"github.com/jasonlvhit/gocron"
)

var ServiceMap sync.Map
func Setup(conf config.Config) error {

	config := consulapi.DefaultConfig()
	config.Address = conf.Consul.Server // consul 服务地址
	client, err := consulapi.NewClient(config)
	if err != nil {
		log.Fatal("consul client error : ", err)
	}
	updateAllService := func() {
		services, err := client.Agent().Services()
		if err != nil {
			log.Warn("[consul][Setup]error,addr",config.Address,err)
			return
		}
		eraseSyncMap(&ServiceMap)
		for k, value := range services {
			ServiceMap.Store(k,value)
		}
		return
	}
	updateAllService()
	serviceUpdate := gocron.NewScheduler()
	_ = serviceUpdate.Every(2).Seconds().Do(func() {
		updateAllService()
	})
	serviceUpdate.Start()
	return nil
}

func eraseSyncMap(m *sync.Map) {
	m.Range(func(key interface{}, value interface{}) bool {
		m.Delete(key)
		return true
	})
}
