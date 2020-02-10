/*
Simple echo server. Connection implements Reader and Writer interfaces so
.Read([]byte) (size, error), .Write([]byte) (size, error) are present.
listener.Accept is blocking function and returns new instance of connection.
Thats why we handle echo in different thread with different connection instance


[michal@localhost hck]$ telnet localhost 20080
Trying ::1...
Connected to localhost.
Escape character is '^]'.
blaah
blaah
Connection closed by foreign host.

{"level":"info","message":"listening on port :20080"}
{"level":"info","message":"received connection connection-57a52e96-55a7-8e81-9a9e-f1e7a72b68a4"}
{"level":"info","message":"received 7 bytes: blaah\r\n\n"}
{"level":"info","message":"writing data  connection-57a52e96-55a7-8e81-9a9e-f1e7a72b68a4"}
{"level":"info","message":"closed  connection-57a52e96-55a7-8e81-9a9e-f1e7a72b68a4"}


*/

package main

import (
	"bufio"
	"fmt"
	"net"

	"hck/common/guard"
	"hck/common/guid"
	"hck/common/log"
)

const port = ":20080"

var logger = log.Log

func main(){

	listener,err := net.Listen("tcp", ":20080")

	guard.FailOnError(err,"unable to bind port %s",port)

	logger.Info().Msgf("listening on port %s",port)

	for {
		//blocking function, open new connection when new request comes
		conn, err := listener.Accept()

		guard.FailOnError(err,"unable to accept connection")

		id, _ := guid.Guid()
		id = fmt.Sprintf("connection-%s",id)

		logger.Info().Msgf("received connection %s",id)

		//function close connection after client disconnects or errors.
		//Don't forget - multiple clients can connect to one echo server
		go echo(conn, id)
	}
}


//look to commented echo where direct conn.Read / Write are called
func echo(conn net.Conn, id string){
	defer closeEcho(conn, id)

	//conn implements Reader, so can be used on this place as parameter

	reader := bufio.NewReader(conn)

	s, err := reader.ReadString('\n')

	guard.FailOnError(err,"unable to read data %s",id)

	logger.Info().Msgf("received %d bytes: %s\n", len(s),s)

	logger.Info().Msgf("writing data  %s",id)

	writer := bufio.NewWriter(conn)

	_, err = writer.Write([]byte(s))

	guard.FailOnError(err,"unable to write data %s",id)

	//writing data to underlying writer - in this case conn instance
	writer.Flush()

}


func closeEcho(conn net.Conn, id string){

	err := conn.Close()

	guard.FailOnError(err,"unable to close connection  %s",id)

	logger.Info().Msgf("closed  %s",id)
}


