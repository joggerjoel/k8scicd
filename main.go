package main

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"os"
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
func getClient(configLocation string) (typev1.CoreV1Interface, error){
    kubeconfig := filepath.Clean(configLocation)
    config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
    if err != nil {
        log.Fatal(err)
    }
    clientset, err := kubernetes.NewForConfig(config)
    if err != nil {
        return nil, err
    }
    return clientset.CoreV1(), nil
}

func getServiceForDeployment(deployment string, namespace string, k8sClient typev1.CoreV1Interface) (*corev1.Service, error){
    listOptions := metav1.ListOptions{}
    svcs, err := k8sClient.Services(namespace).List(listOptions)
    if err != nil{
        log.Fatal(err)
    }
    for _, svc:=range svcs.Items{
        if strings.Contains(svc.Name, deployment){
            fmt.Fprintf(os.Stdout, "service name: %v\n", svc.Name)
            return &svc, nil
        }
    }
    return nil, errors.New("cannot find service for deployment")
}

func getPodsForSvc(svc *corev1.Service, namespace string, k8sClient typev1.CoreV1Interface) (*corev1.PodList, error){
    set := labels.Set(svc.Spec.Selector)
    listOptions:= metav1.ListOptions{LabelSelector: set.AsSelector().String()}
    pods, err:=  k8sClient.Pods(namespace).List(listOptions)
    for _,pod:= range pods.Items{
        fmt.Fprintf(os.Stdout, "pod name: %v\n", pod.Name)
    }
    return pods, err
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
	kubeconfig := filepath.Join(
	os.Getenv("HOME"), ".kube", "config",
	)
	namespace:="FOO"
	k8sClient, err:= getClient(kubeconfig)
	if err!=nil{
	fmt.Fprintf(os.Stderr, "error: %v\n", err)
	os.Exit(1)
	}

	svc, err:=getServiceForDeployment("APP_NAME", namespace, k8sClient)
	if err!=nil{
	fmt.Fprintf(os.Stderr, "error: %v\n", err)
	os.Exit(2)
	}

	pods, err:=getPodsForSvc(svc, namespace, k8sClient)
	if err!=nil{
	fmt.Fprintf(os.Stderr, "error: %v\n", err)
	os.Exit(2)
	}
	
	s := &Server{}
	http.Handle("/", s)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
