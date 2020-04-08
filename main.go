//go:generate echo generate parser.proto
//go:generate protoc -I ./lib/grpc_service/parser_proto --go_out=plugins=grpc:./lib/grpc_service/parser_proto ./lib/grpc_service/parser_proto/parser.proto
//go:generate echo generate gateway.proto
//go:generate protoc -I ./backend/proto --go_out=plugins=grpc:./backend/proto ./backend/proto/gateway.proto

package main

import (
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"open/backend/gateway"
	"open/backend/gateway/mqtt"
	"open/config"
	"open/consul"
	"open/downlink"
	"open/service"
	"open/uplink"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	fmt.Println("start")

	var server = new(uplink.Server)
	var baiDuSer = new( service.BaiduAPI)
	tasks := []func() error{
		setLogLevel,
		setGatewayBackend,
		setupUplink,
		setupDownlink,
		setupServicefind,
		startBaiDuService(baiDuSer),
		startLoRaServer(server),
	}

	for _, t := range tasks {
		if err := t(); err != nil {
			log.Fatal(err)
		}
	}
	sigChan := make(chan os.Signal)
	exitChan := make(chan struct{})
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	log.WithField("signal", <-sigChan).Info("signal received")
	go func() {
		log.Warning("stopping chirpstack-network-server")
		if err := server.Stop(); err != nil {
			log.Fatal(err)
		}
		//if err := gwStats.Stop(); err != nil {
		//	log.Fatal(err)
		//}
		exitChan <- struct{}{}
	}()
	//两种结束程序的方式
	select {
	case <-exitChan:
	case s := <-sigChan:
		log.WithField("signal", s).Info("signal received, stopping immediately")
	}

}

func setLogLevel() error {
	level,_ := log.ParseLevel(config.C.General.LogLevel)
	log.SetLevel(level)
	//log.SetReportCaller(true)
	return nil
}

func setupUplink() error {
	if err := uplink.Setup(config.C); err != nil {
		return errors.Wrap(err, "setup link error")
	}
	return nil
}
func setupDownlink() error {

	if err := downlink.Setup(config.C); err != nil {
		return errors.Wrap(err, "setup setupBaiDuService error")
	}
	return nil

}
func startBaiDuService(baiDuSer *service.BaiduAPI) func() error {
	return func() error {
		baiDuSer = service.NewServer(config.C)
		return  baiDuSer.Start()
	}
}
func startLoRaServer(server *uplink.Server) func() error {
	return func() error {
		*server = *uplink.NewServer()
		return server.Start()
	}
}
func setupServicefind() error {
	if err := consul.Setup(config.C);err != nil {
		return errors.Wrap(err, "setupServicefind error")
	}
	return nil
}
func setGatewayBackend() error {
	var err error
	var gw gateway.Gateway
	//log.Info("[ldm]Backend.Type = ",config.C.NetworkServer.Gateway.Backend.Type)
	switch config.C.General.Type {
	case "mqtt":
		gw, err = mqtt.NewBackend( // 最终赋值给一个接口对象
			config.C,
		)
	case "amqp":
		//gw, err = amqp.NewBackend(config.C)
		log.Info("[setGatewayBackend] amqp unsupport protocol")
	case "gcp_pub_sub":
		//gw, err = gcppubsub.NewBackend(config.C)
		log.Info("[setGatewayBackend] gcp_pub_sub unsupport protocol")
	case "azure_iot_hub":
		//gw, err = azureiothub.NewBackend(config.C)
		log.Info("[setGatewayBackend] azure_iot_hub unsupport protocol")
	default:
		return fmt.Errorf("unexpected gateway backend type: %s", config.C.General.Type)
	}

	if err != nil {
		return errors.Wrap(err, "gateway-backend setup failed")
	}

	gateway.SetBackend(gw) // 将网关对象赋值给网关接口。因为网关可以是mqtt协议的，也可以是其他协议上传来的
	return nil
}


//func startQueueScheduler() error {
//	log.Info("starting downlink device-queue scheduler")
//	go downlink.DeviceQueueSchedulerLoop()
//	return nil
//}