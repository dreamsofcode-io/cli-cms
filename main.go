package main

import (
	"log"

	"github.com/dreamsofcode-io/cli-cms/cmd"
)

func main() {
	log.SetFlags(0)

	if err := cmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
