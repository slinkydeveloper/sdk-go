package main

import (
	"context"
	"log"
	nethttp "net/http"

	nethttp2 "golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"strconv"

	"github.com/google/uuid"

	"github.com/cloudevents/sdk-go/v2/binding"
	"github.com/cloudevents/sdk-go/v2/binding/test"
	"github.com/cloudevents/sdk-go/v2/protocol/http"

	cloudevents "github.com/cloudevents/sdk-go/v2"
)

func main() {
	// Setup HTTP 2 upgrade
	handler := nethttp.HandlerFunc(ServeHttp)
	h2s := &nethttp2.Server{}
	h1s := &nethttp.Server{
		Addr:    ":8080",
		Handler: h2c.NewHandler(handler, h2s),
	}

	log.Printf("will listen on :8080\n")
	if err := h1s.ListenAndServe(); err != nil {
		log.Fatalf("unable to start http server, %s", err)
	}
}

// Example is a basic data struct.
type Example struct {
	Sequence int    `json:"id"`
	Message  string `json:"message"`
}

func ServeHttp(responseWriter nethttp.ResponseWriter, request *nethttp.Request) {
	// Get number of events i need to generate
	numParam := request.URL.Query().Get("num")
	num := 5
	if numParam != "" {
		num, _ = strconv.Atoi(numParam)
	}

	log.Printf("will emit %d events\n", num)

	// Generate some messages
	messages := make([]binding.Message, num)
	for i := range messages {
		event := cloudevents.NewEvent()
		event.SetID(uuid.New().String())
		event.SetSource("https://github.com/cloudevents/sdk-go/v2/cmd/samples/http_multipart/random_multi_event_response")
		event.SetType("com.cloudevents.multipart.sample")
		_ = event.SetData("application/json", &Example{
			Sequence: i,
			Message:  "hello world",
		})

		// This is to show mixed binary and structured representation in the same envelope
		if i%2 != 0 {
			messages[i] = test.MustCreateMockBinaryMessage(event)
		} else {
			messages[i] = test.MustCreateMockStructuredMessage(event)
		}
	}

	// Write the multipart response!
	multiMessage := binding.NewGenericMultiMessage(messages...)
	err := http.WriteJsonSeqResponse(context.Background(), multiMessage, 200, responseWriter)
	if err != nil {
		log.Fatal(err)
	}
}
