package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Redundancy/upto/store"
	"github.com/codegangsta/cli"
)

var app = cli.NewApp()
var extraContentPath string = "./content/"
var mux = http.NewServeMux()

var datastore store.UptoDataStore

func main() {

	var handler = &BinaryMessageHandler{}
	datastore = &SimpleMemoryStore{}

	// start UDP server
	serverAddr, _ := net.ResolveUDPAddr("udp", ":8123")
	udpServer, err := net.ListenUDP("udp", serverAddr)

	if err != nil {
		fmt.Fprint(os.Stderr, err)
		return
	}

	waiter := make(chan bool)

	go func() {
		defer close(waiter)
		interrupt := make(chan os.Signal, 1)
		signal.Notify(interrupt, os.Interrupt)

		buffer := make([]byte, 1024)

	Messages:
		for {
			oneSecondFromNow := time.Now().Add(time.Second)
			if err = udpServer.SetDeadline(oneSecondFromNow); err != nil {
				fmt.Fprint(os.Stderr, err)
				return
			}

			select {
			case <-interrupt:
				return
			default:
				// don't wait
			}

			n, remote, err := udpServer.ReadFrom(buffer)

			if err != nil {
				switch e := err.(type) {
				case *net.OpError:
					if !e.Timeout() {
						fmt.Fprint(os.Stderr, err)
						return
					}
					continue Messages
				default:
					fmt.Fprint(os.Stderr, err)
					return
				}
			}

			var ip string

			switch t := remote.(type) {
			case *net.UDPAddr:
				ip = t.IP.String()
			default:
				ip = t.String()
			}

			handler.ReceiveMessage(datastore, ip, buffer[:n])
		}
	}()

	// start HTTP server
	go func() {
		d := http.Dir(extraContentPath)
		fs := http.FileServer(d)
		log.Println("Serving content from", extraContentPath)
		mux.Handle("/", fs)

		log.Println(
			http.ListenAndServe(
				":8080",
				mux,
			),
		)

		return
	}()

	<-waiter
}
