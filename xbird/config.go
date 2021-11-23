package main

import (
	"os"
	"fmt"
	"encoding/json"
)

type xConf struct {
	CommType string // com or tcp
	CommPort string // 
	Ipaddr string
	TcpPort int
}


func getConfig() (xConf, bool) {
	file, err := os.Open("xBull.json")
	defer file.Close()
	conf := xConf{}
	if err != nil {
		fmt.Println("getConfig Error:", err)
		return conf,false 
	}
	decoder := json.NewDecoder(file)
	
	err = decoder.Decode(&conf)
	if  err != nil {
		fmt.Println("getConfig Error:", err)
		return conf, false 
	}
	return conf, true
}