package config
// Yaml2Go

var ParserDefault  = "smartLamp"

type Config struct {
	Mysql   Mysql   `yaml:"mysql"`
	Mqtt    Mqtt    `yaml:"mqtt"`
	General General `yaml:"general"`
	Consul MyConsul `yaml:"consul"`

}
type MyConsul struct {
	Server               string `yaml:"server"`

}
// Mysql
type Mysql struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Database string `yaml:"database"`
}

// Mqtt
type Mqtt struct {
	ClientID             string `yaml:"clientID"`
	EventTopic           string `yaml:"eventTopic"`
	CommandTopicTemplate string `yaml:"commandTopicTemplate"`
	Qos                  int    `yaml:"qos"`
	Password             string `yaml:"password"`
	Server               string `yaml:"server"`
	CleanSession         bool   `yaml:"cleanSession"`
	Username             string `yaml:"username"`
	MaxReconnectInterval int    `yaml:"maxReconnectInterval"`
	MaxHandleGoroutine int    `yaml:"maxHandleGoroutine"`

}

// General
type General struct {
	Type               string `yaml:"type"`
	DataCacheSec       int    `yaml:"data_cache_sec"`
	CheckCacheInterval int    `yaml:"check_cache_interval"`
	DownChannelSize    int    `yaml:"down_channel_size"`
	HttpPort           int    `yaml:"http_port"`
	LogLevel           string `yaml:"log_level"`
}




