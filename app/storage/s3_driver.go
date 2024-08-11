package storage

import (
	"io"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3Driver struct {
	Client *s3.S3
	Bucket string
}

func NewS3Driver(bucket string) (*S3Driver, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2"),
	})
	if err != nil {
		return nil, err
	}
	return &S3Driver{
		Client: s3.New(sess),
		Bucket: bucket,
	}, nil
}

func (s3d *S3Driver) Put(path string, content io.ReadSeeker) error {
	_, err := s3d.Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(s3d.Bucket),
		Key:    aws.String(path),
		Body:   content,
	})
	return err
}

func (s3d *S3Driver) Get(path string) (io.Reader, error) {
	result, err := s3d.Client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(s3d.Bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		return nil, err
	}
	return result.Body, nil
}

func (s3d *S3Driver) Delete(path string) error {
	_, err := s3d.Client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(s3d.Bucket),
		Key:    aws.String(path),
	})
	return err
}

func (s3d *S3Driver) Exists(path string) (bool, error) {
	_, err := s3d.Client.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(s3d.Bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s3d *S3Driver) List(directory string) ([]string, error) {
	params := &s3.ListObjectsInput{
		Bucket: aws.String(s3d.Bucket),
		Prefix: aws.String(directory),
	}

	resp, err := s3d.Client.ListObjects(params)
	if err != nil {
		return nil, err
	}

	var files []string
	for _, item := range resp.Contents {
		files = append(files, *item.Key)
	}

	return files, nil
}

func (s3d *S3Driver) Missing(path string) (bool, error) {
	exists, err := s3d.Exists(path)
	if err != nil {
		return false, err
	}
	return !exists, nil
}

func (s3d *S3Driver) Download(path string) (io.Reader, error) {
	return s3d.Get(path)
}

func (s3d *S3Driver) URL(path string) (string, error) {
	req, _ := s3d.Client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(s3d.Bucket),
		Key:    aws.String(path),
	})
	return req.Presign(15 * time.Minute)
}

func (s3d *S3Driver) TemporaryURL(path string, expiresIn int64) (string, error) {
	req, _ := s3d.Client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(s3d.Bucket),
		Key:    aws.String(path),
	})
	return req.Presign(time.Duration(expiresIn) * time.Second)
}

func (s3d *S3Driver) Size(path string) (int64, error) {
	output, err := s3d.Client.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(s3d.Bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		return 0, err
	}
	return *output.ContentLength, nil
}

func (s3d *S3Driver) Copy(sourcePath string, destinationPath string) error {
	_, err := s3d.Client.CopyObject(&s3.CopyObjectInput{
		CopySource: aws.String(s3d.Bucket + "/" + sourcePath),
		Bucket:     aws.String(s3d.Bucket),
		Key:        aws.String(destinationPath),
	})
	return err
}

func (s3d *S3Driver) Move(sourcePath string, destinationPath string) error {
	err := s3d.Copy(sourcePath, destinationPath)
	if err != nil {
		return err
	}
	return s3d.Delete(sourcePath)
}
