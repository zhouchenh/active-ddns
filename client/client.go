package client

import (
	"crypto/tls"
	"github.com/zhouchenh/active-ddns/doublable"
	"github.com/zhouchenh/active-ddns/logger"
	"github.com/zhouchenh/active-ddns/neterr"
	"github.com/zhouchenh/active-ddns/ticker"
	"io"
	"io/ioutil"
	"net"
	"time"
)

type Client struct {
	ConnectAddr             string
	NoTLS                   bool
	AllowInsecureTLS        bool
	HeartbeatInterval       time.Duration
	MissedHeartbeatsAllowed int
	idleTimeout             time.Duration
	RedialInterval          *doublable.Duration
	currentIPAddr           net.IP
	OnIPAddrUpdate          func(newIPAddr net.IP)
}

func (c *Client) Run() (err error) {
	c.idleTimeout = c.HeartbeatInterval/2 + c.HeartbeatInterval + time.Duration(c.MissedHeartbeatsAllowed)*c.HeartbeatInterval
	_, err = net.ResolveTCPAddr("tcp", c.ConnectAddr)
	if err != nil {
		return err
	}
	var dial func() (net.Conn, error)
	if c.NoTLS {
		dial = func() (net.Conn, error) {
			return net.Dial("tcp", c.ConnectAddr)
		}
	} else {
		config := &tls.Config{InsecureSkipVerify: c.AllowInsecureTLS}
		dial = func() (net.Conn, error) {
			return tls.Dial("tcp", c.ConnectAddr, config)
		}
	}
	for {
		conn, err := dial()
		if err != nil {
			neterr.LogError(err)
			c.RedialInterval.Double()
			ri := c.RedialInterval.Duration()
			logger.Info().Str("duration", ri.String()).Msg("Waiting for reconnection")
			time.Sleep(ri)
			continue
		}
		c.RedialInterval.Minimize()
		c.handleConn(conn)
	}
}

func (c *Client) handleConn(conn net.Conn) {
	defer conn.Close()
	remoteAddr := conn.RemoteAddr().String()
	defer logger.Info().Str("server", remoteAddr).Msg("Disconnected")
	logger.Info().Str("server", remoteAddr).Msg("Connected")
	buffer := make([]byte, net.IPv6len)
	err := conn.SetReadDeadline(time.Now().Add(c.idleTimeout))
	if err != nil {
		neterr.LogError(err)
		return
	}
	_, err = conn.Read(buffer[:1])
	if err != nil {
		neterr.LogError(err)
		return
	}
	if buffer[0] != net.IPv4len && buffer[0] != net.IPv6len {
		logger.Warning().Str("server", remoteAddr).Int("length", int(buffer[0])).Msg("Received invalid data")
		return
	}
	length := int(buffer[0])
	for read := 0; read < length; {
		err := conn.SetReadDeadline(time.Now().Add(c.idleTimeout))
		if err != nil {
			neterr.LogError(err)
			return
		}
		n, err := conn.Read(buffer[read:length])
		if err != nil {
			neterr.LogError(err)
			return
		}
		read += n
	}
	ip := net.IP(buffer[:length])
	logger.Debug().Str("server", remoteAddr).Str("address", ip.String()).Msg("Received IP address")
	go c.onIPAddrReceived(ip)
	t := ticker.NewTicker(c.HeartbeatInterval)
	defer t.Stop()
	go c.sendHeartbeats(conn, t, remoteAddr)
	c.receiveHeartbeats(conn, remoteAddr)
}

func (c *Client) onIPAddrReceived(ip net.IP) {
	if ip.Equal(c.currentIPAddr) {
		return
	}
	c.currentIPAddr = ip
	c.OnIPAddrUpdate(ip)
}

func (c *Client) sendHeartbeats(conn net.Conn, t *ticker.Ticker, remoteAddr string) {
	logger.Debug().Str("server", remoteAddr).Msg("Heartbeat started")
	for {
		select {
		case <-t.Done:
			logger.Debug().Str("server", remoteAddr).Msg("Heartbeat stopped")
			return
		case <-t.C:
			err := conn.SetWriteDeadline(time.Now().Add(c.idleTimeout))
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
			logger.Debug().Str("server", remoteAddr).Msg("Sent Heartbeat")
		}
	}
}

func (c *Client) receiveHeartbeats(conn net.Conn, remoteAddr string) {
	buffer := make([]byte, 1)
	for {
		err := conn.SetReadDeadline(time.Now().Add(c.idleTimeout))
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
			logger.Debug().Str("server", remoteAddr).Msg("Received Heartbeat")
			continue
		}
		err = conn.SetReadDeadline(time.Now().Add(c.idleTimeout))
		if err != nil {
			neterr.LogError(err)
			return
		}
		_, err = io.CopyN(ioutil.Discard, conn, int64(buffer[0]))
		if err != nil {
			neterr.LogError(err)
			return
		}
		logger.Warning().Str("server", remoteAddr).Int("length", int(buffer[0])).Msg("Received invalid data")
	}
}
