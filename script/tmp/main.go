//package main
//import "C"
//import (
//	"fmt"
//	"unsafe"
//)

// import (
//
//	"fmt"
//	"unsafe"
//
// )
//
//	func main() {
//		//const s = 2.2131e-123
//		//j := 3
//		//x := reflect.ValueOf(&j)
//		//fmt.Println(x.CanAddr())
//		//y := x.Elem()
//		//fmt.Println(y.Kind(), x.Kind(), y.CanAddr())
//		//y.Set(reflect.ValueOf(2))
//		//fmt.Println(j)
//		var f string = "haha"
//		pointer := unsafe.Pointer(&f)
//		fmt.Println(*(*string)(pointer))
//
// }
package main

//
// 引用的C头文件需要在注释中声明，紧接着注释需要有import "C"，且这一行和注释之间不能有空格
//

/*
   #include <stdio.h>
   #include <stdlib.h>
   #include <unistd.h>

   void myprint(char* s) {
   	printf("%s\n", s);
   }

*/
import "C"

import (
	"fmt"
	"unsafe"
)

func main() {
	//使用C.CString创建的字符串需要手动释放。
	cs := C.CString("Hello World\n")
	C.myprint(cs)
	C.free(unsafe.Pointer(cs))
	fmt.Println("call C.sleep for 3s")
	C.sleep(3)
	return
}
