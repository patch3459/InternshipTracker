package main

import (
	"flag"
	"fmt"
	"os"
)

func usage() {
	fmt.Println("Usage")
}

func main() {
	args := os.Args

	if len(args) < 2 {
		usage()
		return
	}
	flag.Parse()

	scrapeNewInternships()

}
