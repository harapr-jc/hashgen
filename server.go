// Copyright 2017 <CompanyName>, Inc. All Rights Reserved.

// Provides server and utilities for generating cryptographic hash for a password.
// In a package for reuse, since one can't import a main package.
package hashgen

import (
	"encoding/base64"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"strings"
	"sync"
	"time"
)

// Server bits
var (
	host       string
	port       string
	hashServer *http.Server
	stats      Stats
	// Evicts after 1000 entries, but all records persisted to disk
	cache *LruCache = NewCache(1000, "backup.json")
)

// For monitoring shutdown request
var (
	// Exported so that calling process can terminate gracefully
	Exit = make(chan struct{})

	// Used to internally coordinate graceful shutdown
	exit = make(chan struct{})

	// The shutdown request is Ctrl-c
	quit = make(chan os.Signal)

	// Reference count for crypto hash jobs
	waiter sync.WaitGroup
)

// TODO: What is the right way to do this, compose with server?
var shutdownState struct {
	sync.RWMutex        // same as adding the methods of sync.RWMutex
	isShutdownRequested bool
}

func HandleGetHashRequest(w http.ResponseWriter, req *http.Request) {

	// TODO: Might be more elegant to wrap this handler in a timer handler
	start := time.Now()
	defer stats.Accumulate(start)

	// Validate parameters (last element of path)
	uuid := path.Base(req.URL.Path)
	if len(uuid) != 36 {
		http.Error(w, "Error: illegal Id", http.StatusBadRequest)
		return
	}

	userRecord, ok := cache.Get(uuid)
	if !ok {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	encodedResult := base64.URLEncoding.EncodeToString(userRecord.HashBytes)
	io.WriteString(w, encodedResult+"\n")
}

func HandleHashRequest(w http.ResponseWriter, req *http.Request) {

	// TODO: Might be more elegant to wrap this handler in a timer handler
	start := time.Now()
	defer stats.Accumulate(start)

	// We can easily support a timeout by creating a context here.
	// See https://blog.golang.org/context

	// Validate parameters
	req.ParseForm()
	password := req.FormValue("password")
	if len(password) == 0 {
		http.Error(w, "Error: missing password", http.StatusBadRequest)
		return
	}
	// Adding salt is optional.
	saltParam := req.FormValue("salt")
	addSalt := false
	if saltParam == "yes" {
		addSalt = true
	}

	// While holding read lock on shutdown not requested, add to wait group
	shutdownState.RLock()
	shutdownPending := shutdownState.isShutdownRequested
	if shutdownPending == false {
		waiter.Add(1)
	}
	shutdownState.RUnlock()
	if shutdownPending {
		http.Error(w, "Error: service shutdown pending", http.StatusServiceUnavailable)
		return
	}

	// Computes unique job identifier
	//uuid, err := getUuid()
	uuid := getUuid()

	// Launch a goroutine to compute the crypo hash
	go func() {
		// Requirement: Must sleep for 5 seconds. Presumably this simulates
		// password stretching and iteration time.
		time.Sleep(5 * time.Second)
		// Decrements the waiter when goroutine completes
		defer waiter.Done()
		// Work against lookup table and rainbow table attacks by randomizing the hash.
		saltSize := 0
		if addSalt {
			saltSize = DefaultDigestSize
		}
		cryptoBytes := getCryptoHash([]byte(password), saltSize)
		cache.Add(uuid, []byte(""), cryptoBytes)
	}()

	// Returns the job identifier
	io.WriteString(w, uuid+"\n")
}

// Returns StatsMessage as Json
func HandleStatsRequest(w http.ResponseWriter, req *http.Request) {

	io.WriteString(w, string(stats.GetJson())+"\n")
}

func StartServer(host string, port string) *http.Server {

	s := []string{host, port}
	address := strings.Join(s, ":")
	hashServer = &http.Server{Addr: address}

	// Register all endpoints
	// Measure duration for hash GET and POST handlers
	http.HandleFunc("/hash", HandleHashRequest)
	http.HandleFunc("/hash/", HandleGetHashRequest)
	http.HandleFunc("/stats", HandleStatsRequest)

	// Subscribe to SIGINT
	signal.Notify(quit, os.Interrupt)

	// The crypto hash server supports a shutdown request. Work in progress
	// completes while shutdown pending. Therefore, run server in goroutine.
	go func() {

		// In the background, wait for a quit (Ctrl-c) request
		go func() {
			// If this function completes with a panic, run recovery.
			defer func() {
				if err := recover(); err != nil {
					log.Printf("Recovered: %+v\n", err)
				}
			}()
			// Blocks routine pending a quit (Ctrl-c) request
			<-quit
			// All crypto hash jobs must complete prior to shutdown,
			// with the output persisted to disk.
			log.Println("Shutdown requested...")
			// Take a write lock for mutating state
			shutdownState.Lock()
			shutdownState.isShutdownRequested = true
			shutdownState.Unlock()
			// Waits for all jobs to complete
			waiter.Wait()
			log.Println("Done waiting")
			log.Println("Shutting down server...")
			// This is a closure, we can use hashServer from calling function
			if err := hashServer.Shutdown(nil); err != nil {
				log.Fatalf("Failed to shutdown: %v", err)
			}
			// Notify thread monitoring the 'exit' channel
			close(exit)
		}()

		log.Printf("Starting crypto hash server at %q\n", address)
		if err := hashServer.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				log.Fatalf("Httpserver: ListenAndServe() error: %s\n", err)
			}
			// Blocks routine pending graceful shutdown
			<-exit
			log.Println("Server shutdown complete. Bye!")
			// Notify main() monitoring the 'Exit' channel
			close(Exit)
		}
	}()

	return hashServer
}
