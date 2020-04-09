package http

import (
	"io"
	"mime"
	"mime/multipart"
	"net/http"

	"github.com/cloudevents/sdk-go/v2/binding"
)

type MultiMessage struct {
	reader     *multipart.Reader
	bodyCloser io.Closer
}

func NewMultiMessageFromHttpRequest(req *http.Request) (binding.MultiMessage, error) {
	ct := req.Header.Get(ContentType)
	if ct == "" {
		return nil, binding.ErrUnknownEncoding
	}
	contentType, param, err := mime.ParseMediaType(ct)
	if err != nil {
		return nil, err
	}
	if contentType == MultipartCloudEvents {
		return &MultiMessage{
			reader:     multipart.NewReader(req.Body, param["boundary"]),
			bodyCloser: req.Body,
		}, nil
	}
	return nil, binding.ErrUnknownEncoding
}

func NewMultiMessageFromHttpResponse(res *http.Response) (binding.MultiMessage, error) {
	ct := res.Header.Get(ContentType)
	if ct == "" {
		return nil, binding.ErrUnknownEncoding
	}
	contentType, param, err := mime.ParseMediaType(ct)
	if err != nil {
		return nil, err
	}
	if contentType == MultipartCloudEvents {
		return &MultiMessage{
			reader:     multipart.NewReader(res.Body, param["boundary"]),
			bodyCloser: res.Body,
		}, nil
	}
	return nil, binding.ErrUnknownEncoding
}

func (m *MultiMessage) Read() (binding.Message, error) {
	p, err := m.reader.NextPart()
	if err != nil {
		return nil, err
	}
	return NewMessage(http.Header(p.Header), p), nil
}

func (m *MultiMessage) Finish(error) error {
	return m.bodyCloser.Close()
}
