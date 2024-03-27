package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
)

const version = "1.0.0"

var (
	lastExecution time.Time
	mutex         sync.Mutex
)

func main() {
	fmt.Printf("Go Webhook version %v\n", version)
	// Define and parse command-line flags
	endpoint := flag.String("u", "", "The endpoint url to listen for requests, starting with \"/\"")
	port := flag.String("p", "7999", "The port number to listen on")
	command := flag.String("c", "./script.sh", "The command to execute when the webhook URL is invoked")
	method := flag.String("m", "GET", "HTTP method of the webhook url")
	throttle := flag.String("w", "10", "Duration in seconds to wait between the command executions")
	flag.Parse()

	if *endpoint == "" || (*endpoint)[0] != '/' {
		fmt.Printf("Please provide an endpoint url with -u (starting with '/').\nExample: -u /webhook/2ff80e9159b517704ce43f0f74e6e247\n")
		os.Exit(1)
		return
	}

	seconds, err := strconv.Atoi(*throttle)
	if err != nil {
		fmt.Printf("Error: Failed to parse -w argument: %v\n", err)
		os.Exit(1)
		return
	}
	waitDur := time.Duration(seconds) * time.Second

	spl := strings.Split(*command, " ")
	args := make([]string, 0, len(spl)-1)
	if len(spl) > 1 {
		args = append(args, spl[1:]...)
	}

	// Set up the HTTP handler
	http.HandleFunc(*endpoint, func(w http.ResponseWriter, r *http.Request) {
		webhookHandler(w, r, spl[0], args, *endpoint, *method, waitDur)
	})

	// Start the server

	fmt.Printf("Webhook test url: %s http://127.0.0.1:%s%s\n", *method, *port, *endpoint)
	if err = http.ListenAndServe(":"+*port, nil); err != nil {
		fmt.Printf("Error: Failed to start the server: %v\n", err)
		os.Exit(1)
		return
	}
}

func webhookHandler(w http.ResponseWriter, r *http.Request, command string, args []string, endpoint string, method string, throttleDur time.Duration) {
	if r.Method != method {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Path != endpoint {
		http.Error(w, "Invalid request URL \""+r.URL.Path+"\"", http.StatusNotFound)
		return
	}

	if !mutex.TryLock() {
		http.Error(w, "Request throttled. Please try again later.", http.StatusTooManyRequests)
		return
	}
	defer mutex.Unlock()

	now := time.Now()
	if now.Sub(lastExecution) < throttleDur {
		http.Error(w, "Request throttled. Please try again later.", http.StatusTooManyRequests)
		return
	}

	log.Printf("Running command: %s %s", command, strings.Join(args, " "))
	started := time.Now()
	cmd := exec.Command(command, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println(err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(output)

	now = time.Now()
	lastExecution = now
	log.Printf("Command executed successfully in %d ms", now.Sub(started).Milliseconds())
}
