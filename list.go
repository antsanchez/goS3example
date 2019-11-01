/*
   List all items from an AWS S3 Bucket
*/
package main

import (
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
)

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
