package http

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/cloudevents/sdk-go/v2/binding"
)

func WriteMultipartRequest(ctx context.Context, m binding.MultiMessage, httpRequest *http.Request, transformers ...binding.Transformer) error {
	var buf bytes.Buffer // We can pool it!
	writer := multipart.NewWriter(&buf)

	// Let's first write the content type and the 200 status code
	httpRequest.Header.Set(ContentType, multipartContentType(writer.Boundary()))

	// Now we can start writing the body
	singleMessage, err := m.Read()
	for err != io.EOF {
		if err != nil {
			return err
		}

		// Now I have a single message and i can reuse the usual write
		pw := partWriter{mw: writer}
		_, writeErr := binding.Write(ctx, singleMessage, &pw, &pw, transformers...)
		if writeErr != nil {
			return writeErr
		}
		singleMessage, err = m.Read()
	}
	err = writer.Close()
	if err != nil {
		return err
	}
	return setRequestBody(httpRequest, &buf)
}

func WriteJsonSeqRequest(ctx context.Context, m binding.MultiMessage, httpRequest *http.Request, transformers ...binding.Transformer) error {
	var buf bytes.Buffer // We can pool it!

	// Let's first write the content type and the 200 status code
	httpRequest.Header.Set(ContentType, JsonSeqCloudEvents)

	// Now we can start writing the body
	singleMessage, err := m.Read()
	for err != io.EOF {
		if err != nil {
			return err
		}

		// Now I have a single message and i can reuse the usual write
		jsonSeqWriter := JsonSeqWriter{writer: &buf}
		_, writeErr := binding.Write(ctx, singleMessage, &jsonSeqWriter, nil, transformers...)
		if writeErr != nil {
			return writeErr
		}
		singleMessage, err = m.Read()
	}
	return setRequestBody(httpRequest, &buf)
}
