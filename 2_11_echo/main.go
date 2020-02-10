/*
Simple echo server. Connection implements Reader and Writer interfaces so
.Read([]byte) (size, error), .Write([]byte) (size, error) are present.
listener.Accept is blocking function and returns new instance of connection.
Thats why we handle echo in different thread with different connection instance


[michal@localhost hck]$ telnet localhost 20080
Trying ::1...
Connected to localhost.
Escape character is '^]'.
hello from CLIENT 2
hello from CLIENT 2


[michal@localhost hck]$ telnet localhost 20080
Trying ::1...
Connected to localhost.
Escape character is '^]'.
CLIENT 1 IS SCREAMING
CLIENT 1 IS SCREAMING


{"level":"info",:"listening on port :20080"}
{"level":"info",:"received connection connection-beb83e31-e404-d052-262a-65dcb106b5dd"}
{"level":"info",:"received connection connection-ed4f906b-798e-236c-7fcf-0f9374f971c9"}

"message":"received 21 bytes: hello from CLIENT 2\
{"level":"info","time","message":"writing data  connection-ed4f906b-798e-236c-7fcf-0f9374f971c9"}

hello from CLIENT 2
hello from CLIENT 2




{"level":"info","message":"received 23 bytes: CLIENT 1 IS SCREAMING\r\n
{"level":"info","message":"writing data  connection-beb83e31-e404-d052-262a-65dcb106b5dd"}

CLIENT 1 IS SCREAMING
CLIENT 1 IS SCREAMING

//now I closed terminal with terminal 2
{"level":"error""message":"client disconnected connection-ed4f906b-798e-236c-7fcf-0f9374f971c9"}
{"level":"info","message":"closed  connection-ed4f906b-798e-236c-7fcf-0f9374f971c9"}

//now closing terminal with Client1
{"level":"error",:"client disconnected connection-beb83e31-e404-d052-262a-65dcb106b5dd"}
{"level":"info","closed  connection-beb83e31-e404-d052-262a-65dcb106b5dd"}


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


func echo(conn net.Conn, id string){
	defer closeEcho(conn, id)

	//buffer
	b := make([]byte,1024)

	for {
		size, err := conn.Read(b[0:])
		if err == io.EOF {
			logger.Error().Msgf("client disconnected %s",id)
			break
		}
		if err != nil {
			logger.Error().Msgf("unexpected error %s",id)
			break
		}

		logger.Info().Msgf("received %d bytes: %s\n", size,string(b))

		logger.Info().Msgf("writing data  %s",id)

		_,err = conn.Write(b[0:size])

		guard.FailOnError(err,"unable to write data %s",id)
	}
}

func closeEcho(conn net.Conn, id string){

	err := conn.Close()

	guard.FailOnError(err,"unable to close connection  %s",id)

	logger.Info().Msgf("closed  %s",id)
}
