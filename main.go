package main

import (
  "os"
  "log"
  "path/filepath"	
  "context"
  "k8s.io/apimachinery/pkg/apis/meta/v1"
  "k8s.io/client-go/kubernetes"
  "k8s.io/client-go/tools/clientcmd"
  "net/http"
  "net"
)

type Server struct{}

// GetLocalIP returns the non loopback local IP of the host
func GetLocalIP() string {
    addrs, err := net.InterfaceAddrs()
    if err != nil {
        return ""
    }
    for _, address := range addrs {
        // check the address type and if it is not a loopback the display it
        if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
            if ipnet.IP.To4() != nil {
                return ipnet.IP.String()
            }
        }
    }
    return ""
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "Hello World: `))
	w.Write([]byte(GetLocalIP()))
	w.Write([]byte(`/`))
	w.Write([]byte(GetLocalIP()))
	
	kubeconfig := filepath.Join(
	    os.Getenv("HOME"), ".kube", "config",
	)
	// uses the current context in kubeconfig
	// path-to-kubeconfig -- for example, /root/.kube/config
	config, _ := clientcmd.BuildConfigFromFlags("", kubeconfig)
	// creates the clientset
	clientset, _ := kubernetes.NewForConfig(config)
	// access the API to list pods
	pods, _ := clientset.CoreV1().Pods("").List(context.TODO(), v1.ListOptions{})

	w.Write([]byte("There are "))
	w.Write([]byte(len(pods.Items)))
	w.Write([]byte("in the cluster"))
	w.Write([]byte(`"}`))
	
	
}

func main() {	
	s := &Server{}
	http.Handle("/", s)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
