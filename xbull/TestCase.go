package main

import (
	"reflect"
	"unsafe"
	"fmt"
)


func testCaseSlice(){
	sa := make([]byte, 0 ,64)
	fmt.Print(0,"len=", len(sa)," cap=", cap(sa), "addr=", unsafe.Pointer(&sa), " val=", string(sa)); fmt.Printf( "addr= %p\n", sa)
	sb := []byte("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	//copy(sa, sb[:])
	sa = sb[:]
	fmt.Print(0,"len=", len(sa)," cap=", cap(sa), "addr=", unsafe.Pointer(&sa[0]), " val=", string(sa)); fmt.Printf( "addr= %p\n", sa)

	nlen := len(sa)
	for i:=1;i<nlen;i++ {
		sa = sa[1:]
		fmt.Print(i,"len=", len(sa)," cap=", cap(sa), "addr=", unsafe.Pointer(&sa[0]), " val=", string(sa)); fmt.Printf( "addr= %p\n", sa)
	}

	sa = append(sa, []byte("abcde")[:]...)
	fmt.Print(99,"len=", len(sa)," cap=", cap(sa), "addr=", unsafe.Pointer(&sa[0]), " val=", string(sa)); fmt.Printf( "addr= %p\n", sa)

	sc := sa


	addr := unsafe.Pointer(&sa[0])
	addr1 := unsafe.Pointer(&sc[0])
	fmt.Println("addr=", addr , " addr1=", addr1, reflect.TypeOf(addr))
	if addr == addr1 {
		fmt.Println("it is same addr!!")
	}
	return 
}