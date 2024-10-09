package filemanagement

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"server/src/helper"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Client struct {
	client *s3.Client
	bucket string
}

func New() (*S3Client, error) {

	accessKey := helper.GetEnvVar("AWS_ACCESS_KEY")
	secretKey := helper.GetEnvVar("AWS_SECRET_KEY")

	client := s3.New(s3.Options{
		Region:      "eu-central-1",
		Credentials: aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
	})

	return &S3Client{client: client, bucket: "luna-s3-bucket"}, nil
}
func (c *S3Client) UploadFile(file *os.File, key string) error {
	file.Seek(0, 0)
	_, err := c.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
		Body:   file,
	})
	if err != nil {
		return err
	}
	defer file.Close()
	return nil
}
func (c *S3Client) DownloadFile(key string) (*os.File, error) {
	tempFile, err := ioutil.TempFile("", "*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary file, %v", err)
	}
	fmt.Println("S3 object key:" + key)
	output, err := c.client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: &c.bucket,
		Key:    &key,
	})
	fmt.Println("S3 object size:", *output.ContentLength)
	if err != nil {
		return nil, fmt.Errorf("failed to download file, %v", err)
	}
	_, err = io.Copy(tempFile, output.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to copy file content, %v", err)
	}
	_, err = tempFile.Seek(0, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to seek file, %v", err)
	}

	return tempFile, nil
}
func (c *S3Client) DeleteFile(key string) error {
	_, err := c.client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: &c.bucket,
		Key:    &key,
	})
	if err != nil {
		return fmt.Errorf("failed to delete file, %v", err)
	}
	return nil
}
