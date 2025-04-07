package s3

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/OutOfStack/game-library/internal/appconf"
	"github.com/OutOfStack/game-library/internal/pkg/observability"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

var tracer = otel.Tracer("s3")

// Client represents dependencies for S3 client
type Client struct {
	log        *zap.Logger
	s3Client   *s3.Client
	bucketName string
	cdnBaseURL string
}

// New constructs Client instance
func New(log *zap.Logger, conf appconf.S3) (*Client, error) {
	httpClient := &http.Client{
		Transport: observability.NewMonitoredTransport(otelhttp.NewTransport(http.DefaultTransport), "s3"),
	}

	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(conf.Region),
		config.WithBaseEndpoint(conf.Endpoint),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			conf.AccessKeyID,
			conf.SecretAccessKey,
			"",
		)),
		config.WithHTTPClient(httpClient),
	)
	if err != nil {
		return nil, fmt.Errorf("loading S3 config: %v", err)
	}

	s3Client := s3.NewFromConfig(cfg)

	return &Client{
		log:        log,
		s3Client:   s3Client,
		bucketName: conf.BucketName,
		cdnBaseURL: conf.CDNBaseURL,
	}, nil
}

// Upload uploads a file to S3 storage
func (c *Client) Upload(ctx context.Context, data io.ReadSeeker, fileName string) (UploadResult, error) {
	ext := strings.ToLower(filepath.Ext(fileName))
	objectKey := uuid.NewString()
	if ext != "" {
		objectKey += ext
	}

	ctx, span := tracer.Start(ctx, "upload", trace.WithAttributes(
		attribute.String("objectKey", objectKey),
		attribute.String("fileName", fileName),
	))
	defer span.End()

	if _, err := data.Seek(0, io.SeekStart); err != nil {
		return UploadResult{}, fmt.Errorf("seek file start: %v", err)
	}

	_, err := c.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(c.bucketName),
		Key:    aws.String(objectKey),
		Body:   data,
		Metadata: map[string]string{
			"fileName": fileName,
		},
	})
	if err != nil {
		return UploadResult{}, fmt.Errorf("putting object: %v", err)
	}

	fileURL := fmt.Sprintf("%s/%s", c.cdnBaseURL, objectKey)

	return UploadResult{
		FileID:  objectKey,
		FileURL: fileURL,
	}, nil
}
