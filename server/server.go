package server

import (
	"crypto/tls"
	"github.com/zhouchenh/active-ddns/logger"
	"github.com/zhouchenh/active-ddns/neterr"
	"github.com/zhouchenh/active-ddns/ticker"
	"io"
	"io/ioutil"
	"net"
	"time"
)

type Server struct {
	ListenAddr              string
	NoTLS                   bool
	CertFile                string
	KeyFile                 string
	HeartbeatInterval       time.Duration
	MissedHeartbeatsAllowed int
	idleTimeout             time.Duration
}

func (s *Server) Run() (err error) {
	s.idleTimeout = s.HeartbeatInterval/2 + s.HeartbeatInterval + time.Duration(s.MissedHeartbeatsAllowed)*s.HeartbeatInterval
	var listen func() (net.Listener, error)
	if s.NoTLS {
		listen = func() (net.Listener, error) {
			return net.Listen("tcp", s.ListenAddr)
		}
	} else {
		var cert tls.Certificate
		cert, err = tls.LoadX509KeyPair(s.CertFile, s.KeyFile)
		if err != nil {
			return err
		}
		listen = func() (net.Listener, error) {
			return tls.Listen("tcp", s.ListenAddr, &tls.Config{Certificates: []tls.Certificate{cert}})
		}
	}
	var listener net.Listener
	listener, err = listen()
	if err != nil {
		return
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			neterr.LogError(err)
			continue
		}
		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	defer conn.Close()
	remoteAddr := conn.RemoteAddr().String()
	defer logger.Info().Str("client", remoteAddr).Msg("Disconnected")
	logger.Info().Str("client", remoteAddr).Msg("Connected")
	tcpAddr, ok := conn.RemoteAddr().(*net.TCPAddr)
	if !ok {
		return
	}
	ip := tcpAddr.IP
	if ipv4 := ip.To4(); ipv4 != nil {
		ip = ipv4
	}
	data := append([]byte{byte(len(ip))}, []byte(ip)...)
	length := len(data)
	for written := 0; written < length; {
		err := conn.SetWriteDeadline(time.Now().Add(s.idleTimeout))
		if err != nil {
			neterr.LogError(err)
			return
		}
		n, err := conn.Write(data[written:])
		if err != nil {
			neterr.LogError(err)
			return
		}
		written += n
	}
	logger.Debug().Str("client", remoteAddr).Msg("Sent IP address")
	t := ticker.NewTicker(s.HeartbeatInterval)
	defer t.Stop()
	go s.sendHeartbeats(conn, t, remoteAddr)
	s.receiveHeartbeats(conn, remoteAddr)
}

func (s *Server) sendHeartbeats(conn net.Conn, t *ticker.Ticker, remoteAddr string) {
	logger.Debug().Str("client", remoteAddr).Msg("Heartbeat started")
	for {
		select {
		case <-t.Done:
			logger.Debug().Str("client", remoteAddr).Msg("Heartbeat stopped")
			return
		case <-t.C:
			err := conn.SetWriteDeadline(time.Now().Add(s.idleTimeout))
			if err != nil {
				neterr.LogError(err)
				conn.Close()
				return
			}
			_, err = conn.Write([]byte("\x00"))
			if err != nil {
				neterr.LogError(err)
				conn.Close()
				return
			}
			logger.Debug().Str("client", remoteAddr).Msg("Sent Heartbeat")
		}
	}
}

func (s *Server) receiveHeartbeats(conn net.Conn, remoteAddr string) {
	buffer := make([]byte, 1)
	for {
		err := conn.SetReadDeadline(time.Now().Add(s.idleTimeout))
		if err != nil {
			neterr.LogError(err)
			return
		}
		_, err = conn.Read(buffer)
		if err != nil {
			neterr.LogError(err)
			return
		}
		if buffer[0] == 0 {
			logger.Debug().Str("client", remoteAddr).Msg("Received Heartbeat")
			continue
		}
		err = conn.SetReadDeadline(time.Now().Add(s.idleTimeout))
		if err != nil {
			neterr.LogError(err)
			return
		}
		_, err = io.CopyN(ioutil.Discard, conn, int64(buffer[0]))
		if err != nil {
			neterr.LogError(err)
			return
		}
		logger.Warning().Str("client", remoteAddr).Int("length", int(buffer[0])).Msg("Received invalid data")
	}
}
