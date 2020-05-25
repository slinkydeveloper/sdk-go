package http

import (
	"io"
	"mime"
	"mime/multipart"
	"net/http"

	"github.com/cloudevents/sdk-go/v2/binding"
)

const (
	MultipartCloudEvents = "multipart/cloudevents"
	JsonSeqCloudEvents   = "application/cloudevents-stream"
)

type MultipartMultiMessage struct {
	reader     *multipart.Reader
	bodyCloser io.Closer
}

type JsonStreamingMultiMessage struct {
	reader     *JsonSeqReader
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
		return &MultipartMultiMessage{
			reader:     multipart.NewReader(req.Body, param["boundary"]),
			bodyCloser: req.Body,
		}, nil
	}
	if contentType == JsonSeqCloudEvents {
		return &JsonStreamingMultiMessage{
			reader:     NewJsonSeqReader(req.Body),
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
		return &MultipartMultiMessage{
			reader:     multipart.NewReader(res.Body, param["boundary"]),
			bodyCloser: res.Body,
		}, nil
	}
	if contentType == JsonSeqCloudEvents {
		return &JsonStreamingMultiMessage{
			reader:     NewJsonSeqReader(res.Body),
			bodyCloser: res.Body,
		}, nil
	}
	return nil, binding.ErrUnknownEncoding
}

func (m *MultipartMultiMessage) Read() (binding.Message, error) {
	p, err := m.reader.NextPart()
	if err != nil {
		return nil, err
	}
	return NewMessage(http.Header(p.Header), p), nil
}

func (m *MultipartMultiMessage) Finish(error) error {
	return m.bodyCloser.Close()
}

func (j *JsonStreamingMultiMessage) Read() (binding.Message, error) {
	return j.reader.ReadNext()
}

func (j *JsonStreamingMultiMessage) Finish(err error) error {
	return j.bodyCloser.Close()
}
