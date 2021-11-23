package main

import (
	"log"
	"time"
	"math/rand"
	"fmt"
	"encoding/hex"
)

type tunsender struct {
	port EndPoint
	tun float64
	protNumb int 
}


func tCase(f float64) []float64 {
	f1 := f / 30
	rand.Seed(int64(time.Now().UnixNano()))

	fList :=[]float64{0,0}
	f2 := 0.0
	for i:=0;i<30;i++ {
		if f2 > f {
			break
		}
		n:= rand.Intn(6)
		for j:=0;j<n;j++ {
			fList = append(fList, f2)
		}
		f2 += (f1 + rand.Float64())
	}
	for i:=0 ;i<25;i++ {
		fList= append(fList, f)
	}

	f2 = f
	for i:=30;i>0;i-- {
		f2 -= (f1+rand.Float64())
		if f2 < 0 {
			break
		}
		n:= rand.Intn(6)
		for j:=0;j<n;j++ {
			fList = append(fList, f2)
		}
	}
	for i:=0 ;i<5;i++ {
		fList= append(fList, 0)
	}
	return fList
}

func(s *tunsender) sendMode3(){
	port := s.port
	//protNumb := s.protNumb
	//if port == nil {
	//	log.Println("tunSender:send Err (port == nil)")
	//	return 
	//}

	//3033313803022b3030303030
	stun := "0318\x03\x02+00000"

	for n:=0;n<50;n++ {
		port.Write([]byte(stun))
		fmt.Println("s:<", stun , "> hex:<" , hex.EncodeToString([]byte(stun)), ">")
		time.Sleep(time.Second)
	}
}

func(s *tunsender) sendList(tunList []float64){
	port := s.port
	protNumb := s.protNumb
	if port == nil {
		log.Println("tunSender:send Err (port == nil)")
		return 
	}
	for _,v := range tunList {
		stun := fmatString(v,protNumb)
		port.Write([]byte(stun))
		fmt.Println("s:<", stun , "> hex:<" , hex.EncodeToString([]byte(stun)), ">")
	}
}
func(s *tunsender) send() {
	tun := s.tun
	tunList := tCase(tun)
	s.sendList(tunList)
}
