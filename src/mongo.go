package main

import (
	"context"
	"net/url"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	CollectionUniqueServers    string = "unique_servers"
	CollectionServerStatistics string = "server_statistics"
)

type MongoDB struct {
	Client   *mongo.Client
	Database *mongo.Database
}

type UniqueServer struct {
	ID          string    `bson:"_id" json:"id"`
	Type        string    `bson:"type" json:"type"`
	Status      string    `bson:"status" json:"status"`
	Hostname    string    `bson:"hostname" json:"hostname"`
	Port        uint16    `bson:"port" json:"port"`
	OnlineCount uint64    `bson:"onlineCount" json:"onlineCount"`
	TotalCount  uint64    `bson:"totalCount" json:"totalCount"`
	CreatedAt   time.Time `bson:"createdAt" json:"createdAt"`
}

func (c *MongoDB) Connect(uri string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)

	defer cancel()

	parsedURI, err := url.Parse(uri)

	if err != nil {
		return err
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))

	if err != nil {
		return err
	}

	if err := client.Ping(ctx, nil); err != nil {
		return err
	}

	c.Client = client
	c.Database = client.Database(strings.TrimPrefix(parsedURI.Path, "/"))

	return nil
}

func (c *MongoDB) UpsertServerStatistics(filter, update bson.M) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)

	defer cancel()

	_, err := c.Database.Collection(CollectionServerStatistics).UpdateOne(ctx, filter, update, &options.UpdateOptions{
		Upsert: PointerOf(true),
	})

	return err
}

func (c *MongoDB) GetNextUniqueServers() ([]*UniqueServer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)

	defer cancel()

	cur, err := c.Database.Collection(CollectionUniqueServers).Aggregate(ctx, []bson.M{
		{
			"$match": bson.M{
				// "lastPingedAt": bson.M{"$lte": time.Now().Add(-config.CycleInterval / 2)}, // TODO test if /2 is the solution
				// TODO chunk by instance ID
			},
		},
	})

	if err != nil {
		return nil, err
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	result := make([]*UniqueServer, 0)

	if err := cur.All(ctx, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (c *MongoDB) UpdateUniqueServerByID(id string, update bson.M) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)

	defer cancel()

	_, err := c.Database.Collection(CollectionUniqueServers).UpdateOne(ctx, bson.M{"_id": id}, update)

	return err
}

func (c *MongoDB) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)

	defer cancel()

	return c.Client.Disconnect(ctx)
}
