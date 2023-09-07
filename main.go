package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/patch3459/InternshipTracker/parse"
)

func usage() {
	fmt.Println("Usage")
}

func main() {
	args := os.Args

	if len(args) < 2 {
		fmt.Println("")
	}
	flag.Parse()

	parse.ScrapeNewInternships()

}
