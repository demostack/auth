package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/demostack/auth/endpoint"

	"github.com/akrylysov/algnhsa"
)

var (
	// AWSAccessKeyID ..
	AWSAccessKeyID = ""

	// AWSSecretAccessKey .
	AWSSecretAccessKey = ""
)

func init() {
	// Verbose logging with file name and line number.
	log.SetFlags(log.Lshortfile)
}

func main() {
	// Set the AWS environment variables if they are passed as build flags.
	if len(AWSAccessKeyID) > 0 {
		os.Setenv("AWS_ACCESS_KEY_ID", AWSAccessKeyID)
	}
	if len(AWSSecretAccessKey) > 0 {
		os.Setenv("AWS_SECRET_ACCESS_KEY", AWSSecretAccessKey)
	}

	c := endpoint.Boot()
	http.HandleFunc("/", c.Hander)

	if len(os.Getenv("AWS_LAMBDA_FUNCTION_NAME")) > 0 {
		// Run lambda server.
		fmt.Println("lamba started")
		algnhsa.ListenAndServe(http.DefaultServeMux, nil)
	} else {
		// Run web server.
		fmt.Println("listening as web server")
		fmt.Println("server started")
		http.ListenAndServe(":8080", nil)
	}

}
