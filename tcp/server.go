package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	csv "readcsv"
)

var (
	data = csv.Loadcsvfile("covid_data.csv")
)

func handleConnection(connection net.Conn) {
	fmt.Println("Handling new connection")

	defer func() {
		fmt.Println("closing connection")
		connection.Close()
	}()

	for {
		dec := json.NewDecoder(connection)
		var req csv.DataRequest
		err := dec.Decode(&req)
		if err != nil {
			fmt.Println("network error:", err)
			return
		}
		result := csv.Search(data, req.Get)

		enc := json.NewEncoder(connection)
		error := enc.Encode(&result)
		if error != nil {
			fmt.Println("failed to send response:", error)
			return

		}
	}
}

func main() {
	var host string
	var network string

	flag.StringVar(&host, "port", ":4040", "Server with port 4040")             //string flag with specified name,value and usage string.
	flag.StringVar(&network, "TCP Network", "tcp", "communication over TCP/IP") //&port is where to store the value of flag
	flag.Parse()

	ln, err := net.Listen(network, host) // listening for incoming connections
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer ln.Close() //A defer statement defers/delay the execution of a function until the surrounding function returns.
	log.Println("**** COVID-19 DETAILS OF PAKISTAN  ***")
	log.Printf("Service started: (%s) %s\n", network, host)

	for {
		connection, err := ln.Accept() // connects to the server, creates a socket
		if err != nil {
			fmt.Println(err)
			connection.Close()
			continue
		}
		log.Println("Connected to ", connection.RemoteAddr())
		go handleConnection(connection)
	}
}
