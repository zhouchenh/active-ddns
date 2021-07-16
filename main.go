package main

import (
	"flag"
	"fmt"
	"github.com/zhouchenh/active-ddns/client"
	"github.com/zhouchenh/active-ddns/doublable"
	"github.com/zhouchenh/active-ddns/info"
	"github.com/zhouchenh/active-ddns/logger"
	"github.com/zhouchenh/active-ddns/server"
	"github.com/zhouchenh/active-ddns/shell"
	"net"
	"os"
	"strings"
	"time"
)

func main() {
	flag.Parse()
	if *version {
		printVersion()
		return
	}
	if *serverListenAddr != "" && *clientConnectAddr != "" {
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "flag -s and -c cannot be set together\n")
		flag.Usage()
		os.Exit(2)
	}
	if *hbiValue <= 0 {
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "invalid value \"%d\" for flag -hbi: value out of range\n", *hbiValue)
		flag.Usage()
		os.Exit(2)
	}
	if *mhbValue <= 0 {
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "invalid value \"%d\" for flag -mhb: value out of range\n", *mhbValue)
		flag.Usage()
		os.Exit(2)
	}
	logger.SetTimestamp(*logTime)
	logger.SetLogLevel(logLevel())
	if *shellArgs != "" {
		shell.Shell = *shellArgs
	}
	if *serverListenAddr != "" {
		if !*noTLS && (*certFilePath == "" || *keyFilePath == "") {
			_, _ = fmt.Fprintf(flag.CommandLine.Output(), "a valid certificate and private key should be specified with -cert and -key\n")
			flag.Usage()
			os.Exit(2)
		}
		runServer()
	} else if *clientConnectAddr != "" {
		if *minRI < 0 {
			_, _ = fmt.Fprintf(flag.CommandLine.Output(), "invalid value \"%d\" for flag -minri: value out of range\n", *minRI)
			flag.Usage()
			os.Exit(2)
		}
		if *maxRI < 0 {
			_, _ = fmt.Fprintf(flag.CommandLine.Output(), "invalid value \"%d\" for flag -maxri: value out of range\n", *maxRI)
			flag.Usage()
			os.Exit(2)
		}
		if *script == "" {
			_, _ = fmt.Fprintf(flag.CommandLine.Output(), "a script should be specified with -script\n")
			flag.Usage()
			os.Exit(2)
		}
		if *keyword == "" {
			_, _ = fmt.Fprintf(flag.CommandLine.Output(), "a non-empty keyword should be specified with -keyword\n")
			flag.Usage()
			os.Exit(2)
		}
		runClient()
	} else {
		flag.Usage()
	}
}

func printVersion() {
	for _, s := range info.VersionStatement() {
		_, _ = fmt.Fprintln(logger.Output(), s)
	}
}

func logLevel() logger.Level {
	switch *logLevelStr {
	case "debug":
		return logger.DebugLevel
	case "info":
		return logger.InfoLevel
	case "warning":
		return logger.WarningLevel
	case "error":
		return logger.ErrorLevel
	case "off":
		return logger.Disabled
	default:
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "invalid value \"%s\" for flag -log: undefined log level\n", *logLevelStr)
		flag.Usage()
		os.Exit(2)
	}
	return 0
}

func runServer() {
	s := &server.Server{
		ListenAddr:              *serverListenAddr,
		NoTLS:                   *noTLS,
		CertFile:                *certFilePath,
		KeyFile:                 *keyFilePath,
		HeartbeatInterval:       time.Duration(*hbiValue) * time.Millisecond,
		MissedHeartbeatsAllowed: *mhbValue,
	}
	printVersion()
	logger.Fatal().Msg(s.Run().Error())
}

func runClient() {
	c := &client.Client{
		ConnectAddr:             *clientConnectAddr,
		NoTLS:                   *noTLS,
		AllowInsecureTLS:        *insecureTLS,
		HeartbeatInterval:       time.Duration(*hbiValue) * time.Millisecond,
		MissedHeartbeatsAllowed: *mhbValue,
		RedialInterval:          &doublable.Duration{Min: time.Duration(*minRI) * time.Millisecond, Max: time.Duration(*maxRI) * time.Millisecond},
		OnIPAddrUpdate:          onIPAddrUpdate,
	}
	printVersion()
	logger.Fatal().Msg(c.Run().Error())
}

func onIPAddrUpdate(newIPAddr net.IP) {
	shell.Script(strings.ReplaceAll(*script, *keyword, newIPAddr.String())).Run()
}
