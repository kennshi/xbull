package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	//"os"
	//"flag"
)

type EndPointAttr struct {
	Ipaddr   string `json:"ipaddr"`
	Tcpport int    `json:"tcpport"`
	Unit    int    `json:"unit"` // 1,2,3
	ProtNb     int    `json:"protnb"`   //   -5 ～ 5
	EpId   string    `json:"epid"`
}

type Config struct {
	IdColliery string	`json:"id_colliery"`
	Dsn    string         `json:"dsn"`
	IntervalHeartbeatApp int  `json:"interval_heartbeat_app"`
	IntervalReConnect int `json:"interval_reconnect"`
	TimeoutEndPointData int `json:"timeout_endpoint_data"`
	DurationEndPointData int `json:"duration_endpoint_data"`
	MinEndPointData float64 `json:"min_endpoint_data"`
	NumEndPoint int `json:"num_endpoint"`
	EventEPDataChange bool `json:"event_epdata_change"`
	SaveEPDataZero bool `json:"save_epdata_zero"`
	EndPoints []EndPointAttr `json:"endpoints"`
}

func NewConfig() *Config {
	return &Config{}
}

func (m *Config) SaveToFile(cfg string) error {
	bytes, err := json.MarshalIndent(m, "", "\t")
	if err != nil {
		//log.Println("encoding xconf faild!")
		return err
	} else {
		//fmt.Println("encoding data :")
		//fmt.Println(bytes)
		//fmt.Println(string(bytes))
		ioutil.WriteFile(cfg, bytes, 0666) //os.ModeAppend)
	}
	return nil
}

func (m *Config) toJsonString() string {
	b, _ := json.MarshalIndent(m,"","\t")
	return string(b)
}

func (m *Config) fromJsonString(js string) {
	json.Unmarshal([]byte(js), m)
}

func (m *Config) LoadFromFile(cfg string) error {
	ra, err := ioutil.ReadFile(cfg)
	if err != nil {
		//log.Println(err)
		return err
	}
	err = json.Unmarshal(ra, m)
	if err != nil {
		//log.Println(err)
		return err
	}
	return nil
}

func ConfigtestCase(){

	m := Config {
		"140000000",
		"root:13811237916sS@tcp(localhost:3306)/bridge",
		600,
		60,
		600,
		40,
		1.50,
		3,
		false,
		true,
		[]EndPointAttr{
			{
				"192.168.12.1",
				8001,
				0,
				1,
				"001",
			},
			{
				"192.168.13.1",
				8002,
				1,
				2,
				"002",
			},
			{
				"192.168.13.1",
				8002,
				1,
				3,
				"003",
			},
		},
	}

	fmt.Println(m.toJsonString())
	m.SaveToFile("./xbull.json")

	m1 := NewConfig()
	m1.fromJsonString(m.toJsonString())
	fmt.Println(m1.toJsonString())
}
