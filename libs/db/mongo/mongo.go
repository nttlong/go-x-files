package mongo

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	clientMu sync.Mutex
	cache    = make(map[string]*mongo.Client)
)

// GetClient retrieves a mongo.Client from the cache or creates a new one if not present.
func GetClient(uri string) *mongo.Client {
	if uri == "" {
		panic(fmt.Errorf("uri cannot be empty"))
	}

	clientMu.Lock()
	defer clientMu.Unlock()

	// Check if the client already exists in cache
	if client, ok := cache[uri]; ok {
		return client // Client already exists in cache
	}

	// Create a new context with a timeout for connection attempts
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		panic(fmt.Errorf("failed to create client for uri '%s': %w", uri, err))
	}

	// Optionally, you can check the connection status here
	// if err := client.Ping(ctx, nil); err != nil {
	//    return nil, fmt.Errorf("failed to ping client: %w", err)
	// }

	cache[uri] = client
	return client
}

// Optional: A function to disconnect and remove a client from the cache
func DisconnectClient(uri string) error {
	clientMu.Lock()
	defer clientMu.Unlock()

	if client, ok := cache[uri]; ok {
		if err := client.Disconnect(context.TODO()); err != nil {
			return fmt.Errorf("failed to disconnect client for uri '%s': %w", uri, err)
		}
		delete(cache, uri) // Remove from cache after disconnecting
	}
	return nil
}

func InsertOne(client *mongo.Client, dbName, collectionName string, document interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	coll := client.Database(dbName).Collection(collectionName)
	result, err := coll.InsertOne(ctx, document)
	if err != nil {
		return fmt.Errorf("failed to insert document: %w", err)
	}
	fmt.Println("Inserted document with ID:", result.InsertedID)
	return nil
}
func GetAllTags(t reflect.Type) (map[string]string, error) {
	if t.Kind() == reflect.Ptr {
		t = t.Elem() // Dereference pointer
	}

	if t.Kind() != reflect.Struct {
		return nil, fmt.Errorf("input must be a struct or a pointer to a struct, but got %v", t.Kind())
	}

	tags := make(map[string]string)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := string(field.Tag) // Get the entire tag string

		if tag != "" {
			tags[field.Name] = tag
		}

		// Handle embedded structs recursively
		if field.Type.Kind() == reflect.Struct && field.Anonymous {
			embeddedTags, err := GetAllTags(field.Type)
			if err != nil {
				return nil, fmt.Errorf("error getting tags from embedded struct %s: %w", field.Name, err)
			}
			for k, v := range embeddedTags {
				tags[k] = v
			}
		}
	}

	return tags, nil
}
