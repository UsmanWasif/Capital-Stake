package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"os"
	csv "readcsv"
)

func main() {
	var host string
	var network string

	flag.StringVar(&host, "port", ":4040", "Server with port 4040")             //string flag with specified name,value and usage string.
	flag.StringVar(&network, "TCP Network", "tcp", "communication over TCP/IP") //&port is where to store the value of flag
	flag.Parse()

	connection, err := net.Dial(network, host)
	if err != nil {
		fmt.Println("failed to create socket", err)
		os.Exit(1) //os.exit(1) means there is some error that is why program is exitting
	} else if connection == nil {
		fmt.Println("failed to create a connection successfully")
		os.Exit(1)
	}
	defer connection.Close()
	var param string
	for {
		_, err = fmt.Scanf("%s", &param)
		if err != nil {
			fmt.Println("Usage: <get sindh or get date>")
			continue
		}
		req := csv.DataRequest{Get: param}
		// Send request:
		// use json encoder to encode value of type csv.DataRequest
		// and stream it to the server via net.Conn.
		if err := json.NewEncoder(connection).Encode(&req); err != nil {
			switch err := err.(type) {
			case net.Error:
				fmt.Println("failed to send request:", err)
				os.Exit(1)
			default:
				fmt.Println("failed to encode request:", err)
				continue
			}
		}
		// Display response
		var covidDetails []byte
		err = json.NewDecoder(connection).Decode(&covidDetails)
		if err != nil {
			switch err := err.(type) {
			case net.Error:
				fmt.Println("failed to receive response:", err)
				os.Exit(1)
			default:
				fmt.Println("failed to decode response:", err)
				continue
			}
		}
		fmt.Println(string(covidDetails))
	}
}
