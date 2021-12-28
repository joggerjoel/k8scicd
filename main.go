package main

import (
  "log"
  "flag"
	"context"
	"fmt"
	"time"
	"path/filepath"
	"io/ioutil"

        "os"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	//
	// Uncomment to load all auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth"
	//
	// Or uncomment to load specific auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth/azure"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/openstack"
	
  "net/http"
  "net"
  "strconv"
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
	
// creates the in-cluster config
	path, err := os.Getwd()
	if err != nil {
            panic(err.Error())
	}
	fmt.Printf("GO PATH: %s\n",path)  // for example /home/user
	files, err := ioutil.ReadDir("../")
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
	    	fmt.Println(f.Name())
	}
	
	
	var kubeconfig *string
	kubeconfig = flag.String("kubeconfig", filepath.Join("/var/lib/jenkins/workspace", "config"), "(optional) absolute path to the kubeconfig file")

	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	for {
		// get pods in all the namespaces by omitting namespace
		// Or specify namespace to get pods in particular namespace
		pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))
		w.Write([]byte("There are "))
		w.Write([]byte(strconv.Itoa(len(pods.Items))))
		w.Write([]byte("pods in the cluster"))

		// Examples for error handling:
		// - Use helper functions e.g. errors.IsNotFound()
		// - And/or cast to StatusError and use its properties like e.g. ErrStatus.Message
		_, err = clientset.CoreV1().Pods("default").Get(context.TODO(), "example-xxxxx", metav1.GetOptions{})
		if errors.IsNotFound(err) {
			fmt.Printf("Pod example-xxxxx not found in default namespace\n")
		} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
			fmt.Printf("Error getting pod %v\n", statusError.ErrStatus.Message)
		} else if err != nil {
			panic(err.Error())
		} else {
			fmt.Printf("Found example-xxxxx pod in default namespace\n")
		}

		time.Sleep(10 * time.Second)
	}
	w.Write([]byte(`"}`))
	
	
}

func main() {	
	s := &Server{}
	http.Handle("/", s)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
