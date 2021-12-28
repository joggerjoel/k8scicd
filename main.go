package main

import (
	"log"
	"net/http"
	"net"
)

type Server struct{}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	
	var ip net.IP
	ifaces, err := net.Interfaces()
	// handle err
	for _, i := range ifaces {
	    addrs, err := i.Addrs()
	    // handle err
	    for _, addr := range addrs {
		switch v := addr.(type) {
		case *net.IPNet:
			ip = v.IP
		case *net.IPAddr:
			ip = v.IP
		}
		// process IP address
	    }
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "Hello World: ", ip}`))
	
	
}

func main() {
	s := &Server{}
	http.Handle("/", s)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
