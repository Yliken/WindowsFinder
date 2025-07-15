package main

/*
#include <stdio.h>

// C 语言函数
int add(int a, int b) {
    return a + b;
}
*/
import "C"
import "fmt"

func main() {
	a := 5
	b := 7

	// 调用 C 语言的 add 函数
	result := C.add(C.int(a), C.int(b))
	fmt.Printf("C.add(%d, %d) = %d\n", a, b, int(result))
}
