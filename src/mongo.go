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
	CollectionUsers         string = "users"
	CollectionSessions      string = "sessions"
	CollectionApplications  string = "applications"
	CollectionServers       string = "servers"
	CollectionUniqueServers string = "unique_servers"
	CollectionTokens        string = "tokens"
	CollectionRequestLog    string = "request_log"
)

type MongoDB struct {
	Client   *mongo.Client
	Database *mongo.Database
}

type User struct {
	ID        string    `bson:"_id" json:"id"`
	Email     string    `bson:"email" json:"email"`
	Password  string    `bson:"password" json:"-"`
	Type      string    `bson:"type" json:"type"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
}

type Session struct {
	ID        string    `bson:"_id" json:"id"`
	User      string    `bson:"user" json:"user"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
}

type Application struct {
	ID               string    `bson:"_id" json:"id"`
	Name             string    `bson:"name" json:"name"`
	ShortDescription string    `bson:"shortDescription" json:"shortDescription"`
	User             string    `bson:"user" json:"user"`
	Token            string    `bson:"token" json:"token"`
	RequestCount     uint64    `bson:"requestCount" json:"requestCount"`
	CreatedAt        time.Time `bson:"createdAt" json:"createdAt"`
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

type Server struct {
	ID        string    `bson:"_id" json:"id"`
	Name      string    `bson:"name" json:"name"`
	User      string    `bson:"user" json:"user"`
	ServerID  string    `bson:"serverID" json:"serverID"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
}

type FormattedServer struct {
	ID        string        `bson:"_id" json:"id"`
	Name      string        `bson:"name" json:"name"`
	User      string        `bson:"user" json:"user"`
	Server    *UniqueServer `bson:"server" json:"server"`
	CreatedAt time.Time     `bson:"createdAt" json:"createdAt"`
}

type Token struct {
	ID           string     `bson:"_id" json:"id"`
	Name         string     `bson:"name" json:"name"`
	Token        string     `bson:"token" json:"token"`
	RequestCount uint64     `bson:"requestCount" json:"requestCount"`
	Application  string     `bson:"application" json:"application"`
	CreatedAt    time.Time  `bson:"createdAt" json:"createdAt"`
	LastUsedAt   *time.Time `bson:"lastUsedAt" json:"lastUsedAt"`
}

type RequestLog struct {
	ID           string    `bson:"_id" json:"_id"`
	Application  string    `bson:"application" json:"application"`
	Token        string    `bson:"token" json:"token"`
	Timestamp    time.Time `bson:"timestamp" json:"timestamp"`
	RequestCount int64     `bson:"requestCount" json:"requestCount"`
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

func (c *MongoDB) GetNextUniqueServers() ([]*UniqueServer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)

	defer cancel()

	cur, err := c.Database.Collection(CollectionUniqueServers).Aggregate(ctx, []bson.M{
		{
			"$match": bson.M{
				"lastPingedAt": bson.M{"$lte": time.Now().Add(-config.CycleInterval)},
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
