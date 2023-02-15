package uploadcare

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"path"

	"github.com/OutOfStack/game-library/internal/appconf"
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

var tracer = otel.Tracer("")

// Client represents dependencies for uploadcare client
type Client struct {
	log *zap.Logger
	hc  *http.Client
	uc  *ucare.Client
}

// New constructs Client instance
func New(log *zap.Logger, conf appconf.Uploadcare) (*Client, error) {
	creds := ucare.APICreds{
		SecretKey: conf.SecretKey,
		PublicKey: conf.PublicKey,
	}

	otelClient := otelhttp.DefaultClient

	uClient, err := ucare.NewClient(creds, &ucare.Config{
		SignBasedAuthentication: true,
		HTTPClient:              otelClient,
	})
	if err != nil {
		return nil, fmt.Errorf("creating uploadcare client: %v", err)
	}

	return &Client{
		log: log,
		uc:  &uClient,
		hc:  otelClient,
	}, nil
}

// UploadImageFromURL - uploads image from image url and returns new image url
func (c *Client) UploadImageFromURL(ctx context.Context, imageURL string) (newURL string, err error) {
	ctx, span := tracer.Start(ctx, "uploadcare.uploadimagefromurl")
	defer span.End()

	request, err := http.NewRequestWithContext(ctx, "GET", imageURL, nil)
	if err != nil {
		return "", fmt.Errorf("creating get image by url request: %v", err)
	}

	resp, err := c.hc.Do(request)
	if err != nil {
		return "", fmt.Errorf("get image by url: %v", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response body: %v", err)
	}

	reader := bytes.NewReader(data)
	fileName := path.Base(request.URL.Path)

	service := upload.NewService(*c.uc)
	fileParams := upload.FileParams{
		Data:    reader,
		Name:    fileName,
		ToStore: ucare.String(upload.ToStoreTrue),
	}

	fCtx, fSpan := tracer.Start(ctx, "uploadcare/uploadcare-go/upload.file", trace.WithAttributes(attribute.String("filename", fileName)))
	fileID, err := service.File(fCtx, fileParams)
	if err != nil {
		fSpan.End()
		return "", fmt.Errorf("upload image to ucare: %v", err)
	}
	fSpan.End()

	return getFileURL(fileID), nil
}

func getFileURL(fileID string) string {
	return path.Join(uploadcareCDNURL, fileID) + "/"
}
