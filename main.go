package main

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"os"
	"os/signal"
	"time"
	"log"
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

//
// Pod infos
//

func GetPodDetails() () {

    // creates the in-cluster config
    config, err := rest.InClusterConfig()
    if err != nil {
        panic(err.Error())
    }
    // creates the clientset
    clientset, err := kubernetes.NewForConfig(config)
    if err != nil {
        panic(err.Error())
    }

    IP = GetLocalIP()
    for {
        if IP != "" {
            break
        } else {
            log.Printf("No IP for now.\n")
        }

        pods, err := clientset.CoreV1().Pods("default").List(metav1.ListOptions{})
        if err != nil {
            panic(err.Error())
        }
        for _, pod := range pods.Items {
            pod, _ := clientset.CoreV1().Pods("default").Get(pod.Name, metav1.GetOptions{})
            if pod.Name == os.Getenv("HOSTNAME") {
                IP = pod.Status.PodIP
            }
        }

        log.Printf("Waits...\n")
        time.Sleep(1 * time.Second)
    }

    name = os.Getenv("HOSTNAME")
    log.Printf("Trying os.Getenv(\"HOSTNAME/IP\"): [%s][%s]\n", name, IP)

    return IP, name
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "Hello World: `))
	w.Write([]byte(GetLocalIP()))
	w.Write([]byte(`/`))
	w.Write([]byte(GetLocalIP()))
	w.Write([]byte(`"}`))
	
	
}

func main() {
	s := &Server{}
	http.Handle("/", s)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
