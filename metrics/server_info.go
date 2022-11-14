package metrics

import (
	"log"
	"os"
	"strconv"
	"strings"
)

type ServerInfo struct {
	ServiceName string

	ServiceInstance string
}

func NewServerInfo(serviceName, serviceInstance string) ServerInfo {
	return ServerInfo{
		ServiceName:     serviceName,
		ServiceInstance: serviceInstance,
	}
}

func NewServerInfoInK8sCluster() ServerInfo {
	podName := os.Getenv("MY_POD_NAME") // inject pod name by k8s Downward API
	if len(podName) != 0 {
		log.Println("get pod name from env[MY_POD_NAME], pod name: " + podName)
	} else {
		podName = os.Getenv("HOSTNAME")
		if len(podName) != 0 {
			log.Println("get pod name from env[HOSTNAME], pod name: " + podName)
		}
	}
	if podName == "" {
		podName = "unknown"
	}

	serviceName := os.Getenv("MY_SERVICE_NAME")
	if len(serviceName) != 0 {
		log.Println("get service name from env[MY_SERVICE_NAME], service name: " + serviceName)
	} else {
		serviceName = guessServiceName(serviceName)
	}
	log.Printf("pod name: %s, service name: %s\n", podName, serviceName)
	return ServerInfo{
		ServiceName:     serviceName,
		ServiceInstance: podName,
	}
}

// guessServiceName guess service name by pod name
func guessServiceName(podName string) string {
	if podName == "unknown" {
		return "unknown"
	}
	split := strings.Split(podName, "-")
	if len(split) == 0 { // maybe use pod.yaml
		return podName
	}
	last := split[len(split)-1]

	stsNum, err := strconv.Atoi(last)
	if err == nil {
		if stsNum >= 0 && stsNum < 99 {
			return podName[:strings.LastIndex(podName, "-")]
		}
	}

	//maybe deployment, job
	if len(split) > 2 {
		temp := podName[:strings.LastIndex(podName, "-")]
		return temp[0:strings.LastIndex(temp, "-")]
	}
	//unknown
	return podName
}
