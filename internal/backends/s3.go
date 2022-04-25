package backends

import (
	"bytes"
	"net/url"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/enmand/cf-tf-diff/internal/terraform/state"
	"github.com/jbowes/cling"
)

type S3 struct {
	client *s3.S3
	dl     *s3manager.Downloader
}

func NewS3Backend(cfg map[string]interface{}) (Backend, error) {
	sess := session.Must(session.NewSession()) // todo: pass in config

	client := s3.New(sess)
	dl := s3manager.NewDownloader(sess)
	return &S3{
		client: client,
		dl:     dl,
	}, nil
}

func (b *S3) GetStateFile() (*state.State, error) {
	s := []byte{}
	u := url.URL{
		Host: "todo",
		Path: "todo",
	}

	params := &s3.GetObjectInput{
		Bucket: aws.String(u.Host),
		Key:    aws.String(u.Path),
	}
	_, err := b.dl.Download(aws.NewWriteAtBuffer(s), params)
	if err != nil {
		return nil, cling.Wrap(err, "unable to download state file")
	}

	return state.ReadState(&bytes.Buffer{})
}
