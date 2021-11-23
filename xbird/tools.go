package main
import (
	"fmt"
	"strings"
	"math"
)

func sayHello()  {
	fmt.Println("hello world")
	//what exactly does he mean by that

}

func invertString(s string) string {
	bytes := [] byte(s)
	for from , to :=0, len(bytes)-1; from < to;from,to = from+1, to-1{
		bytes[from], bytes[to] = bytes[to], bytes[from]
	}
	s = string(bytes)
	return s
}


func fmatString1( s string ) string {
	fmt.Println(" fmatString  < ", s , " > len= < " , len(s) ," > ")
	s = invertString(s)
	fmt.Println(" fmatString  < ", s , " > len= < " , len(s) ," > ")
	return  s + strings.Repeat("0", 7- len(s)) + "="
}


func fmatString( f float64 , protNumb int) string {

	f = math.Round(f*1000) / 1000
	s := fmt.Sprintf("%.3f", f)
	if protNumb == 1 {
		if len(s) > 7 {
			s = s[:7]
		}
		s = "=" + strings.Repeat("0", 7- len(s)) + s
		return invertString(s)
	} else if protNumb == 2 {
		if len(s) > 8{
			s = s[:8]
		}
		s = "=" + strings.Repeat("0", 8- len(s)) + s
		return invertString(s)
	} else if protNumb == 3 {
		if len(s) > 7 {
			s = s[:7]
		}
		s = strings.Repeat("0", 7- len(s)) + s
		buf := []byte {}
		buf = append(buf,0x02,'+')
		b := []byte(s)
		for _,c := range b {
			if c != '.'{
				buf = append(buf,c)
			}
		}

		dotpos := strings.Index(s,".")
		dotpos = 6-dotpos
		dotpos = dotpos + 0x30
		buf = append(buf,byte(dotpos))

		xor := byte(0)

		for i:=1 ;i<=8;i++ {
			xor = xor ^ buf[i]
		}

		h := (xor & 0xf0) >> 4
		l := (xor & 0x0f)

		buf = append(buf,h+0x30)
		buf = append(buf,l+0x30)
		buf = append(buf,0x03)
		s = string(buf)
		//fmt.Println("---",s,"--", []byte(s),"---")
	}
	//fmt.Println(" fmatString  < ", s , " > len= < " , len(s) ," > ")
	return  s
}
