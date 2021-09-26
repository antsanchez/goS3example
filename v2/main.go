/*
   This is a simple code example for connecting, uploading, downloading and listing files
   from an AWS S3 Bucket using the AWS SDK v2 for Go.
   Author: Antonio Sanchez antonio@asanchez.dev
*/

package main

import (
	"context"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const (
	AWS_S3_REGION = "" // Region
	AWS_S3_BUCKET = "" // Bucket
)

// We will be using this client everywhere in our code
var awsS3Client *s3.Client

func main() {
	configS3()

	http.HandleFunc("/upload", handlerUpload)     // Upload: /upload (upload file named "file")
	http.HandleFunc("/download", handlerDownload) // Download: /download?key={key of the object}&filename={name for the new file}
	http.HandleFunc("/list", handlerList)         // List: /list?prefix={prefix}&delimeter={delimeter}
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// configS3 creates the S3 client
func configS3() {

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(AWS_S3_REGION))
	if err != nil {
		log.Fatal(err)
	}

	awsS3Client = s3.NewFromConfig(cfg)
}

func showError(w http.ResponseWriter, r *http.Request, status int, message string) {
	http.Error(w, message, status)
}
