package service

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io/ioutil"
	"math/rand"
	"net/http"
	"open/backend/gateway"
	"open/config"
	"open/downlink"
	"open/helpers"
	"open/uplink/data"
	"strconv"
	"time"
	"context"
	pb "open/backend/proto"
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"

)

const(
	BAIDU_TURN_OFF = "TurnOffRequest"
	BAIDU_TURN_ON = "TurnOnRequest"

)
type attr struct {
	Name                      string `json:"name"`
	Value                     string `json:"value"`
	Scale                     string `json:"scale"`
	TimestampOfSample         int    `json:"timestampOfSample"`
	UncertaintyInMilliseconds int    `json:"uncertaintyInMilliseconds"`
}

type Details struct {
	ExtraDetail1 string `json:"extraDetail1"`
	ExtraDetail2 string `json:"extraDetail2"`
	ExtraDetail3 string `json:"extraDetail3"`
	ExtraDetail4 string `json:"extraDetail4"`
}

type Appliance struct {
	Actions                    []string `json:"actions"`
	ApplianceTypes             []string `json:"applianceTypes"`
	AdditionalApplianceDetails Details  `json:"additionalApplianceDetails"`
	ApplianceID                string   `json:"applianceId"`
	FriendlyDescription        string   `json:"friendlyDescription"`
	FriendlyName               string   `json:"friendlyName"`
	IsReachable                bool     `json:"isReachable"`
	ManufacturerName           string   `json:"manufacturerName"`
	ModelName                  string   `json:"modelName"`
	Version                    string   `json:"version"`
	Attributes                 []attr   `json:"attributes"`
}

type reqBody struct {
	Header struct {
		Namespace      string `json:"namespace"`
		Name           string `json:"name"`
		MessageID      string `json:"messageId"`
		PayloadVersion string `json:"payloadVersion"`
	} `json:"header"`
	Payload struct {
		AccessToken string `json:"accessToken"`
		OpenUID     string `json:"openUid"`
		Appliance struct {
			AdditionalApplianceDetails struct {
				ExtraDetail1 string `json:"extraDetail1"`
				ExtraDetail2 string `json:"extraDetail2"`
				ExtraDetail3 string `json:"extraDetail3"`
				ExtraDetail4 string `json:"extraDetail4"`
			} `json:"additionalApplianceDetails"`
			ApplianceID string `json:"applianceId"`
		} `json:"appliance"`
	} `json:"payload"`

}
type findResult struct {
	Header struct {
		Namespace      string `json:"namespace"`
		Name           string `json:"name"`
		MessageID      string `json:"messageId"`
		PayloadVersion string `json:"payloadVersion"`
	} `json:"header"`
	Payload struct {
		DiscoveredAppliances []Appliance `json:"discoveredAppliances"`
		DiscoveredGroups     []struct {
			GroupName              string   `json:"groupName"`
			ApplianceIds           []string `json:"applianceIds"`
			GroupNotes             string   `json:"groupNotes"`
			AdditionalGroupDetails struct {
				ExtraDetail1 string `json:"extraDetail1"`
				ExtraDetail2 string `json:"extraDetail2"`
				ExtraDetail3 string `json:"extraDetail3"`
			} `json:"additionalGroupDetails"`
		} `json:"discoveredGroups"`
	} `json:"payload"`
}

type ControlRespond struct {
	Header struct {
		Namespace string `json:"namespace"`
		Name string `json:"name"`
		MessageID string `json:"messageId"`
		PayloadVersion string `json:"payloadVersion"`
	} `json:"header"`
	Payload struct {
		Attributes []interface{} `json:"attributes"`
		Err errorInfo `json:"errorInfo"`
	} `json:"payload"`
}
type errorInfo struct {
	Code string `json:"code"`
	Description string `json:"description"`
}

type errorResp struct {
	Header struct {
		Namespace string `json:"namespace"`
		Name string `json:"name"`
		MessageID string `json:"messageId"`
		PayloadVersion string `json:"payloadVersion"`
	} `json:"header"`
	Payload struct {

	} `json:"payload"`
}

type BaiduAPI struct {
	server *http.Server
	port string

}

func findHandle(c *gin.Context,find *reqBody) {
	var findResp findResult
	findResp.Header.Namespace = "DuerOS.ConnectedHome.Discovery"
	findResp.Header.MessageID = find.Header.MessageID
	findResp.Header.Name = "DiscoverAppliancesResponse"
	findResp.Header.PayloadVersion = find.Header.PayloadVersion
	attrSlice := []attr{
		attr{Name: "name", Value: "卧室的灯", Scale: "", TimestampOfSample: int(time.Now().Unix()), UncertaintyInMilliseconds: 10},
	}

	app1 := []Appliance{Appliance{
		Actions:        []string{"turnOn", "turnOff"},
		ApplianceTypes: []string{"LIGHT"},
		AdditionalApplianceDetails: Details{
			ExtraDetail1: "optionalDetailForSkillAdapterToReferenceThisDevice",
			ExtraDetail2: "There can be multiple entries",
			ExtraDetail3: "but they should only be used for reference purposes",
			ExtraDetail4: "This is not a suitable place to maintain current device state",
		},
		ApplianceID:         "uniqueLightDeviceId",
		FriendlyDescription: "friendlyDescription",
		FriendlyName:        "电灯",
		IsReachable:         true,
		ManufacturerName:    "设备制造商的名称",
		ModelName:           "fancyLight",
		Version:             "your software version number here.",
		Attributes:          attrSlice,
	}}
	findResp.Payload.DiscoveredAppliances = app1
	findResp.Payload.DiscoveredGroups = nil
	sjson, _ := json.Marshal(findResp)
	log.Info("respon to find req ", string(sjson))
	c.JSON(http.StatusOK, findResp)
	return

}

func controlHandle (c *gin.Context,req *reqBody) {
	var resp ControlRespond
	var cmd pb.Operation
	if req.Header.Name == BAIDU_TURN_ON {
		cmd =  pb.Operation_ON
	}else if req.Header.Name == BAIDU_TURN_OFF {
		cmd = pb.Operation_OFF
	}else {
		resp.Header.Namespace = "DuerOS.ConnectedHome.Control"
		resp.Header.Name = "UnsupportedTargetError"
		resp.Header.MessageID = req.Header.MessageID
		resp.Header.PayloadVersion = req.Header.PayloadVersion
		//resp.Payload.Attributes= []interface{}{}
		//log.Debug("[controlHandle]control respond")
		log.Info("[controlHandle]unsupported cmd,error",req.Header.Name)
		c.JSON(http.StatusOK, resp)
		return
	}
	downMsg := &pb.Payload{
		Kind:uint32(pb.Category_CMD) ,
		Key:uint32(pb.Device_LAMP),
		Val:[]byte{uint8(cmd)},
	}
	pData, err := proto.Marshal(downMsg)
	if err != nil {
		log.Debug("[Marshal]proto Marshal error")
		panic(err)
	}

	aes_data,err :=helpers.Encrypt(pData)
	if err != nil {
		log.Debug("[Marshal]Encrypt error")
		panic(err)
	}
	log.Debugf("origin data 0x %02X",pData)
	log.Debugf("after aes data, 0x%02X",aes_data)

	command := pb.DownlinkFrame{
		FrameType:gateway.ConfirmedDataDown,
		DevAddr:[]byte{0x30,0x31,0x32,0x33},
		FrameNum:1,
		Port:2,
		DownlinkId: uint32(rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(10000-1000)+1000),
		PhyPayload:aes_data,

	}
	downlink.AddDownlinkFrame(command)
	ch := make(chan pb.UplinkFrame,1)
	data.DownlinkTaskCache.SetDefault(strconv.Itoa(int(command.DownlinkId)),ch)
	log.Debug("wait up -------------")
	upFrame := <- ch
	if upFrame.UplinkId != command.DownlinkId  {
		log.Error("[controlHandle]UplinkId != DownlinkId",zap.Uint32("DownlinkId",command.DownlinkId),
			zap.Uint32("DownlinkId",upFrame.UplinkId))
		panic("UplinkId != DownlinkId")
	}
	//log.Debug("downID=",command.DownlinkId,",upID",upFrame.UplinkId)
	//log.Debugf("befor Decrypt upFrame PhyPayload = 0x%02x",upFrame.PhyPayload)
	//log.Debugf("befor Decrypt upFrame PhyPayload = ",string(upFrame.PhyPayload))

	payDecrypt,err := helpers.Decrypt(upFrame.PhyPayload)
	if err != nil {
		log.Error("[controlHandle]Decrypt error")
		panic(err)
	}
	upData :=  &pb.Payload{}
	err = proto.Unmarshal(payDecrypt, upData)
	if err != nil {
		log.Error("[controlHandle] proto Unmarshal error")
		panic(err)
	}

	if upData.Kind != uint32(pb.Category_CMD)  || upData.Key != uint32(pb.Device_LAMP){
		log.Errorf("[controlHandle] uplink's kind or key conflict。kind=%d,key=%d\n",upData.Kind,upData.Key)
		panic(err)
	}
	log.Debugf("upData val = 0x %02X",upData.Val)

	if len(upData.Val) >0 &&  upData.Val[0] == downMsg.Val[0] {
		// 执行成功。

		resp.Header.Namespace = "DuerOS.ConnectedHome.Control"
		resp.Header.Name = "TurnOnConfirmation"
		resp.Header.MessageID = req.Header.MessageID
		resp.Header.PayloadVersion = req.Header.PayloadVersion
		resp.Payload.Attributes= []interface{}{}
		log.Debug("[controlHandle]command success")

	}else {
		//执行失败

		resp.Header.Namespace = "DuerOS.ConnectedHome.Control"
		resp.Header.Name = "UnableToSetValueError"
		resp.Header.MessageID = req.Header.MessageID
		resp.Header.PayloadVersion = req.Header.PayloadVersion
		//resp.Payload.Attributes= []interface{}{}
		resp.Payload.Err = errorInfo{"DEVICE_AJAR","A custom description of the error.."}
		log.Info("[controlHandle]command failed")

	}

	c.JSON(http.StatusOK, resp)
	return
}
func ServiceAll(c *gin.Context){
	//header, _ := json.Marshal(c.Request.Header)
	input, _ := ioutil.ReadAll(c.Request.Body)
	// bodyjson, _ := json.Marshal(body)
	// encodeString := base64.StdEncoding.EncodeToString(input)
	// decodeJson, _ := base64.StdEncoding.DecodeString(string(input))
	// fmt.Println("body=", decodeJson)
	// log.Println("body=", string(decodeJson))
	//log.Debug("[ServiceAll]body=", string(input))
	var findReq reqBody
	json.Unmarshal(input, &findReq)
	log.Debug("[ServiceAll]Namespace=",findReq.Header.Namespace,"body=",string(input))
	switch findReq.Header.Namespace {
	case "DuerOS.ConnectedHome.Discovery" :
		findHandle(c,&findReq)
	case "DuerOS.ConnectedHome.Control":
		controlHandle(c,&findReq)
	default:
		log.Println("no match case")
		return
	}
	// c.String(http.StatusOK, "Hello World test2.go,/test")
}

func (ser *BaiduAPI) Start() error {
	router := gin.Default()
	//gin.SetMode(gin.ReleaseMode)

	router.POST("/", func(c *gin.Context) {
		fmt.Println("POST", "/")
		c.String(http.StatusOK, "Hello World test2.go,/")
	})
	router.GET("/", func(c *gin.Context) {
		fmt.Println("POST", "/")
		c.String(http.StatusOK, "Hello World test2.go,/")
	})
	router.POST("/test",  ServiceAll )
	router.GET("/test", func(c *gin.Context) {
		fmt.Println("POST", "test")
		c.String(http.StatusOK, "Hello World test2.go,/test")
	})
	router.GET("/oath/redirect", func(c *gin.Context) {
		fmt.Println("GET", "/oath/redirect")
		c.String(http.StatusOK, "Hello World test2.go,/test")
	})
	ser.server = &http.Server{
		Addr:    ":"+ser.port,
		Handler: router,
	}
	log.Debug("listen on :",ser.server.Addr)
	go func() {
		if err := ser.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()


	return nil
}
func (ser *BaiduAPI)Stop(){
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) //实际退出时间不是5s
	defer cancel()
	if err := ser.server.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:" + err.Error())
		//logger.GetLogger().Fatal("Server Shutdown:" + err.Error())
	}
	log.Debug("Server exiting")

}
func NewServer(conf config.Config) *BaiduAPI {

	return &BaiduAPI{
		port: strconv.Itoa(conf.General.HttpPort),
	}
}
func (ser *BaiduAPI) GetPort() string{
	return ser.port
}