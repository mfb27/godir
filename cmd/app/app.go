package main

import (
	"flag"
	"fmt"
)

var c = flag.String("c", "", "要打印的内容")

func main() {
	flag.Parse()
	fmt.Println(*c)
}
