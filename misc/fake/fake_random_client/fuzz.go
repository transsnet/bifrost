// +build gofuzz

package random

import (
	"crypto/tls"
	"errors"
	"flag"
	"net"
	"net/url"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/net/websocket"
)

func Fuzz(data []byte) int {
	//var addr *string = flag.String("addr", "ssl://127.0.0.1:8000", "Connection Address")
	//	var addr *string = flag.String("addr", "tcp://127.0.0.1:1883", "Connection Address")
	flag.Parse()

	//	url, _ := url.Parse(*addr)
	url, _ := url.Parse("tcp://127.0.0.1:1883")
	tlsc := NewTLSConfig()
	timeout := time.Duration(1e9)
	conn, err := openConnection(url, tlsc, timeout)
	if err != nil {
		log.Errorf("open connection failed. %#v", err.Error())
		return -1
	}
	n, err := conn.Write(data)
	if err != nil {
		log.Errorf("conn write %#v", err.Error())
		return -1
	}
	var buf []byte
	n, err = conn.Read(buf)
	if err != nil {
		log.Errorf("conn read %#v", err.Error())
		return -1
	}
	log.Infof("read data n=%d", n)
	return 0
}

func NewTLSConfig() *tls.Config {
	// Create tls.Config with desired tls properties
	return &tls.Config{
		ClientAuth: tls.NoClientCert,
		// ClientCAs = certs used to validate client cert.
		ClientCAs: nil,
		// InsecureSkipVerify = verify that cert contents
		// match server. IP matches what is in cert etc.
		InsecureSkipVerify: true,
	}
}

func openConnection(uri *url.URL, tlsc *tls.Config, timeout time.Duration) (net.Conn, error) {
	switch uri.Scheme {
	case "ws":
		conn, err := websocket.Dial(uri.String(), "mqtt", "ws://localhost")
		if err != nil {
			return nil, err
		}
		conn.PayloadType = websocket.BinaryFrame
		return conn, err
	case "wss":
		config, _ := websocket.NewConfig(uri.String(), "ws://localhost")
		config.Protocol = []string{"mqtt"}
		config.TlsConfig = tlsc
		conn, err := websocket.DialConfig(config)
		if err != nil {
			return nil, err
		}
		conn.PayloadType = websocket.BinaryFrame
		return conn, err
	case "tcp":
		conn, err := net.DialTimeout("tcp", uri.Host, timeout)
		if err != nil {
			return nil, err
		}
		return conn, nil
	case "ssl":
		fallthrough
	case "tls":
		fallthrough
	case "tcps":
		conn, err := tls.DialWithDialer(&net.Dialer{Timeout: timeout}, "tcp", uri.Host, tlsc)
		if err != nil {
			return nil, err
		}
		return conn, nil
	}
	return nil, errors.New("Unknown protocol")
}
