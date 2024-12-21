package glogs

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// LogLevel represents the severity of a log message.
type LogLevel int

const (
	Trace LogLevel = iota
	Debug
	Info
	Warning
	Error
	Critical
)

var logLevelStrings = []string{
	"TRACE",
	"DEBUG",
	"INFO",
	"WARNING",
	"ERROR",
	"CRITICAL",
}

func (l LogLevel) String() string {
	if l < Trace || l > Critical {
		return "UNKNOWN"
	}
	return logLevelStrings[l]
}

// LogEntry represents a log message
type LogEntry struct {
	Timestamp time.Time              `bson:"timestamp"`
	Level     LogLevel               `bson:"level"`
	Message   string                 `bson:"message"`
	Context   map[string]interface{} `bson:"context,omitempty"` // Contextual information
}

// ILogger is the interface for logging messages.
type ILogger interface {
	Log(level LogLevel, message string, ctx map[string]interface{})
	Trace(message string, ctx map[string]interface{})
	Debug(message string, ctx map[string]interface{})
	Info(message string, ctx map[string]interface{})
	Warning(message string, ctx map[string]interface{})
	Error(message string, ctx map[string]interface{})
	Critical(message string, ctx map[string]interface{})
}

// FileLogger logs to a file.
type FileLogger struct {
	file *os.File
	log  *log.Logger
}

// NewFileLogger creates a new FileLogger.
func NewFileLogger(filename string) (ILogger, error) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	return &FileLogger{file: file, log: log.New(file, "", log.LstdFlags)}, nil
}

func (l *FileLogger) Log(level LogLevel, message string, ctx map[string]interface{}) {
	l.log.Printf("%s: %s %+v\n", level, message, ctx)
}

func (l *FileLogger) Trace(message string, ctx map[string]interface{}) { l.Log(Trace, message, ctx) }
func (l *FileLogger) Debug(message string, ctx map[string]interface{}) { l.Log(Debug, message, ctx) }
func (l *FileLogger) Info(message string, ctx map[string]interface{})  { l.Log(Info, message, ctx) }
func (l *FileLogger) Warning(message string, ctx map[string]interface{}) {
	l.Log(Warning, message, ctx)
}
func (l *FileLogger) Error(message string, ctx map[string]interface{}) { l.Log(Error, message, ctx) }
func (l *FileLogger) Critical(message string, ctx map[string]interface{}) {
	l.Log(Critical, message, ctx)
}

func (l *FileLogger) Close() error {
	return l.file.Close()
}

// MongoLogger logs to MongoDB.
type MongoLogger struct {
	collection *mongo.Collection
}

// NewMongoLogger creates a new MongoLogger.
func NewMongoLogger(uri, dbName, collectionName string) (ILogger, error) {
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	collection := client.Database(dbName).Collection(collectionName)

	return &MongoLogger{collection: collection}, nil
}

func (l *MongoLogger) Log(level LogLevel, message string, ctx map[string]interface{}) {
	entry := LogEntry{
		Timestamp: time.Now(),
		Level:     level,
		Message:   message,
		Context:   ctx,
	}

	_, err := l.collection.InsertOne(context.Background(), entry)
	if err != nil {
		log.Printf("Error logging to MongoDB: %v", err)
	}
}

func (l *MongoLogger) Trace(message string, ctx map[string]interface{}) { l.Log(Trace, message, ctx) }
func (l *MongoLogger) Debug(message string, ctx map[string]interface{}) { l.Log(Debug, message, ctx) }
func (l *MongoLogger) Info(message string, ctx map[string]interface{})  { l.Log(Info, message, ctx) }
func (l *MongoLogger) Warning(message string, ctx map[string]interface{}) {
	l.Log(Warning, message, ctx)
}
func (l *MongoLogger) Error(message string, ctx map[string]interface{}) { l.Log(Error, message, ctx) }
func (l *MongoLogger) Critical(message string, ctx map[string]interface{}) {
	l.Log(Critical, message, ctx)
}
