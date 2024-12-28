// this package is used to manage mongodb connections and perform CRUD operations
package dbcontext

import (
	"sync"

	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// create cache to store db connections avoif multiple connections to same db
var cnn_cache = make(map[string]*DBContext)
var db_cache = make(map[string]*DB)

// DbAccess interface to perform CRUD operations
type DB struct {
	Client *mongo.Client
	DBName string
}

type DBContext struct {
	UriCnn string
	Client *mongo.Client
}
type AggregateStates[T any] struct {
}

func NewDBContext(uriCnn string) (*DBContext, error) {
	// check if connection already exists in cache
	if _, ok := cnn_cache[uriCnn]; ok {
		return cnn_cache[uriCnn], nil
	}
	// lock to avoid multiple connections to same db
	var lock sync.Mutex
	lock.Lock()
	defer lock.Unlock()
	// check again if connection already exists in cache
	if _, ok := cnn_cache[uriCnn]; ok {
		return cnn_cache[uriCnn], nil
	}
	// create new connection
	clientOptions := options.Client().ApplyURI(uriCnn)

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	// check if connection is valid
	err = client.Ping(ctx, nil)
	if err != nil {
		panic(err)
	}
	// add connection to cache
	cnn_cache[uriCnn] = &DBContext{UriCnn: uriCnn, Client: client}
	return cnn_cache[uriCnn], nil
}
func (db *DBContext) GetDB(dbName string) *DB {
	//change to check if db already exists in cache
	if _, ok := db_cache[dbName]; ok {
		return db_cache[dbName]
	}
	// lock to avoid multiple connections to same db
	var lock sync.Mutex
	lock.Lock()
	defer lock.Unlock()
	// check again if db already exists in cache
	if _, ok := db_cache[dbName]; ok {
		return db_cache[dbName]
	}
	// create new db
	db_cache[dbName] = &DB{Client: db.Client, DBName: dbName}
	return db_cache[dbName]
}
func FindOneToDict(db *DB, filter string) (map[string]interface{}, error) {

	panic("not implemented")
}
func InsertOneByDict(db *DB, data map[string]interface{}) error {
	panic("not implemented")
}
func UpdateOneByDict(db *DB, filter string, data map[string]interface{}) error {
	return nil
}

// implementation
func FindOne[T any](db DB, filter string) (T, error) {
	panic("not implemented")
}
