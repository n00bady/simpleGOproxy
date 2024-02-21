package main

import (
	"flag"
	"log"
)

func main() {
	log.Println("\t\t>>>> SimpleGOProxy >>>>\t\t")

	reverse := flag.Bool("reverse", false, "Use this option to launch the reverse proxy.")

	flag.Parse()

	if *reverse {
		log.Println("Staring the test http server...")
		StartReverseProxy()
	} else {
		log.Println("Starting proxy...")
		StartProxy()
	}
}
