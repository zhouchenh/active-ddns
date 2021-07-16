package main

import "flag"

var (
	serverListenAddr  = flag.String("s", "", "Run as a server and listen at the specific address")
	clientConnectAddr = flag.String("c", "", "Run as a client and connect to the specific address")
	script            = flag.String("script", "", "Specify the script to be executed when the IP address is updated")
	keyword           = flag.String("keyword", "{}", "Specify the keyword in the script to be replaced by the updated IP address")
	shellArgs         = flag.String("shell", "", "Specify the shell and arguments which is used to run the DDNS script")
	certFilePath      = flag.String("cert", "", "Specify the path to the certificate file")
	keyFilePath       = flag.String("key", "", "Specify the path to the private key file")
	noTLS             = flag.Bool("notls", false, "Do not use TLS")
	insecureTLS       = flag.Bool("insecuretls", false, "Allow insecure TLS")
	hbiValue          = flag.Int("hbi", 5000, "Specify the interval between heartbeats in milliseconds")
	mhbValue          = flag.Int("mhb", 3, "Specify the number of missed heartbeats allowed before disconnection")
	minRI             = flag.Int("minri", 1000, "Specify the minimal interval between reconnections in milliseconds")
	maxRI             = flag.Int("maxri", 15000, "Specify the maximal interval between reconnections in milliseconds")
	logLevelStr       = flag.String("log", "info", "Specify the log level { debug | info | warning | error | off }")
	logTime           = flag.Bool("logtime", false, "Output logs with timestamps")
	version           = flag.Bool("version", false, "Print version information and exit")
)
