package db

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Client structure represents client entity described in the task
type Client struct {
	SalesforceID   int32   `json:"salesforceId" bson:"salesforceId"`
	Country        string  `json:"country,omitempty"`
	Owner          string  `json:"owner,omitempty"`
	Manager        string  `json:"manager,omitempty"`
	PredictedUsage []int64 `json:"predictedUsage,omitempty"`
	ActualUsage    []int64 `json:"actualUsage,omitempty"`
}

// DB represents database client and options
type DB struct {
	Opts   *Opts
	client *mongo.Client
}

// Opts represents options passed via command line or env vars
type Opts struct {
	MongoURI   string `long:"conn_uri" env:"MONGO_URI" description:"mongo connection URI" default:"mongodb://localhost:27017"`
	DBName     string `long:"db_name" env:"DB_NAME" description:"mongo database name" default:"sg"`
	Mode       string `long:"mode" env:"MODE" description:"mode" default:"prod"`
	BuildTime  string `long:"build_time" env:"BUILD_TIME" description:"app build time" default:""`
	CommitHash string `long:"commit_hash" env:"COMMIT_HASH" description:"app commit hash" default:""`
}

// DefaultOpts represents default options (used to run application tests locally)
func DefaultOpts() *Opts {
	mongoUri := os.Getenv("MONGO_URI")
	dbName := os.Getenv("DB_NAME")

	if mongoUri == "" {
		mongoUri = "mongodb://localhost:27017"

	}

	if dbName == "" {
		dbName = "sg-test"

	}
	return &Opts{MongoURI: mongoUri, DBName: dbName}

}

// NewDB creates new database client
func NewDB(opts *Opts) (*DB, error) {

	// Set client options
	clientOptions := options.Client().ApplyURI(opts.MongoURI)

	// Connect to MongoDB
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	client, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		return nil, err
	}
	// Check the connection
	err = client.Ping(ctx, nil)

	if err != nil {
		log.Println("ERR. Cannot connect to Database. Check your connection details:", opts.MongoURI)
		return nil, err
	}
	if opts.Mode != "testing" {
		log.Println("Successfully connected to MongoDB!")
	}

	return &DB{client: client, Opts: opts}, nil
}

//Cleanup test data (used in testing)
func (s *DB) Cleanup() error {
	clientsCollection := s.client.Database(s.Opts.DBName).Collection("clients")
	_, err := clientsCollection.DeleteMany(context.Background(), bson.D{{}})
	if err != nil {
		return err
	}
	return nil
}

// NewClient creates new client
func (s *DB) NewClient(client *Client) (*Client, error) {
	collection := s.client.Database(s.Opts.DBName).Collection("clients")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := collection.InsertOne(ctx, client)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// NewClient returns all clients
func (s *DB) ListClients() ([]*Client, error) {
	collection := s.client.Database(s.Opts.DBName).Collection("clients")

	// Pass these options to the Find method
	findOptions := options.Find()
	findOptions.SetLimit(10000)

	// Here's an array in which you can store the decoded documents
	results := make([]*Client, 0)

	// Passing bson.D{{}} as the filter matches all documents in the collection
	cur, err := collection.Find(context.TODO(), bson.D{{}}, findOptions)
	if err != nil {
		return nil, err
	}

	// Finding multiple documents returns a cursor
	// Iterating through the cursor allows us to decode documents one at a time
	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var elem Client
		err := cur.Decode(&elem)
		if err != nil {
			return nil, err
		}

		results = append(results, &elem)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	// Close the cursor once finished
	err = cur.Close(context.TODO())
	if err != nil {
		log.Println("[WARN] Problem closing cursor")
	}
	return results, nil
}

// NewClient gets client by id
func (s *DB) GetClient(id int64) (*Client, error) {
	collection := s.client.Database(s.Opts.DBName).Collection("clients")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var client *Client
	query := bson.D{{"salesforceId", id}}
	err := collection.FindOne(ctx, query).Decode(&client)
	if err != nil {
		return nil, err
	}
	return client, nil
}
