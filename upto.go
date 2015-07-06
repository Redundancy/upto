package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Redundancy/upto/message"
	"github.com/codegangsta/cli"
)

var app = cli.NewApp()
var extraContentPath string
var mux = http.NewServeMux()

func main() {

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
			if err = udpServer.SetDeadline(time.Now().Add(time.Second)); err != nil {
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

			if udpAddr, ok := remote.(*net.UDPAddr); ok {
				receiveUDPMessage(udpAddr.IP.String(), buffer[:n])
			} else {
				receiveUDPMessage(remote.String(), buffer[:n])
			}
		}
	}()

	// start HTTP server
	go func() {
		mux.Handle("/", http.FileServer(http.Dir(extraContentPath)))

		log.Println(
			http.ListenAndServe(
				"localhost:8080",
				mux,
			),
		)

		return
	}()

	<-waiter
}

func receiveUDPMessage(ip string, buffer []byte) {
	m := &message.UDPMessage{}
	m.UnmarshalMsg(buffer)

	if m.FillHostWithIP {
		m.Host = ip
	}
	fmt.Printf("Message: %#v\n", m)
}
