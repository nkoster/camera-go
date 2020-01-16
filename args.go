package main

import (
	"fmt"
	"os"
	"strings"
)

func args() (string, string, string) {
	var serverHost string = "localhost"
	var serverPort string = "8086"
	var serverHTML string = "./"
	for i, arg := range os.Args {
		if arg == "-host" {
			if i+1 < len(os.Args) {
				serverHost = os.Args[i+1]
			}
			arg = ""
		}
		if arg == "-port" {
			if i+1 < len(os.Args) {
				serverPort = os.Args[i+1]
			}
			arg = ""
		}
		if arg == "-html" {
			if i+1 < len(os.Args) {
				serverHTML = os.Args[i+1]
				serverHTML = strings.TrimRight(serverHTML, "/")
			}
			arg = ""
		}
		if arg == "--help" || arg == "-help" || arg == "/h" {
			fmt.Printf("usage: %s [[-host <host>] [-port <port>] [-html <path>]]\n", os.Args[0])
			os.Exit(0)
		}
	}
	return serverHost, serverPort, serverHTML
}
