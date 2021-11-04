package s3

import (
	"bytes"
	"email_action/logging"
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/s3/s3manager/s3manageriface"
)

var (
	log               = logging.NewZapLogger()
	ErrFileInvalid    = errors.New("file is invalid")
	ErrFileKeyInvalid = errors.New("file key is invalid")
)

type S3Adapter struct {
	s3Svc            s3iface.S3API
	uploader         s3manageriface.UploaderAPI
	downloader       s3manageriface.DownloaderAPI
	fileUploadBucket string
}

func NewS3Adapter(env string) *S3Adapter {
	log.Infof("NewS3Adapter(): env=%s", env)
	awsConfig := &aws.Config{
		Region: aws.String("us-west-2"),
	}
	sess := session.Must(session.NewSession(awsConfig))
	svc := s3.New(sess)
	uploader := s3manager.NewUploaderWithClient(svc)
	downloader := s3manager.NewDownloaderWithClient(svc)

	return &S3Adapter{
		s3Svc:            svc,
		uploader:         uploader,
		downloader:       downloader,
		fileUploadBucket: "email-action-file-upload-" + env,
	}
}

func (h *S3Adapter) UploadFile(file []byte, key string) (string, error) {
	log.Infof("UploadFile(): %s", key)
	if file == nil {
		return "", ErrFileInvalid
	}
	if key == "" {
		return "", ErrFileKeyInvalid
	}

	reader := bytes.NewReader(file)
	upParams := &s3manager.UploadInput{
		Bucket: &h.fileUploadBucket,
		Key:    aws.String(key),
		Body:   reader,
	}

	s3Output, err := h.uploader.Upload(upParams)
	if err != nil {
		return "", err
	}
	return s3Output.Location, nil
}

func (h *S3Adapter) DownloadFile(key string) ([]byte, error) {
	log.Infof("DownloadFile(): %s", key)
	if key == "" {
		return []byte{}, ErrFileKeyInvalid
	}

	req := &s3.GetObjectInput{
		Bucket: &h.fileUploadBucket,
		Key:    aws.String(key),
	}

	buff := &aws.WriteAtBuffer{}

	_, err := h.downloader.Download(buff, req)
	if err != nil {
		return nil, err
	}
	return buff.Bytes(), nil
}

func (h *S3Adapter) DeleteFile(key string) error {
	log.Infof("DeleteFile(): %s", key)
	if key == "" {
		return ErrFileKeyInvalid
	}

	req := &s3.DeleteObjectInput{
		Bucket: &h.fileUploadBucket,
		Key:    aws.String(key),
	}

	_, err := h.s3Svc.DeleteObject(req)
	if err != nil {
		log.Errorf("DeleteFile(): %v", err)
		return err
	}
	return nil
}

func (h *S3Adapter) ListFiles(prefix string) ([]string, error) {
	log.Infof("ListFiles(): %s", prefix)
	if prefix == "" {
		return nil, ErrFileKeyInvalid
	}

	req := &s3.ListObjectsInput{
		Bucket: &h.fileUploadBucket,
		Prefix: aws.String(prefix + "/"),
	}

	output, err := h.s3Svc.ListObjects(req)
	if err != nil {
		log.Errorf("ListFiles(): %v", err)
		return nil, err
	}

	var urls []string
	for _, content := range output.Contents {
		urls = append(urls, h.compileS3Url(*content.Key))
	}
	return urls, nil
}

func (h *S3Adapter) compileS3Url(key string) string {
	return "https://" + h.fileUploadBucket + ".s3-us-west-2.amazonaws.com/" + key
}
