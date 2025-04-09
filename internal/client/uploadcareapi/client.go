package uploadcareapi

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/OutOfStack/game-library/internal/appconf"
	"github.com/OutOfStack/game-library/internal/pkg/observability"
	"github.com/uploadcare/uploadcare-go/ucare"
	"github.com/uploadcare/uploadcare-go/upload"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

const (
	uploadcareCDNURL = "https://ucarecdn.com/"
)

var tracer = otel.Tracer("uploadcareapi")

// Client represents dependencies for uploadcare client
type Client struct {
	log    *zap.Logger
	upload upload.Service
}

// New constructs Client instance
func New(log *zap.Logger, conf appconf.Uploadcare) (*Client, error) {
	creds := ucare.APICreds{
		SecretKey: conf.SecretKey,
		PublicKey: conf.PublicKey,
	}

	client := &http.Client{
		Transport: observability.NewMonitoredTransport(otelhttp.NewTransport(http.DefaultTransport), "uploadcare"),
	}

	ucareClient, err := ucare.NewClient(creds, &ucare.Config{
		SignBasedAuthentication: true,
		HTTPClient:              client,
	})
	if err != nil {
		return nil, fmt.Errorf("creating uploadcare client: %v", err)
	}

	return &Client{
		log:    log,
		upload: upload.NewService(ucareClient),
	}, nil
}

// UploadImage - uploads image and returns new image url
func (c *Client) UploadImage(ctx context.Context, data io.ReadSeeker, fileName string) (string, error) {
	fCtx, fSpan := tracer.Start(ctx, "uploadFile", trace.WithAttributes(attribute.String("filename", fileName)))
	defer fSpan.End()

	fileID, err := c.upload.File(fCtx, upload.FileParams{
		Data:    data,
		Name:    fileName,
		ToStore: ucare.String(upload.ToStoreTrue),
	})
	if err != nil {
		return "", fmt.Errorf("upload image to ucare: %v", err)
	}

	return getFileURL(fileID)
}

func getFileURL(fileID string) (string, error) {
	s, err := url.JoinPath(uploadcareCDNURL, fileID, "/")
	if err != nil {
		return "", fmt.Errorf("get uploadcare file url for fileID %s: %v", fileID, err)
	}
	return s, nil
}
