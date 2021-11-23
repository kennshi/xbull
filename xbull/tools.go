package main
import (
	"fmt"
	"sort"
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


func sortMap(){
	m := map[string]int {
		"something" :10,
		"yo":		20,
		"blah":		30,
	}

	type kv struct {
		Key string
		Value int
	}

	ss := []kv{}

	for k,v := range m {
		ss = append(ss, kv{k,v})
	}

	sort.Slice(ss, func(i,j int) bool {
		return ss[i].Value > ss[j].Value //降序
		//return ss[i].Value < ss[j].Value //升序
	})

	for _,kv := range ss {
		fmt.Printf("%s, %d\n", kv.Key, kv.Value)
	}

}

func array(){

	{
		var a [3]int 
		a[0] = 12
		a[1] = 78
		a[2] = 50
		fmt.Println(a)
	}

	{
		a:= [3]int{12,78,50}
		fmt.Println(a)
	}

	{
		a := [...]int{12,78,50}
		fmt.Println(a)
	}

	//数组长度是数组类型的一部分 因此 [5]int 和  [25]int 是两个不同类型的数组。正因为如此， 一个数组不能动态改变长度。 slices （切片）可以弥补这个不足
	//Go语言中， 数组是值类型， 意味着数组变量被赋值时，将会获得原先数组的拷贝。。。。
	//range for 遍历数组

	{
		a := [5]int{76, 77, 78, 79, 80}
   		var b []int = a[1:4] //creates a slice from a[1] to a[3]
    	fmt.Println(b)
		//通过 a[start:end] 这样的语法创建了一个从 a[start] 到 a[end -1] 的切片。

		c := []int{6, 7, 8} //creates and array and returns a slice reference
		fmt.Println(c)
		//切片本身不包含任何数据。它仅仅是底层数组的一个上层表示。对切片进行的任何修改都将反映在底层数组中。
	}

	//map 是引用类型








}

func init() {
	//Log.Println("Here is tools.go", VERSION)
}

