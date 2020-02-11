package main

/*
newe netcat runs
Connected to localhost.
Escape character is '^]'.
ls -la
total 56
drwxrwxr-x. 12 michal michal 4096 Feb 11 16:43 .
drwxrwxr-x. 26 michal michal 4096 Feb 11 13:32 ..
drwxrwxr-x.  2 michal michal 4096 Feb 10 23:31 1_portscanner
drwxrwxr-x.  2 michal michal 4096 Feb 10 23:31 2_11_echo
drwxrwxr-x.  2 michal michal 4096 Feb 10 23:31 2_12_echo
drwxrwxr-x.  2 michal michal 4096 Feb 10 23:31 2_13_echo
drwxrwxr-x.  2 michal michal 4096 Feb 10 23:31 2_13_forward_proxy
drwxrwxr-x.  2 michal michal 4096 Feb 11 16:43 2_14_https_forward_proxy
drwxrwxr-x.  2 michal michal 4096 Feb 11 17:00 2_15_netcat_pipes
drwxrwxr-x.  5 michal michal 4096 Feb 10 11:26 common
drwxrwxr-x.  8 michal michal 4096 Feb 11 16:44 .git
-rw-rw-r--.  1 michal michal  202 Feb  5 17:02 .gitignore
drwxrwxr-x.  2 michal michal 4096 Feb 11 17:01 .idea
-rw-rw-r--.  1 michal michal    6 Feb  6 14:47 README.md

*/


import (
	"io"
	"log"
	"net"
	"os/exec"
)

func handle(conn net.Conn){
	//explicitly calling /bin/sh using -i for interactive mode
	//so than we can use stdin and stdout
	//for windows use exec.Command("cmd.exe")
	cmd := exec.Command("/bin/sh","-i")

	rp, wp := io.Pipe()

	cmd.Stdin = conn

	cmd.Stdout = wp

	//preventing code from blocking
	//thread is blocked until connection ends
	//Copy(dst, src)
	go io.Copy(conn, rp)

	//run is blocking function.
	//to exit type "exit", "exit"
	cmd.Run()

	println("closed")
	conn.Close()
}



func main() {
	listener, err := net.Listen("tcp", ":20080")
	if err != nil {
		log.Fatalln(err)
	}

	for {
		//blocking and unblock only when someone new connects i.e. through telnet
		conn, err := listener.Accept()
		println("new connection")
		if err != nil {
			log.Fatalln(err)
		}
		//each new connection calls this
		go handle(conn)
	}
}