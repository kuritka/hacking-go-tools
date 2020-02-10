package main

import (
	"flag"
	"fmt"
	"io"
	"net"

	"hck/common/guard"
	"hck/common/log"
)

var logger = log.Log


//forwards traffic from source to specified url and back
func forward(src net.Conn, url string,  port int ) {

	//working for www.columbia.edu:80, than http://localhost:8443 automatically redirects
	//ulozto.cz:443 - is HTTPS and listeners close connection as https
	dst, err := net.Dial("tcp", fmt.Sprintf("%s",url))
	guard.FailOnError(err, "unable to make connection to %s or host is unreachable ", url)

	defer dst.Close()
	defer src.Close()

	go func() {
		//copy is blocking func
		_, err := io.Copy(dst, src)
		guard.FailOnError(err, "unable to copy from %s to source", url)
	}()

	//copy is blocking func
	_, err = io.Copy(src, dst)
	guard.FailOnError(err, "unable to copy from source to %s", url)
}

func main() {
	var url = flag.String("url", "https://ulozto.cz", "url adress where is proxy pointing i.e. https://forbidden.cz")

	var port = flag.Int("port", 443, "port where is proxy listening ")

	flag.Parse()

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))

	guard.FailOnError(err, "unable to bind port :%v", *port)

	defer listener.Close()

	logger.Info().Msgf("forwarding on %s listening on port %d", *url,*port)

	for {
		//blocking operation accepts new connections from multiple sources
		conn, err := listener.Accept()

		guard.FailOnError(err, "unable to accept connection")

		go forward(conn, *url, *port)
	}

}
