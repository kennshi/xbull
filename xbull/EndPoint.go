package main

import (
	"fmt"
	"math"
	"net"
	"time"

	//"sync"
	//"log"
	"encoding/hex"
)

type EndPoint struct {
	done       chan int
	conn       net.Conn
	sqlChannel chan string
	config     Config
	epcfg      EndPointAttr
}

func (this *EndPoint) Init(config Config, epcfg EndPointAttr, sqlChannel chan string) {
	this.sqlChannel = sqlChannel
	this.config = config
	this.epcfg = epcfg
	ipaddr := fmt.Sprint(epcfg.Ipaddr, ":", epcfg.Tcpport)
	Log.Println("EndPoint Init ipaddr:", ipaddr, " epid:", epcfg.EpId, " unit:", epcfg.Unit, " protNb:", epcfg.ProtNb)
	this.done = make(chan int, 1)
	go this.RunLoop()
}

func (this *EndPoint) RunLoop() {

	var dialer net.Dialer
	dialer.LocalAddr = &net.TCPAddr{IP: net.ParseIP("0.0.0.0")}
	defer func() {
		//this.RecvBuff <- []byte {}
		this.done <- 0
	}()

	config := this.config
	epcfg := this.epcfg
	sqlChannel := this.sqlChannel
	bufChannel := make(chan []byte, 8)
	valChannel := make(chan float64)
	defer func() { close(bufChannel); close(valChannel) }()

	ipaddr := fmt.Sprint(epcfg.Ipaddr, ":", epcfg.Tcpport)
	p := NewParser(epcfg.ProtNb)
	if p == nil {
		Err.Fatalln(ipaddr, "Create Protocol Parser Error! protNb:", epcfg.ProtNb)
	}
	p.Init(epcfg.ProtNb, valChannel)

	quit := make(chan int, 1)
	done := make(chan int, 1)
	defer func() { close(quit); close(done) }()
	for {
		sqlChannel <- buildEventSql(config.IdColliery, epcfg.EpId, EVENT_EP_Connecting, fmt.Sprintf("(%s:%d)", epcfg.Ipaddr, epcfg.Tcpport))
		// add connect status for UIConfig 2019-12-18
		status := fmt.Sprintf("status,%s,connecting", epcfg.EpId)
		sqlChannel <- status

		conn, err := dialer.Dial("tcp", ipaddr)
		if err != nil {
			Log.Println(ipaddr, "EndPoint::RunLoop dial failed:", err)
			// add connect status for UIConfig 2019-12-18
			status = fmt.Sprintf("status,%s,connectnk", epcfg.EpId)
			sqlChannel <- status

			time.Sleep(time.Duration(config.IntervalReConnect) * time.Second)
			continue
		}
		sqlChannel <- buildEventSql(config.IdColliery, epcfg.EpId, EVENT_EP_Connected, fmt.Sprintf("(%s:%d)", epcfg.Ipaddr, epcfg.Tcpport))

		// add connect status for UIConfig 2019-12-18
		status = fmt.Sprintf("status,%s,connected", epcfg.EpId)
		sqlChannel <- status

		go func() {
			//Log.Println(ipaddr, "Parser thrd Start!!")

			quitv := make(chan int, 1)
			donev := make(chan int, 1)

			defer func() {
				Log.Println(ipaddr, "Parser thrd stopped!!")
				close(quitv)
				close(donev)
			}()

			go func() {
				//Log.Println(ipaddr, "valueCollect thrd start!!")
				defer func() {
					Log.Println(ipaddr, "valueCollect thrd stopped!!")
				}()

				lastVal := 0.0
				lastValCount := 0
				lastDataTime := time.Now()
				lastWrited := false

				haveCar := false
				begTime := time.Now()
				endTime := time.Now()
				valTime := time.Now()
				savedValue := 0.0
				savedCount := 0

			ExitFor:
				for {
					select {
					case val := <-valChannel:
						if epcfg.Unit == 1 {
							val = val / 1000
						}

						if val != lastVal {
							if config.EventEPDataChange {
								datasecs := int(math.Round(time.Now().Sub(lastDataTime).Seconds()))
								info := fmt.Sprintf("[%s] N:%d S:%d Val:%.3f -> %.3f", lastDataTime.Format(standFormat), lastValCount, datasecs, lastVal, val)
								sql := buildEventSql(config.IdColliery, epcfg.EpId, EVENT_EP_DATA, info)
								this.sqlChannel <- sql
							}

							if !haveCar && lastVal == 0.0 && val != 0.0 {
								haveCar = true
								begTime = time.Now()
								savedValue = val
								savedCount = 1
								// haveCar begin......
								sql := buildProcSql(config.IdColliery, epcfg.EpId, 0.0, "carcome")
								this.sqlChannel <- sql
							}

							if haveCar && lastVal != 0.0 && val != 0.0 {
								if lastValCount > savedCount {
									savedValue = lastVal
									savedCount = lastValCount
									valTime = lastDataTime
								}
							}

							if haveCar && lastVal != 0.0 && val == 0.0 {
								datasecs := int(math.Round(time.Now().Sub(valTime).Seconds()))
								if savedValue >= config.MinEndPointData && datasecs >= config.DurationEndPointData {
									endTime = time.Now()
									begS := begTime.Format(standFormat)
									endS := endTime.Format(standFormat)
									valS := valTime.Format(standFormat)
									sql := buildDataSql(config.IdColliery, epcfg.EpId, begS, valS, endS, savedValue)
									this.sqlChannel <- sql
								}
								{
									sql := buildProcSql(config.IdColliery, epcfg.EpId, 0.0, "cargo")
									this.sqlChannel <- sql
								}
								haveCar = false
							}

							lastVal = val
							lastDataTime = time.Now()
							lastValCount = 1
							lastWrited = false
						} else {
							lastValCount += 1
							datasecs := int(math.Round(time.Now().Sub(lastDataTime).Seconds()))
							if !lastWrited && (lastVal >= config.MinEndPointData && datasecs >= config.DurationEndPointData) {
								sql := buildProcSql(config.IdColliery, epcfg.EpId, val, "")
								this.sqlChannel <- sql
								lastWrited = true
							}
						}
					case <-quitv:
						donev <- 1
						break ExitFor
					}
				}
			}()

			lastDataTime := time.Now()
		ExitFor:
			for {
				select {
				case buf := <-bufChannel:
					{
						lastDataTime = time.Now()
						n := len(buf)

						//Log.Println("main RunLoop n=",n, " msg:", string(buf))
						if n == 0 {
							break ExitFor
						}
						p.Parse(buf)
					}
				case <-time.After(time.Duration(15) * time.Second): //just for check dataTimeout
					{
						secs := int(math.Round(time.Now().Sub(lastDataTime).Seconds()))
						if secs > config.TimeoutEndPointData {
							Log.Println(ipaddr, "Now:", time.Now().Format(standFormat), " lastTime:", lastDataTime.Format(standFormat), " secs:", secs)
							lastDataTime = time.Now()
							sqlChannel <- buildEventSql(config.IdColliery, epcfg.EpId, EVENT_EP_TIMEOUT, "")
						}
					}
				case <-quit:
					quitv <- 1
					<-donev
					done <- 1
					break ExitFor
				}
			}
		}()

		this.conn = conn
		//buffer := make([]byte, 32)
		Log.Println(ipaddr, "EndPoint RunLoop RAddr:", conn.RemoteAddr(), " LAddr:", conn.LocalAddr())
		Log.Println(ipaddr, "EndPoint RunLoop start...")

		packetLen := p.getPacketLen()
		for {
			buffer := make([]byte, packetLen)
			n, err := conn.Read(buffer)
			if err != nil {
				Log.Println(ipaddr, "EndPoint RunLoop Read failed:", err)
				sqlChannel <- buildEventSql(config.IdColliery, epcfg.EpId, EVENT_EP_ConnBreak, fmt.Sprintf("(%s:%d)", epcfg.Ipaddr, epcfg.Tcpport))
				break
			}
			Log.Println(ipaddr, "Recved n=", n, " buffer:", string(buffer[:n]), " hex:", hex.EncodeToString(buffer[:n]))
			bufChannel <- buffer[:n]
		}
		quit <- 1
		<-done
		conn.Close()
		//this.conn = nil
	}
}

func (this *EndPoint) Close() {
	defer close(this.done)
	if this.conn != nil {
		this.conn.Close()
	}
	<-this.done
}

func NewEndPoint() *EndPoint {
	return &EndPoint{}
}
