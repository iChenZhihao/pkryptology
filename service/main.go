package main

import (
	"flag"
	"github.com/coinbase/kryptology/service/initial"
)

func main() {
	err := flag.Set("logtostderr", "true")
	if err != nil {
		return
	} // 将日志输出到控制台
	flag.Parse()

	initial.Run()
}
