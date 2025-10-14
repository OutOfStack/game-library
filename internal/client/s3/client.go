package s3

import (
	"context"
	"fmt"
	"io"
	"mime"
	"net/http"
	"time"

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
	"go.uber.org/zap"
)

const (
	defaultTimeout = 10 * time.Second
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
		Timeout:   defaultTimeout,
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
func (c *Client) Upload(ctx context.Context, data io.ReadSeeker, contentType string, md map[string]string) (UploadResult, error) {
	ctx, span := tracer.Start(ctx, "upload")
	defer span.End()

	// get file extension and construct object key
	var fileContentType *string
	fileExt, err := getExtensionByContentType(contentType)
	if err != nil {
		c.log.Warn("detect content type", zap.String("type", contentType), zap.Error(err))
	} else {
		fileContentType = aws.String(contentType)
	}
	objectKey := uuid.NewString() + fileExt

	span.SetAttributes(attribute.String("objectKey", objectKey))

	_, err = data.Seek(0, io.SeekStart)
	if err != nil {
		return UploadResult{}, fmt.Errorf("seek file start: %v", err)
	}

	_, err = c.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(c.bucketName),
		Key:         aws.String(objectKey),
		Body:        data,
		ContentType: fileContentType,
		Metadata:    md,
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

var contentTypeExtensionOverrides = map[string]string{
	"image/jpeg": ".jpg",
}

func getExtensionByContentType(contentType string) (string, error) {
	if ext, ok := contentTypeExtensionOverrides[contentType]; ok {
		return ext, nil
	}
	exts, err := mime.ExtensionsByType(contentType)
	if err != nil {
		return "", err
	}
	if len(exts) > 0 {
		return exts[0], nil
	}
	return "", nil
}
