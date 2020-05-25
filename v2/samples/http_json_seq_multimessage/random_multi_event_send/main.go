package main

import (
	"context"
	"fmt"
	"log"
	nethttp "net/http"

	"github.com/google/uuid"

	"github.com/cloudevents/sdk-go/v2/binding"
	"github.com/cloudevents/sdk-go/v2/binding/test"
	"github.com/cloudevents/sdk-go/v2/protocol/http"

	cloudevents "github.com/cloudevents/sdk-go/v2"
)

// Example is a basic data struct.
type Example struct {
	Sequence int    `json:"id"`
	Message  string `json:"message"`
}

func main() {
	// Generate some messages
	messages := make([]binding.Message, 5)
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

	client := nethttp.Client{}
	req, _ := nethttp.NewRequest("POST", "http://localhost:8080", nil)

	// Write the multipart request!
	multiMessage := binding.NewGenericMultiMessage(messages...)
	err := http.WriteJsonSeqRequest(context.Background(), multiMessage, req)
	if err != nil {
		log.Fatal(err)
	}

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Received %d status code\n", res.StatusCode)

}
