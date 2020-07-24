package stub

import (
	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/pubsub/pstest"
	"context"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
)

func CreateTestPubSubServer(topicId string, ctx context.Context) (*grpc.ClientConn, *pubsub.Client) {
	// Start a fake server running locally.
	srv := pstest.NewServer()
	// Connect to the server without using TLS.
	conn, _ := grpc.Dial(srv.Addr, grpc.WithInsecure())
	// Use the connection when creating a pubsub client.
	client, _ := pubsub.NewClient(ctx, "project", option.WithGRPCConn(conn))
	topic, _ := client.CreateTopic(ctx, topicId)
	_ = topic
	return conn, client
}