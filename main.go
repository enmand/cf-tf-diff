package main

import (
	"log"

	"github.com/enmand/cf-tf-diff/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
