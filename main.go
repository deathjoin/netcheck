package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/go-ping/ping"
)

var (
	serverHost      string
	serverPort      int
	checkInterval   int
	checkTimeout    int
	printOnlyErrors bool
	errorLogger     *log.Logger
	defaultLogger   *log.Logger
)

func init() {
	// script arguments
	flag.StringVar(&serverHost, "host", "", "Server ip or name to check.")
	flag.IntVar(&serverPort, "port", 0, "Server TCP port to check.")
	flag.IntVar(&checkInterval, "interval", 1, "Check interval in seconds.")
	flag.IntVar(&checkTimeout, "timeout", 5, "Connection timeout in seconds.")
	flag.BoolVar(&printOnlyErrors, "only_errors", false, "Print only fails.")
	flag.Parse()

	// loggers
	errorLogger = log.New(os.Stderr, "ERROR: ", log.LstdFlags)
	defaultLogger = log.New(os.Stdout, "", log.LstdFlags)
	if printOnlyErrors {
		defaultLogger.SetOutput(ioutil.Discard)
	}

	// check for argument values
	if serverHost == "" {
		errorLogger.Fatalln("Server host not set.")
	}

	if serverPort == 0 {
		errorLogger.Fatalln("Server port not set.")
	}

	// Test for sudo privileges
	pinger, _ := ping.NewPinger(serverHost)
	pinger.SetPrivileged(true)
	pinger.Count = 1
	pinger.Timeout = time.Second * time.Duration(checkTimeout)
	pingErr := pinger.Run()
	if pingErr != nil {
		errorLogger.Println("Cannot send ICMP without SUDO.")
	}
}

func main() {
	serverAddress := net.JoinHostPort(serverHost, strconv.Itoa(serverPort))
	timeout := time.Second * time.Duration(checkTimeout)

	defaultLogger.Printf("Starting tcp port check: %s\n", serverAddress)
	for {
		tcpConn, tcpErr := net.DialTimeout("tcp", serverAddress, timeout)

		pinger, _ := ping.NewPinger(serverHost)
		pinger.SetPrivileged(true)
		pinger.Count = 1
		pinger.Timeout = timeout
		pingErr := pinger.Run()

		pingResult := "FAIL"
		if pingErr != nil {
			pingResult = "UNKNOWN"
		} else {
			if pinger.Statistics().PacketLoss == 0 {
				pingResult = "OK"
			}
		}

		tcpResult := "FAIL"
		if tcpErr == nil {
			tcpResult = "OK"
		}

		// fmt.Printf("Sent: %d RECV: %d Loss: %v\n", pinger.Statistics().PacketsSent, pinger.Statistics().PacketsRecv, pinger.Statistics().PacketLoss)

		if pingResult == "OK" && tcpResult == "OK" {
			tcpConn.Close()
			defaultLogger.Printf("Connection success to \"%s\"\n", serverAddress)
		} else {
			errorLogger.Printf(
				"Connection fail to \"%s\" within %d seconds. PING: %s, TCP_PORT: %s\n",
				serverAddress,
				checkTimeout,
				pingResult,
				tcpResult,
			)
		}
		time.Sleep(time.Second * time.Duration(checkInterval))
	}
}
