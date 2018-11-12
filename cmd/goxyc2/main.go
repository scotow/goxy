package main

import (
	"flag"
	goxy "github.com/scotow/goxy/client2"
	log "github.com/sirupsen/logrus"
	"io"
	"net"
)

var (
	localTCP   = flag.String("l", ":2222", "local listening TCP address (host:port)")
	remoteHTTP = flag.String("r", "localhost:8080", "remote Goxy server HTTP address (host:port)")
)

func main() {
	flag.Parse()

	localTCPAddr, err := net.ResolveTCPAddr("tcp", *localTCP)
	if err != nil {
		log.WithFields(log.Fields{
			"address": *localTCP,
			"error":   err,
		}).Error("Invalid listening address.")
	}

	remoteHTTPAddr, err := net.ResolveTCPAddr("tcp", *remoteHTTP)
	if err != nil {
		log.WithFields(log.Fields{
			"address": *remoteHTTP,
			"error":   err,
		}).Error("Invalid remote HTTP address.")
	}

	listener, err := net.ListenTCP("tcp", localTCPAddr)
	if err != nil {
		log.WithFields(log.Fields{
			"address": *localTCP,
			"error":   err,
		}).Error("Cannot start TCP listener.")
	}

	log.WithFields(log.Fields{
		"local":  *localTCP,
		"remote": *remoteHTTP,
	}).Info("Goxy client started.")

	for {
		tcpConn, err := listener.Accept()
		if err != nil {
			log.WithFields(log.Fields{
				"address": *localTCP,
				"error":   err,
			}).Fatal("Cannot accept TCP connection.")
		}

		log.WithFields(log.Fields{
			"local":  tcpConn.LocalAddr(),
			"remote": tcpConn.RemoteAddr(),
		}).Info("TCP connection accepted.")

		goxyConn, err := goxy.Dial(remoteHTTPAddr)
		if err != nil {
			log.WithFields(log.Fields{
				"address": *remoteHTTP,
				"error":   err,
			}).Error("Cannot open Goxy connection.")
			continue
		}

		log.WithFields(log.Fields{
			"address": *remoteHTTP,
		}).Info("Goxy connection created.")

		go io.Copy(tcpConn, goxyConn)
		go io.Copy(goxyConn, tcpConn)
	}
}
