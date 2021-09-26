/*
   Download an item from an AWS S3 Bucket
*/
package main

import (
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func handlerDownload(w http.ResponseWriter, r *http.Request) {

	// We get the name of the file on the URL
	filename := strings.Replace(r.URL.Path, "/get/", "", 1)

	f, err := os.Create(filename)
	if err != nil {
		showError(w, r, http.StatusBadRequest, "Something went wrong creating the local file")
		return
	}

	// Write the contents of S3 Object to the file
	downloader := s3manager.NewDownloader(sess)
	_, err = downloader.Download(f, &s3.GetObjectInput{
		Bucket: aws.String(AWS_S3_BUCKET),
		Key:    aws.String(filename),
	})
	if err != nil {
		showError(w, r, http.StatusBadRequest, "Something went wrong retrieving the file from S3")
		return
	}

	http.ServeFile(w, r, filename)
}
