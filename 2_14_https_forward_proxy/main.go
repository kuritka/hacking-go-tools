package main

import (
	"crypto/tls"
	"hck/common/guard"
	"io"
	"net/http"
	"net/url"

	"hck/common/log"
)

const (
	key  = "/etc/lb/certs/key.pem"
	cert = "/etc/lb/certs/cert.pem"
	port = ":9443"
	host = "ulozto.cz"
)

type webRequest struct {
	r      *http.Request
	w      http.ResponseWriter
	doneCh chan struct{}
}

var (
	requestCh = make(chan *webRequest)
	logger    = log.Log
)

var (
	//TLS error occurs when false. Requires real certificate
	transport = http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	client http.Client
)

func init() {
	http.DefaultClient = &http.Client{Transport: &transport}
	client = http.Client{Transport: &transport}
}

func processRequests() {

	for request := range requestCh {

		println("request")

		go processRequest(host, request)
	}
}

func processRequest(host string, request *webRequest) {
	//build url for new host
	hostUrl, _ := url.Parse(request.r.URL.String())
	hostUrl.Scheme = "https"
	hostUrl.Host = host
	println(host)
	println(hostUrl.String())
	req, _ := http.NewRequest(request.r.Method, hostUrl.String(), request.r.Body)
	//because request headers in go is map of slice of strings we must translate into string of headers to new request
	for k, v := range request.r.Header {
		values := ""
		for _, headerValue := range v {
			values += headerValue + " "
		}
		//to slice of strings
		req.Header.Add(k, values)
	}

	resp, err := client.Do(req)
	if err != nil {
		request.w.WriteHeader(http.StatusInternalServerError)
		request.doneCh <- struct{}{}
		return
	}
	//now we have response headers to work with
	//in production we will need some exceptions here
	//we don't want send any headers making security issues within organisation
	for key, header := range resp.Header {
		headers := ""
		for _, headerValue := range header {
			headers += headerValue + " "
		}
		request.w.Header().Add(key, headers)
	}
	_, _ = io.Copy(request.w, resp.Body)

	request.doneCh <- struct{}{}
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		doneCh := make(chan struct{})
		requestCh <- &webRequest{r: r, w: w, doneCh: doneCh}
		//waits until LoadBalancer resend request to chosen app server and resend response back
		//or error happens
		<-doneCh
	})

	go processRequests()

	//load balancing
	//if nil DefaultServerMux is used and DSM gets registered handler
	logger.Info().Msgf("listening on %s", port)
	err := http.ListenAndServeTLS(port, cert, key, nil)

	guard.FailOnError(err, "server didn't start")
	//log.Println("server started, press <ENTER> to exit")
	//_, _ = fmt.Scanln()

}
