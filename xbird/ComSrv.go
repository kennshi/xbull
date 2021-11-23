package main

import (
	"fmt"
	"net"
	"os"

	//"strconv"
	"go.bug.st/serial.v1"
)

type ComSrv struct {
	_conn serial.Port
	_addr string
}

func ListComPorts() {
	//create serial port Object
	ports, err := serial.GetPortsList()
	if err != nil {
		fmt.Println("ListCommPorts err:", err)
	}
	if len(ports) == 0 {
		fmt.Println("ListCommPorts err:", "No serial ports found!")
	}

	for _, port := range ports {
		fmt.Printf("Found port: %v\n", port)
	}

}

func handleReq(conn net.Conn) {
	defer conn.Close()
	for {

	}
}

func (s ComSrv) Init(comport string) error {
	mode := &serial.Mode{
		BaudRate: 9600,
		DataBits: 8,
		Parity:   serial.NoParity,
		StopBits: serial.OneStopBit,
	}

	port, err := serial.Open(comport, mode)
	defer port.Close()
	if err != nil {
		fmt.Println("ComSrv Init Err:", err)
		return err
	}
	s._conn = port
	return nil
}

func (s ComSrv) Read(p []byte) (n int, err error) {
	return 0, nil
}

func (s ComSrv) Write(p []byte) (n int, err error) {
	return s._conn.Write(p)
}

func (s *ComSrv) IsConnected() bool {
	if s._conn == nil {
		return false
	}
	return true
}

func (s ComSrv) Close() error {
	os.Exit(0)
	return nil
}
