package main

import (
	//"encoding/hex"
	//"unsafe"
)

type Parser3 struct {
	Parser 
}

func (m *Parser3) Parse(buf []byte){
	n := len(buf)
	if n == 0 {
		return
	}
	
	//Log.Println("Parser3.Parse Bgn buffer:", hex.EncodeToString(m.buffer))

	buffer := m.buffer
	buffer = append(buffer, buf...)

	packetLen := m.packetLen
	for {
		bufLen := len(buffer)
		if bufLen < packetLen {
			break
		}

		bufChanged := false
		for i := 0; i < bufLen; i++ {
			if buffer[i] == 0x03 && i >= (packetLen-1) && buffer[i-packetLen+1] == 0x02 {
				//protocol packet must be 12bytes

				b := (i - packetLen + 1)
				//buffer[b] == 0x02
				buf := buffer[b:(b + packetLen)]
				//log.Println("buffer[b:(b + packetLen)]:", buf)
				val := buf[2:8]

				//sflag := string(val)
				//log.Println("val[2:8]:", val)
				//dot must be 0-6
				dot := buf[8] - 0x30
				//log.Println("dot:", buf[8]," ",dot, 6-dot)
				dot = 6 - dot
				//must use copy..
				tmp := make([]byte, len(val[dot:]))
				copy(tmp, val[dot:])
				//log.Println("tmp:", tmp)
				val = val[:dot]
				//log.Println("val1:", val)
				val = append(val, '.')
				//log.Println("val2:", val, "tmp:", tmp)
				val = append(val, tmp...)
				//log.Println("val3:", val)

				sval := string(val)
				v, _ := m.toFloat(sval)
				m.valChannel <- v
				buffer = buffer[(i + 1):]
				bufChanged = true
				break
			}
		}

		if !bufChanged {
			break
		}

	}
	
	m.buffer = buffer
	//Log.Println("Parser3.Parse End buffer:", hex.EncodeToString(m.buffer))
}
