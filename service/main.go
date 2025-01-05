package main

import (
	"flag"
	"github.com/coinbase/kryptology/service/initial"
	"log"
)

func main() {
	err := flag.Set("logtostderr", "true")
	if err != nil {
		return
	} // 将日志输出到控制台
	//flag.Parse()

	if err := initial.Execute(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}
