package main

import (
	"strconv"
	"sort"
)

type IParser interface {
	Parse(buf []byte)
	Init(protNb int, valChannel chan float64)
	getPacketLen() int
}

type Parser struct {
	buffer []byte 
	//haveCar bool
	count map[string]int
	packetLen int
	valChannel chan float64
}

func(m *Parser) getPacketLen() int {
	return m.packetLen
}

func(m *Parser) Init(protNb int,valChannel chan float64){
	mp := map[int] int {
		1: 8,
		2: 9,
		3: 12,
	}	
	ir,ok := mp[protNb]
	if ok {
		m.packetLen = ir
	} else {
		Err.Fatalln("Parser:Init ProtNb Error !")
	}

	m.valChannel = valChannel
	if m.count == nil {
		m.count = make(map[string]int)
	}
}

func(m *Parser) toFloat(s string) (float64, error) {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil{
		return -1.0, err
	}
	return f,err
}

func(m *Parser) saveValue(val string){
	// if v != 0.0 {
	// 	if !m.haveCar {
	// 		for k := range m.count{
	// 			delete(m.count,k) //删除整个字典的数据
	// 		}
	// 		m.haveCar = true
	// 	}
	// 	m.count[sval] ++
	// } else if v == 0.0{
	// 	if m.haveCar {
	// 		v :=m.FindValue()
	// 		ch <- v
	// 		m.haveCar = false
	// 	}
	// }	
}

func(m *Parser) FindValue() float64 {
	type kv struct {
        Key   string
        Value int
    }

    var ss []kv
    for k, v := range m.count {
        ss = append(ss, kv{k, v})
    }

    sort.Slice(ss, func(i, j int) bool {
        return ss[i].Value > ss[j].Value  // 降序
        // return ss[i].Value < ss[j].Value  // 升序
    })

	s := ss[0].Key;

	Log.Println("FindValue :", s)

	s = invertString(s)
	v , err :=strconv.ParseFloat(s,64)
	if err != nil {
		Log.Println("FindValue Err! :", err)
	}
	return v
}


func(m *Parser) Parse(buf []byte) {
	n := len(buf)
	if n == 0 {
		return
	}

	buffer :=  m.buffer
	buffer = append(buffer, buf...)
	packetLen := m.packetLen
	for {
		bufLen := len(buffer)
		if bufLen < packetLen {
			break
		}

		bufChanged := false 
		for i:=0 ;i<bufLen;i++ {
			if buffer[i] == '=' {
				val := buffer[:(i)]
				sval := invertString(string(val))
				v, _ := m.toFloat(sval)
				m.valChannel <- v
				//log.Println("Parser1.Parse val:", string(val))
				buffer = buffer[(i+1):]
				bufChanged = true
				break
			}
		}

		if !bufChanged {
			break
		}
	}
	m.buffer = buffer
}

func NewParser(pt int) IParser {
	m := map[int] IParser {
		1: &Parser{},
		2: &Parser{},
		3: &Parser3{},
	}
	
	ir,ok := m[pt]
	if ok {
		return ir
	} else {
		return nil
	}
}