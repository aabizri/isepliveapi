// TODO: Validity checker for necessary fields in Publications
// TODO: Validity checker for file size
// TODO: Login method to get the cookie
package main

import (
	"flag"
	"log"
	"net/http"
	"path/filepath"
)

var portFlag = flag.String("p", "8080", "Port to serve")
var port string

var noAuthFlag = flag.Bool("noauth", false, "Pass-through all request (WON'T WORK PAST LOG)")

var tempDirFlag = flag.String("temp", "/tmp/", "Temporary directory for image storage")
var tempDirPath string

func init() {
	flag.Parse()
	port = ":" + *portFlag

	tempDirPath = filepath.Clean(*tempDirFlag)
}

// The session
type session struct {
	// Client
}

func main() {
	// Create a goil client & session
	sess := &session{}

	// Get a router
	router := newRouter()

	// Use auth
	mainroute := sess.genAuthHandler(router)

	// Log errors
	log.Fatal(http.ListenAndServe(port, mainroute))
}
