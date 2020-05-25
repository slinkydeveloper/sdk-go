module github.com/cloudevents/sdk-go/v2/samples

go 1.13

replace github.com/cloudevents/sdk-go/v2 => ../../v2

replace github.com/cloudevents/sdk-go/v2/protocol/pubsub => ../../v2/protocol/pubsub

replace github.com/cloudevents/sdk-go/v2/protocol/amqp => ../../v2/protocol/amqp

replace github.com/cloudevents/sdk-go/v2/protocol/stan => ../../v2/protocol/stan

replace github.com/cloudevents/sdk-go/v2/protocol/nats => ../../v2/protocol/nats

replace github.com/cloudevents/sdk-go/v2/protocol/kafka_sarama => ../../v2/protocol/kafka_sarama

require (
	contrib.go.opencensus.io/exporter/prometheus v0.1.0
	github.com/Azure/go-amqp v0.12.7
	github.com/Shopify/sarama v1.19.0
	github.com/cloudevents/sdk-go/v2 v2.0.0-00010101000000-000000000000
	github.com/cloudevents/sdk-go/v2/protocol/amqp v0.0.0-00010101000000-000000000000
	github.com/cloudevents/sdk-go/v2/protocol/kafka_sarama v0.0.0-00010101000000-000000000000
	github.com/cloudevents/sdk-go/v2/protocol/nats v0.0.0-00010101000000-000000000000
	github.com/cloudevents/sdk-go/v2/protocol/pubsub v0.0.0-00010101000000-000000000000
	github.com/cloudevents/sdk-go/v2/protocol/stan v0.0.0-00010101000000-000000000000
	github.com/google/uuid v1.1.1
	github.com/gorilla/mux v1.7.3
	github.com/kelseyhightower/envconfig v1.4.0
	go.opencensus.io v0.22.3
	golang.org/x/net v0.0.0-20200324143707-d3edc9973b7e
)
