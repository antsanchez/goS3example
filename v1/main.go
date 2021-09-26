/*
   This is a simple code example for connecting, uploading, downloading and listing files
   from an AWS S3 Bucket.
   Author: Antonio Sanchez antonio@asanchez.dev
*/

package main

import (
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

const (
	AWS_S3_REGION = ""
	AWS_S3_BUCKET = ""
)

var sess = connectAWS()

func connectAWS() *session.Session {
	sess, err := session.NewSession(&aws.Config{Region: aws.String(AWS_S3_REGION)})
	if err != nil {
		panic(err)
	}
	return sess
}

func main() {

	http.HandleFunc("/upload/", handlerUpload) // Upload
	http.HandleFunc("/get/", handlerDownload)  // Get the file
	http.HandleFunc("/list/", handlerList)     // List all files
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func showError(w http.ResponseWriter, r *http.Request, status int, message string) {
	http.Error(w, message, status)
}
