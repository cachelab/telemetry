//+build !test

package main

import (
	"log"
	"os"
)

const version = "1.1.1"

func main() {
	var svc Service

	err := svc.Init()
	if err != nil {
		log.Fatal(err.Error())
	}

	if err != nil {
		os.Exit(-1)
	}
}
