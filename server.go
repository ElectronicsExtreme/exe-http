package exeserver

import (
	"crypto/tls"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
)

const ()

var ()

type ServerConfig struct {
	ServeMux        *http.ServeMux
	Address         string
	TlsEnable       bool
	TlsCertificates []tls.Certificate
}

func ListenAndServe(configs []ServerConfig) uint {
	var status uint = 0
	var serversWg sync.WaitGroup
	signalStopFlag := false

	// listen for termination signals
	termsig := make(chan os.Signal, 1)
	signal.Notify(termsig, os.Interrupt)

	if configs == nil {
		log.Println("")
		signal.Stop(termsig)
		close(termsig)
		return 1
	}

	// listen requests
	listeners, terminateChannel, err := startServer(configs)
	if err != nil {
		signal.Stop(termsig)
		log.Println(err)
		status |= uint(1)
		if len(listeners) == 0 {
			log.Println("no server has been started")
			close(terminateChannel)
			close(termsig)
			return status
		} else {
			log.Println("stopping activated servers")
		}
	} else {
		log.Println("all requests listeners started")
	}
	serversWg.Add(len(listeners))
	allTerminated := make(chan int, 1)
	go func() {
		serversWg.Wait()
		allTerminated <- 1
	}()
	for {
		select {
		case <-termsig:
			if !signalStopFlag {
				signal.Stop(termsig)
				signalStopFlag = true
				log.Println("received termination request. stopping")
				for serverId := range listeners {
					if listeners[serverId] == nil {
						continue
					}
					if err := listeners[serverId].Close(); err != nil {
						log.Printf("server #%v: failed to close requests listener: %v\n", serverId, err)
						serversWg.Done()
						status |= uint(2)
					}
					listeners[serverId] = nil
				}
			}
		case terminate := <-terminateChannel:
			log.Printf("server #%v: requests listener stopped with status %v\n", terminate.serverId, terminate.status)
			status |= terminate.status
			if listeners[terminate.serverId] != nil {
				listeners[terminate.serverId].Close()
				listeners[terminate.serverId] = nil
			}
			serversWg.Done()
		case <-allTerminated:
			if status&uint(2) == uint(2) {
				log.Println("requests listeners stopped with some errors")
			} else {
				log.Println("all requests listeners stopped")
			}
			if signalStopFlag != true {
				signal.Stop(termsig)
			}

			// clean up
			close(terminateChannel)
			close(termsig)

			return status
		}
	}

}

func startServer(configs []ServerConfig) ([]net.Listener, chan terminateSignal, error) {
	terminate := make(chan terminateSignal, len(configs))
	listeners := make([]net.Listener, 0, len(configs))
	for serverId, config := range configs {
		log.Printf("starting requests listener on %v\n", config.Address)
		// create listener
		listener, err := createListener(config)
		if err != nil {
			for listenerId := range listeners {
				listeners[listenerId].Close()
				listeners[listenerId] = nil
			}
			return listeners, terminate, err
		} else {
			listeners = append(listeners, listener)
		}

		// start server
		var wg sync.WaitGroup
		server := &http.Server{
			Handler: config.ServeMux,
			ConnState: func(conn net.Conn, state http.ConnState) {
				switch state {
				case http.StateNew:
					wg.Add(1)
				case http.StateClosed:
					wg.Done()
				}
			},
		}
		server.SetKeepAlivesEnabled(false)

		go func(server *http.Server, listener net.Listener, terminate chan<- terminateSignal, serverId int) {
			var status uint
			if err := server.Serve(listener); err != nil {
				if !strings.Contains(err.Error(), "use of closed network connection") {
					log.Println(err)
					status = 1
				} else {
					status = 0
				}
			} else {
				status = 0
			}

			log.Printf("server #%v: performing graceful close on all clients\n", serverId)
			wg.Wait()

			terminate <- terminateSignal{
				status:   status,
				serverId: serverId,
			}
		}(server, listener, terminate, serverId)
	}
	return listeners, terminate, nil
}

func createListener(config ServerConfig) (net.Listener, error) {
	listener, err := net.Listen("tcp", config.Address)
	if err != nil {
		return nil, err
	}

	if config.TlsEnable {
		// use TLS connection
		tlsConfig := &tls.Config{}
		tlsConfig.Certificates = config.TlsCertificates
		return tls.NewListener(listener, tlsConfig), nil
	} else {
		return listener, nil
	}
}

type terminateSignal struct {
	status   uint
	serverId int
}
