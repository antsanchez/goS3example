/*
   This is a simple code example for connecting, uploading, downloading and listing files
   from an AWS S3 Bucket.
   Author: Antonio Sanchez antonio@asanchez.dev
*/

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

var sess = connectAWS()

func connectAWS() *session.Session {
	sess, err := session.NewSession(&aws.Config{Region: aws.String(AWS_S3_REGION)})
	if err != nil {
		panic(err)
	}
	return sess
}

const (
	AWS_S3_REGION = ""
	AWS_S3_BUCKET = ""
)

func main() {

	http.HandleFunc("/upload/", handlerUpload) // Upload
	http.HandleFunc("/get/", handlerDownload)  // Get the file
	http.HandleFunc("/list/", handlerList)     // List all files
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handlerUpload(w http.ResponseWriter, r *http.Request) {

	r.ParseMultipartForm(10 << 20)

	// Get a file from the form input name "file"
	file, header, err := r.FormFile("file")
	if err != nil {
		showError(w, r, http.StatusInternalServerError, "Something went wrong retrieving the file from the form")
		return
	}
	defer file.Close()

	filename := header.Filename
	fmt.Println(filename)

	// Upload the file to S3.
	uploader := s3manager.NewUploader(sess)

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(AWS_S3_BUCKET), // Bucket
		Key:    aws.String(filename),      // Name of the file to be saved
		Body:   file,                      // File
	})
	if err != nil {
		// Do your error handling here
		showError(w, r, http.StatusInternalServerError, "Something went wrong uploading the file")
		return
	}

	fmt.Fprintf(w, "Successfully uploaded to %q\n", AWS_S3_BUCKET)
	return
}

func handlerDownload(w http.ResponseWriter, r *http.Request) {

	// We get the name of the file on the URL
	filename := strings.Replace(r.URL.Path, "/get/", "", 1)

	downloader := s3manager.NewDownloader(sess)

	f, err := os.Create(filename)
	if err != nil {
		showError(w, r, http.StatusBadRequest, "Something went wrong creating the local file")
		return
	}

	// Write the contents of S3 Object to the file
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

func handlerList(w http.ResponseWriter, r *http.Request) {

	svc := s3.New(sess)
	input := &s3.ListObjectsInput{
		Bucket: aws.String(AWS_S3_BUCKET),
	}

	result, err := svc.ListObjects(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchBucket:
				fmt.Println(s3.ErrCodeNoSuchBucket, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		showError(w, r, http.StatusBadRequest, "Something went wrong listing the files")
		fmt.Println(err)
		return
	}

	w.Header().Set("Content-Type", "text/html")

	for _, item := range result.Contents {
		fmt.Fprintf(w, "<li>File %s</li>", *item.Key)
	}

	return
}

func showError(w http.ResponseWriter, r *http.Request, status int, message string) {
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprintf(w, message)
}
