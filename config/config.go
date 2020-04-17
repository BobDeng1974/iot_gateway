package config

import (
	"fmt"
	"io/ioutil"
	log "github.com/sirupsen/logrus"
	"os"
	"gopkg.in/yaml.v2"

	"path/filepath"
)
var C Config //在root.go中对config进行初始化

func init() {
	runDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	configFilePath := "config.yaml"
	runEnv := os.Getenv("NS_ENV")
	if runEnv != "" {
		configFilePath = "config_" + runEnv + ".yaml"
	}
	configFilePath = filepath.Join(runDir, configFilePath)
	yamlFile, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		fmt.Println(err)
	}
	err = yaml.Unmarshal(yamlFile, &C)
	if err != nil {
		log.Fatal("config init fail")
	}
	fmt.Println("[init]config=", C)

}