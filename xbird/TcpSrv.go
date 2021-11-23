package main

import (
	"net"
	"fmt"
	"os"
	"time"
	//"strconv"
)

type TcpSrv struct{
	_conn net.Conn
	_addr string

	done chan bool
}

func(s *TcpSrv) StartSrv() {
	fmt.Println("TcpSrv StartSrv.........")
	defer func(){
		fmt.Println("TcpSrv quited 1.........")
		s.done <- true
	}()
	l, err := net.Listen("tcp", s._addr)
	if err != nil {
		fmt.Println("TcpSrv Error listening:", err)
		os.Exit(1)
	}
	fmt.Println("TcpSrv Listening on ", s._addr)

	go func(){
		for {
			fmt.Println("TcpSrv Wait for new conntion..")
			conn,err := l.Accept()
			if err != nil {
				fmt.Println("TcpSrv Error accepting:", err)
				break
			}
	
			fmt.Printf("TcpSrv Reveived Conn %s -> %s \n", conn.RemoteAddr(), conn.LocalAddr())
			s._conn = conn
			buff := make([]byte,1)
			n,err := conn.Read(buff)
			//s._conn = nil
			//conn.Close()
			fmt.Println("TcpSrv conn.Read n=",n, "Err:",err)
		}
		s.done <- true
	}()
	
	done := <- s.done
	if done {
		if s._conn != nil {
			s._conn.Close()
		}
		l.Close()
	}
	<- s.done
	

	fmt.Println("TcpSrv quited 2.........")
	//s.done <- true
}


func(s *TcpSrv) Init(port string) error {
	fmt.Println("TcpSrv Init....")
	s._addr = "0.0.0.0:" + port
	s.done = make(chan bool)
	//defer close(s.done)
	go s.StartSrv()
	fmt.Println("TcpSrv Init end....")
	return nil
}

func(s *TcpSrv) Read(p []byte) (n int , err error){
	return 0, nil
}

func(s *TcpSrv) Write(p [] byte) (n int , err error){
	if s._conn == nil {
		return 0,nil
	}
	n,err = s._conn.Write(p)
	time.Sleep(100 * time.Millisecond)
	return n,err
}

func(s *TcpSrv) IsConnected() bool{
	if s._conn == nil {
		return false
	} else {
		return true
	}
}

func(s *TcpSrv) Close() error{
	fmt.Println("TcpSrv will Close.....")
	defer close(s.done)
	s.done <- true
	time.Sleep(100 * time.Millisecond)
	<- s.done
	fmt.Println("TcpSrv  Closed.....")
	return nil
}