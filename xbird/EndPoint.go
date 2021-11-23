package main

type EndPoint interface {
	Init(attr string) error
	Read(p []byte) (n int, err error)
	Write(p []byte) (n int, err error)
	IsConnected() bool
	Close() error
}


func NewEndPoint(ep string) EndPoint {
	if ep == "com" {
		return &ComSrv{}
	} else if ep == "tcp" {
		return &TcpSrv{}
	}
	return nil
}