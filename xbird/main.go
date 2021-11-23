package main

import (
	"fmt"
	//"time"
	"strconv"
	"log"
	"flag"
	"strings"
	//"math/rand"
	"math"
)

var comtype = flag.String("comtype", "tcp","comm type, 'com' or 'tcp'")
var tcpport = flag.String("tcpport", "6000", "tcp port number when comtype=='tcp'")
var comport = flag.String("comport","COM2", "name of comport when comtype=='com'")
var protocol = flag.Int("protocol",1,"protocol Number: 1,2,3")
var help = flag.Bool("h",false, "print this help")
var test = flag.Bool("test",false, "execute testcase..")
var listcomport = flag.Bool("listcomport", false, "list com port ")
var done = make(chan int)



func CmdHander(c chan string){
	//defer close(c)
	for{
		var sinput string
		fmt.Scanln(&sinput)
		fmt.Println(" your input is ", sinput)
		c <- sinput
		if sinput == "quit" {
			break;
		} 
	}
	fmt.Println("Quit from CmdHander")
	done <- 1
}

func DataHandler(cmd string){

}

func toString(f float64 ){
	v := math.Round(f * 1000) /1000
	s := fmt.Sprintf("%.3f",v)
	s1 := s[:]
	if len(s) > 8 {
		s1 =  s[:8]
	}
	s2 :=  strings.Repeat("0", 8- len(s1)) + s1
	fmt.Println(f, s2)
}

func main() {

	flag.Parse()
	if *help {
		flag.PrintDefaults()
		return
	}

	if *listcomport {
		ListComPorts()
		return 
	}

	if *test {
		//tCase(45678.9)
		//tCase(22.5)
		//tCase(1225.3)
		toString(0.0)
		toString(2)
		toString(0.1234567)
		toString(1.1)
		toString(123.456)
		toString(1234.56789)
		toString(12345.12345678)
		toString(1234.98765)
		toString(12345678)

		fmt.Println("------------")



		fmt.Println("0       :", fmatString(0,1))
		fmt.Println("0       :", fmatString(0,2))
		fmt.Println("0       :", fmatString(0,3))

		fmt.Println("445.6456   :", fmatString(445.6456,1))
		fmt.Println("445.6456   :", fmatString(445.6456,2))
		fmt.Println("445.6456   :", fmatString(445.6456,3))


		fmt.Println("55545.6456 :", fmatString(55545.6456,1))
		fmt.Println("55545.6456 :", fmatString(55545.6456,2))
		fmt.Println("55545.6456 :", fmatString(55545.6456,3))

		fmt.Println("555456.6456 :", fmatString(555456.6456,1))
		fmt.Println("555456.6456 :", fmatString(555456.6456,2))
		fmt.Println("555456.6456 :", fmatString(555456.6456,3))
		return 
	}

	log.Println(    "xBird Start.........")
	defer func(){
		<- done
		log.Println("xBird End...........")
	}()

	chInput := make(chan string)
	defer close(chInput)
	go CmdHander(chInput)

	protNumb := 1
	switch(*protocol){
	case 1:protNumb = 1 ;break
	case 2:protNumb = 2 ;break
	case 3:protNumb = 3 ;break
	default:break
	}

	commtype := *comtype
	if commtype != "com" && commtype != "tcp" {
		commtype = "tcp"
	}

	port := NewEndPoint(commtype)
	defer port.Close()

	fmt.Println("Create EndPoint :", port)

	if commtype == "tcp" {
		port.Init(*tcpport)
	}else {
		port.Init(*comport)
	}
	//fmt.Println(port, port._addr, port._conn)
	//go port.StartSrv()
	//time.Sleep(time.Millisecond * 500)
	//fmt.Println(port, port._addr, port._conn)
	datStr := ".000000=.000000=.000000=.000000=.000000=.000000=.000000=.000000=.000000=.000000=.000000=.000000=.000000=.000000=.000000=.000000=.000000=.000000=.000000=.000000=.000000=.000000=.000000=.000000=.000000=.000000=.000000=.000000=.000000=.000000=.080000=.082000=.026000=.084500=.047500=.049500=.049500=.049500=.029500=.009500=.083500=.063500=.083500=.085500=.029500=.043600=.003410=.006410=.042510=.044510=.007510=.000610=.007810=.087810=.028810=.028810=.068810=.048910=.088910=.020020=.003020=.047020=.063120=.028530=.089630=.069730=.042830=.047830=.063930=.063740=.048740=.081840=.061840=.061840=.041840=.041840=.041840=.061840=.061840=.061840=.061840=.081840=.002840=.022840=.062840=.003840=.043840=.004840=.048840=.068840=.068840=.068840=.068840=.068840=.068840=.068840=.068840=.068840=.068840=.068840=.068840=.068840=.088840=.088840=.088840=.088840=.088840=.088840=.009840=.009840=.029840=.029840=.049840=.049840=.049840=.069840=.089840=.089840=.089840=.089840=.000940=.000940=.000940=.000940=.000940=.000940=.000940=.000940=.000940=.000940=.000940=.000940=.000940=.000940=.000940=.000940=.000940=.000940=.000940=.000940=.089840=.089840=.089840=.089840=.089840=.089840=.089840=.089840=.089840=.000940=.020940=.020940=.020940=.020940=.020940=.020940=.020940=.020940=.020940=.020940=.020940=.020940=.020940=.020940=.000940=.000940=.000940=.000940=.000940=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.049840=.069840=.020940=.020940=.020940=.020940=.020940=.020940=.020940=.000940=.000940=.000940=.000940=.000940=.000940=.089840=.089840=.089840=.089840=.089840=.089840=.089840=.089840=.089840=.089840=.089840=.089840=.089840=.089840=.089840=.089840=.089840=.089840=.089840=.089840=.089840=.000940=.000940=.000940=.000940=.000940=.000940=.000940=.000940=.000940=.000940=.000940=.000940=.000940=.000940=.000940=.000940=.000940=.000940=.000940=.000940=.000940=.020940=.020940=.020940=.020940=.020940=.020940=.020940=.020940=.000940=.089840=.069840=.029840=.000640=.008540=.063540=.080540=.088440=.066440=.004240=.003240=.022240=.022240=.041240=.009140=.027140=.042140=.045040=.064820=.025720=.041720=.041720=.021720=.021720=.001720=.060720=.000720=.086620=.041620=.063520=.049010=.067900=.021900=.027800=.081800=.025700=.020000=.020000=.020000=.020000=.020000=.000000=.000000=.000000=.000000=.000000=.000000=.000000=.000000=.000000=.000000=.000000=.000000=.000000=.000000=.000000=.000000=.000000=.000000=.000000=.000000=.000000=.000000=.000000=.000000=.000000=.000000=.000000=.000000="
	datList:=strings.Split(datStr,"=")
	tunList := []float64{}

	for _, str := range datList {
		s := invertString(str)
		v, err := strconv.ParseFloat(s, 64)
		if err == nil {
			tunList = append(tunList, v)
		}
	
	}
	//fmt.Println(tunList)

	for{

		
		sCmd := <- chInput
		if sCmd == "quit" {
			break
		}

		if !port.IsConnected() {
			log.Println("wait for port Connecting....")
			continue
		}

		if sCmd == "mode3" {
			sender := tunsender { port, 0.0 , protNumb }
			sender.sendMode3()
		} else if sCmd == "data" {
			sender := tunsender { port, 0.0 , protNumb }
			sender.sendList(tunList)
		}else {
			tun, err := strconv.ParseFloat(sCmd,64)
			if err == nil{
				if tun >= 0 && tun <= 999999 {
					sender := tunsender { port, tun , protNumb }
					sender.send()
				}else{
					fmt.Println(" Sorry , floatNumber too big or == zero !!")
				}
			}else {
				fmt.Println("strconv.ParseFloat Error : ",  err)
			}

		}
	}

	//随机数，生成称重列表
	//r := rand.New(rand.NewSource(time.Now().Unix()))
	//ton := 35.55
	//tonbuff := make([]float32, 1024)
	//模拟一台地磅的数据发送过程,无车时发送全0,超重时发送全9, 一直发送间隔为100ms
	fmt.Println("Quit from Main")
}
