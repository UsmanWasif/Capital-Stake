package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	csv "readcsv"
	"strings"
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
		cmdline := make([]byte, (1024 * 4))
		n, err := connection.Read(cmdline)
		if n == 0 || err != nil {
			log.Println("connection read error", err)
			return
		}
		cmd, param, param2 := parsecommand(string(cmdline[0:n]))

		if cmd == "query" && param == "date" || param == "region" {
			result := csv.Search(data, param2)
			for _, lim := range result {
				_, err := connection.Write([]byte(
					fmt.Sprintf(string(lim)),
				))
				if err != nil {
					log.Println("failed to write response", err)
					return
				}
			}
		} else {
			log.Println("trimspace error")
		}
	}
}
func parsecommand(cmdline string) (cmd, param, param2 string) {
	parts := strings.Split(cmdline, ":")
	res1 := strings.TrimSpace(parts[0])
	res2 := strings.TrimSpace(parts[1])
	res3 := strings.TrimSpace(parts[2])
	cmd = strings.TrimLeft(res1, "{")
	param = strings.TrimLeft(res2, "{")
	param2 = strings.TrimRight(res3, " }} ")
	return
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
		log.Println("<Usage: {query: {region: Sindh}} OR {query: {date: 3/11/2020}} >")
		go handleConnection(connection)
	}
}
