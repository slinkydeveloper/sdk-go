package http

import (
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"

	"github.com/cloudevents/sdk-go/v2/binding"
	"github.com/cloudevents/sdk-go/v2/binding/format"
	"github.com/cloudevents/sdk-go/v2/binding/spec"
	"github.com/cloudevents/sdk-go/v2/types"
)

func WriteMultipartResponse(ctx context.Context, m binding.MultiMessage, status int, rw http.ResponseWriter, transformers ...binding.Transformer) error {
	writer := multipart.NewWriter(rw)

	// Let's first write the content type and the 200 status code
	rw.Header().Set(ContentType, multipartContentType(writer.Boundary()))
	rw.WriteHeader(status)

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
	return writer.Close()
}

func WriteJsonSeqResponse(ctx context.Context, m binding.MultiMessage, status int, rw http.ResponseWriter, transformers ...binding.Transformer) error {
	// Let's first write the content type and the 200 status code
	rw.Header().Set(ContentType, JsonSeqCloudEvents)
	rw.WriteHeader(status)

	// Now we can start writing the body
	singleMessage, err := m.Read()
	for err != io.EOF {
		if err != nil {
			return err
		}

		// Now I have a single message and i can reuse the usual write
		jsonSeqWriter := JsonSeqWriter{writer: rw}
		_, writeErr := binding.Write(ctx, singleMessage, &jsonSeqWriter, nil, transformers...)
		if writeErr != nil {
			return writeErr
		}
		singleMessage, err = m.Read()
	}
	return nil
}

func multipartContentType(boundary string) string {
	return MultipartCloudEvents + "; boundary=" + boundary
}

type partWriter struct {
	mw      *multipart.Writer
	headers textproto.MIMEHeader
	body    io.Reader
}

func (b *partWriter) SetStructuredEvent(ctx context.Context, format format.Format, event io.Reader) error {
	b.headers = textproto.MIMEHeader{
		ContentType: []string{format.MediaType()},
	}
	b.body = event
	return b.finalizeWriter()
}

func (b *partWriter) Start(ctx context.Context) error {
	b.headers = make(textproto.MIMEHeader, 4)
	return nil
}

func (b *partWriter) SetAttribute(attribute spec.Attribute, value interface{}) error {
	mapping := attributeHeadersMapping[attribute.Name()]
	if value == nil {
		delete(b.headers, mapping)
	}

	// Http headers, everything is a string!
	s, err := types.Format(value)
	if err != nil {
		return err
	}
	b.headers[mapping] = append(b.headers[mapping], s)
	return nil
}

func (b *partWriter) SetExtension(name string, value interface{}) error {
	if value == nil {
		delete(b.headers, extNameToHeaderName(name))
	}
	// Http headers, everything is a string!
	s, err := types.Format(value)
	if err != nil {
		return err
	}
	b.headers[extNameToHeaderName(name)] = []string{s}
	return nil
}

func (b *partWriter) SetData(reader io.Reader) error {
	b.body = reader
	return nil
}

func (b *partWriter) End(ctx context.Context) error {
	return b.finalizeWriter()
}

func (b *partWriter) finalizeWriter() error {
	writer, err := b.mw.CreatePart(b.headers)
	if err != nil {
		return err
	}

	_, err = io.Copy(writer, b.body)

	return err
}

var _ binding.StructuredWriter = (*partWriter)(nil) // Test it conforms to the interface
var _ binding.BinaryWriter = (*partWriter)(nil)     // Test it conforms to the interface
