package main

import (
    "fmt"
    "net/http"
    "os"
    "os/signal"
    "runtime"
    "sync"
    "syscall"
)

var mu = &sync.Mutex{}
var ROOT string

var running = true

func cleanup() {
    running = false

    fmt.Println("cleaning up...")

    if config.PersistQueue {
        fmt.Println("queue persistence enabled, saving queue...")
        WriteQueue()
    }
}

func main() {
    ROOT = os.Getenv("ROOT")

    if ROOT == "" {
        fmt.Println("root unset")
        os.Exit(1)
    }

    runtime.GOMAXPROCS(runtime.NumCPU())

    // set the main queue
    queue = append(queue, job_set{})

    // init config
    ConfigLoad()

    // handle sigterm
    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt, syscall.SIGTERM)
    go func() {
        <-c
        cleanup()
        os.Exit(1)
    }()

    // handle service endpoints
    http.HandleFunc("/", ServiceStatus)
    http.HandleFunc("/jobs", ServiceJobs)
    http.HandleFunc("/schedule", ServiceSchedule)
    http.HandleFunc("/receive_block", ServiceReceiveBlock)
    http.HandleFunc("/receive_noblock", ServiceReceiveNoBlock)

    fmt.Println("server started on", config.Address)

    http.ListenAndServe(config.Address, nil)
}
