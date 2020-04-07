package transformer

import (
	"github.com/cloudevents/sdk-go/v2/binding"
	"github.com/cloudevents/sdk-go/v2/event"
)

// ExtractExtensions is a TransformerFactory which extracts a set of extensions from transformed mesages. An instance of this transformer should only be used to process a single message, after which the extracted extensions can be read from the map. All extension which are present as keys in the map will be extracted and other extensions will be ignored.
type ExtractExtensions map[string]interface{}

func (e ExtractExtensions) StructuredTransformer(_ binding.StructuredWriter) binding.StructuredWriter {
	return nil
}

func (e ExtractExtensions) BinaryTransformer(writer binding.BinaryWriter) binding.BinaryWriter {
	return extractExtensionsWriter{
		BinaryWriter: writer,
		extensions:   e,
	}
}

func (e ExtractExtensions) EventTransformer() binding.EventTransformer {
	return func(event *event.Event) error {
		for name := range e {
			e[name] = event.Extensions()[name]
		}
		return nil
	}
}

type extractExtensionsWriter struct {
	binding.BinaryWriter
	extensions ExtractExtensions
}

func (e extractExtensionsWriter) SetExtension(name string, value interface{}) error {
	if _, ok := e.extensions[name]; ok {
		e.extensions[name] = value
	}
	return e.BinaryWriter.SetExtension(name, value)
}

// ExtractExtensions is a TransformerFactory which extracts a set of extensions from transformed mesages. An instance of this transformer should only be used to process a single message, after which the extracted extensions can be read from the map. All extension which are present as keys in the map will be extracted and other extensions will be ignored.
type ExtractExtension struct {
	Name  string
	Value interface{}
}

func (e *ExtractExtension) StructuredTransformer(_ binding.StructuredWriter) binding.StructuredWriter {
	return nil
}

func (e *ExtractExtension) BinaryTransformer(writer binding.BinaryWriter) binding.BinaryWriter {
	return extractExtensionWriter{
		BinaryWriter: writer,
		extension:    e,
	}
}

func (e *ExtractExtension) EventTransformer() binding.EventTransformer {
	return func(event *event.Event) error {
		e.Value = event.Extensions()[e.Name]
		return nil
	}
}

type extractExtensionWriter struct {
	binding.BinaryWriter
	extension *ExtractExtension
}

func (e extractExtensionWriter) SetExtension(name string, value interface{}) error {
	if name == e.extension.Name {
		e.extension.Value = value
	}
	return e.BinaryWriter.SetExtension(name, value)
}
