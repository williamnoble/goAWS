package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"io/ioutil"
	"os"
	"strings"
)

var (
	s3session *s3.S3
)

const (
	bucketName = "william-s3-test-12-08-2021"
	region     = "eu-west-2" // london
	folder     = "images"
)

func main() {
	cfg := &aws.Config{
		Region: aws.String(region),
	}
	sess := s3.New(session.Must(session.NewSession(cfg)))

	files, _ := ioutil.ReadDir(folder)
	fmt.Println(files)
	for _, file := range files {
		if file.IsDir() {
			continue
		} else {
			uploadObject(sess, folder+"/"+file.Name())
		}
	}

	fmt.Println(listObjects(sess))

	for _, object := range listObjects(sess).Contents {
		getObject(sess, *object.Key) //dereference
		//deleteObject(sess, *object.Key)
	}

	fmt.Println(listObjects(sess))

}

func getObject(s *s3.S3, filename string) {
	fmt.Println("downloading ", filename)
	input := &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(filename),
	}
	resp, err := s.GetObject(input)
	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	err = ioutil.WriteFile(filename, body, 0644)
	if err != nil {
		panic(err)
	}
}

func deleteObject(s *s3.S3, filename string) (resp *s3.DeleteObjectOutput) {
	fmt.Println("Deleting: ", filename)
	input := &s3.DeleteObjectInput{Bucket: aws.String(bucketName), Key: aws.String(filename)}
	resp, err := s.DeleteObject(input)
	if err != nil {
		panic(err)
	}
	return resp
}

func uploadObject(s *s3.S3, filename string) (resp *s3.PutObjectOutput) {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	fmt.Println("uploading", filename)
	input := &s3.PutObjectInput{
		Body:   f,
		Bucket: aws.String(bucketName),
		Key:    aws.String(strings.Split(filename, "/")[1]),
		ACL:    aws.String(s3.BucketCannedACLPrivate),
	}
	resp, err = s.PutObject(input)
	if err != nil {
		panic(err)
	}
	return resp
}

func listObjects(s *s3.S3) *s3.ListObjectsV2Output {
	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
	}
	resp, err := s.ListObjectsV2(input)

	if err != nil {
		panic(err)
	}

	return resp
}
