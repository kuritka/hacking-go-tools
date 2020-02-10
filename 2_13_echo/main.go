/*
Logger configured
{"level":"info","message":"listening on port :20080"}
{"level":"info","message":"received connection connection-2b34682f-c53a-53fd-64ab-73e330c542e7"}


[michal@localhost hck]$ telnet localhost 20080
Trying ::1...
Connected to localhost.
Escape character is '^]'.
ababa
ababa
babab
babab


{"level":"info","message":"transfered 25 , connection-2b34682f-c53a-53fd-64ab-73e330c542e7"}
{"level":"info","message":"closed  connection-2b34682f-c53a-53fd-64ab-73e330c542e7"}


*/

package main

import (
	"fmt"
	"io"
	"net"

	"hck/common/guard"
	"hck/common/guid"
	"hck/common/log"
)

const port = ":20080"

var logger = log.Log

func main() {

	listener, err := net.Listen("tcp", ":20080")

	guard.FailOnError(err, "unable to bind port %s", port)

	logger.Info().Msgf("listening on port %s", port)

	for {
		//blocking function, open new connection when new request comes
		conn, err := listener.Accept()

		guard.FailOnError(err, "unable to accept connection")

		id, _ := guid.Guid()
		id = fmt.Sprintf("connection-%s", id)

		logger.Info().Msgf("received connection %s", id)

		//function close connection after client disconnects or errors.
		//Don't forget - multiple clients can connect to one echo server
		go echo(conn, id)
	}
}

//look to commented echo where direct conn.Read / Write are called
func echo(conn net.Conn, id string) {

	defer closeEcho(conn, id)

	//conn implements Reader, WRITER, so can be used on this place as parameter
	//Copy is blocking function - as Read or BufRead for instance, and exiting when error occurs or client disconnects
	//thats why you cannot control what is copied
	s, err := io.Copy(conn, conn)

	guard.FailOnError(err, "unable to copy data %s", id)

	logger.Info().Msgf("transfered %d , %s",s, id)
}

func closeEcho(conn net.Conn, id string) {

	err := conn.Close()

	guard.FailOnError(err, "unable to close connection  %s", id)

	logger.Info().Msgf("closed  %s", id)
}
