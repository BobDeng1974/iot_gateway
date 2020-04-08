package data

import (
	"context"
	"strconv"

	"github.com/patrickmn/go-cache"
	//log "github.com/sirupsen/logrus"
	pb "open/backend/proto"
	"open/config"
	"time"
)
var DownlinkTaskCache *cache.Cache // key是token，val是channel

func Handle(ctx context.Context, rxPacket *pb.UplinkFrame) error {
	//go cach. 如果上行帧中有UplinkId 字段，则说明刚刚有任务下行，此上线包中包含对任务的反馈
	if  ch ,found := DownlinkTaskCache.Get(strconv.Itoa(int(rxPacket.UplinkId))); found {

		 ch.(chan *pb.UplinkFrame) <- rxPacket
	}

	return nil
}


func Setup( conf config.Config) error{
	DownlinkTaskCache = cache.New(time.Duration(conf.General.DataCacheSec)*time.Second, time.Duration(conf.General.CheckCacheInterval)*time.Second)
	return nil
}