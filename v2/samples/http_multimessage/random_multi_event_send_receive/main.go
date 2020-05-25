package main

import (
	"context"
	"fmt"
	"io"
	"log"
	nethttp "net/http"

	"github.com/google/uuid"
	nethttp2 "golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/binding"
	"github.com/cloudevents/sdk-go/v2/binding/transformer"
	"github.com/cloudevents/sdk-go/v2/protocol/http"
)

// Example is a basic data struct.
type Example struct {
	Sequence int    `json:"id"`
	Message  string `json:"message"`
}

func main() {
	// Setup HTTP 2 upgrade
	handler := nethttp.HandlerFunc(ServeHttp)
	h2s := &nethttp2.Server{}
	h1s := &nethttp.Server{
		Addr:    ":8080",
		Handler: h2c.NewHandler(handler, h2s),
	}

	// Start HTTP server
	go func() {
		log.Printf("will listen on :8080\n")
		if err := h1s.ListenAndServe(); err != nil {
			log.Fatalf("unable to start http server, %s", err)
		}
	}()

	// Generate some events
	events := make([]cloudevents.Event, 5)
	for i := range events {
		event := cloudevents.NewEvent()
		event.SetID(uuid.New().String())
		event.SetSource("https://github.com/cloudevents/sdk-go/v2/cmd/samples/http_multipart/random_multi_event_response")
		event.SetType("com.cloudevents.multipart.sample")
		_ = event.SetData("application/json", &Example{
			Sequence: i,
			Message:  "hello world",
		})

		fmt.Printf("Send Event %d: %v\n", i, event)

		events[i] = event
	}

	client := nethttp.Client{}
	req, _ := nethttp.NewRequest("POST", "http://localhost:8080", nil)

	// Write the multipart request!
	multiMessage := binding.NewEventMultiMessage(events...)
	err := http.WriteMultipartRequest(context.Background(), multiMessage, req)
	if err != nil {
		log.Fatal(err)
	}

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\n--- Status code %d, Status Message: %s ---\n\n", res.StatusCode, res.Status)

	if res.StatusCode != 200 {
		return
	}

	responseMessage, err := http.NewMultiMessageFromHttpResponse(res)
	if err != nil {
		log.Fatal(err)
	}

	singleMessage, err := responseMessage.Read()
	for err != io.EOF {
		if err != nil {
			log.Fatal(err)
		}

		// Now I have a single message and i can reuse the usual write
		event, readErr := binding.ToEvent(context.TODO(), singleMessage)
		if readErr != nil {
			log.Fatal(readErr)
		}
		fmt.Printf("Received Event: %v\n", event)
		singleMessage, err = responseMessage.Read()
	}
}

func ServeHttp(responseWriter nethttp.ResponseWriter, request *nethttp.Request) {
	// Pipe in and out and apply transformations
	multiMessage, err := http.NewMultiMessageFromHttpRequest(request)
	if err != nil {
		log.Print(err)
		nethttp.Error(responseWriter, err.Error(), 400)
		return
	}
	defer multiMessage.Finish(nil)
	err = http.WriteMultipartResponse(
		context.Background(),
		multiMessage,
		200,
		responseWriter,
		transformer.AddExtension("proxy", "sdk-go"),
		transformer.AddTimeNow,
	)
	if err != nil {
		log.Print(err)
		nethttp.Error(responseWriter, err.Error(), 500)
	}
}
