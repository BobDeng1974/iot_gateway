package service

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"time"
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
	} `json:"payload"`
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
	fmt.Println("respon to find ", string(sjson))
	c.JSON(http.StatusOK, findResp)
	return

}

func controlHandle (c *gin.Context,req *reqBody) {
	var resp ControlRespond
	resp.Header.Namespace = "DuerOS.ConnectedHome.Control"
	resp.Header.Name = "TurnOnConfirmation"
	resp.Header.MessageID = req.Header.MessageID
	resp.Header.PayloadVersion = req.Header.PayloadVersion
	resp.Payload.Attributes= []interface{}{}
	log.Println("[controlHandle]control respond")
	c.JSON(http.StatusOK, resp)
	return
}
func ServiceAll(c *gin.Context){
	fmt.Println("POST", "test")
	header, _ := json.Marshal(c.Request.Header)
	log.Println("header", string(header))
	input, _ := ioutil.ReadAll(c.Request.Body)
	// bodyjson, _ := json.Marshal(body)
	// encodeString := base64.StdEncoding.EncodeToString(input)
	// decodeJson, _ := base64.StdEncoding.DecodeString(string(input))
	// fmt.Println("body=", decodeJson)
	// log.Println("body=", string(decodeJson))
	log.Println("bodyxxxx=", string(input))
	var findReq reqBody
	json.Unmarshal(input, &findReq)
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