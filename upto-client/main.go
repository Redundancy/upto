package main

import (
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/Redundancy/upto/message"
	"github.com/codegangsta/cli"
)

var app = cli.NewApp()

func parseHosts(c *cli.Context) (destination, hostname string, err error) {
	remoteHost := c.GlobalString("remotehost")
	port := c.GlobalInt("remoteport")

	destination = fmt.Sprintf("%v:%v", remoteHost, port)
	hostname, err = os.Hostname()

	if !c.GlobalBool("autohost") {
		hostname = c.GlobalString("hostname")
	}

	return
}

// Send a state to
func SendState(c *cli.Context, messageType message.MessageType) {
	destination, hostname, err := parseHosts(c)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing options: %v", err)
		os.Exit(1)
	}

	if len(c.Args()) < 2 {
		fmt.Fprintf(os.Stderr, "Must provide an event name: %v", c.Args())
		os.Exit(1)
	}

	context := c.Args()[0]
	eventName := c.Args()[1]

	m := &message.UDPMessage{
		Context: context,
		Name:    strings.Split(eventName, "."),
		Type:    messageType,
		Time:    time.Now(),
		Host:    hostname,
	}

	err = sendmessage(m, destination)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not send message: %v", err)
		os.Exit(1)
	}
}

func CreateTimelineAndContext(c *cli.Context) {
	destination, _, err := parseHosts(c)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing options: %v", err)
		os.Exit(1)
	}

	if len(c.Args()) < 1 {
		fmt.Fprintf(os.Stderr, "Must provide a context name: %v", c.Args())
		os.Exit(1)
	}

	context := c.Args()[0]

	m := &message.UDPMessage{
		Context: context,
		Type:    message.MessageNewTimeline,
	}

	err = sendmessage(m, destination)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not send message: %v", err)
		os.Exit(1)
	}
}

func sendmessage(m *message.UDPMessage, destination string) error {
	con, err := net.Dial("udp", destination)

	if err != nil {
		return err
	}

	defer con.Close()

	b, e := m.MarshalMsg(nil)

	if e != nil {
		return e
	}

	_, err = con.Write(b)

	return err
}

func main() {
	app.Name = "upto-client"
	app.Author = "Daniel Speed"
	app.Commands = []cli.Command{
		{
			Name:    "new",
			Aliases: []string{"n"},
			Usage: "Create a new context instance. Usually done before a deployment etc." +
				" Will also create the context if needed.",
			Action: func(c *cli.Context) {
				println("added task: ", c.Args().First())
			},
		},
		{
			Name:    "start",
			Aliases: []string{"s"},
			Usage:   "send a start event",
			Action: func(c *cli.Context) {
				SendState(c, message.MessageStartEvent)
			},
		},
		{
			Name:    "end",
			Aliases: []string{"e"},
			Usage:   "send an end event",
			Action: func(c *cli.Context) {
				SendState(c, message.MessageEndEvent)
			},
		},
		{
			Name:    "event",
			Aliases: []string{"v"},
			Usage:   "send an event",
			Action: func(c *cli.Context) {
				SendState(c, message.MessageSingleEvent)
			},
		},
	}

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "remotehost, rh",
			Value: "localhost",
			Usage: "The server to which to send the message",
		},
		cli.IntFlag{
			Name:  "remoteport, rp",
			Value: 8123,
			Usage: "The remote server port to send to",
		},
		cli.StringFlag{
			Name:  "hostname, n",
			Value: "",
			Usage: "The host name for the event",
		},
		cli.BoolFlag{
			Name:  "autohost, a",
			Usage: "determine and fill in the host field automatically",
		},
	}

	app.Usage = `
    Send events to an upto server / agent.
    Syntax:
        upto-client [options] start context parentEvent.eventName
    `

	app.Run(os.Args)
}
