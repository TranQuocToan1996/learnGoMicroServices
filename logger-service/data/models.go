package data

import (
	"context"
	"fmt"
	"log"
	"time"

	"log-service/utilities"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client     *mongo.Client
	loggerDb   *mongo.Database
	loggerColl *mongo.Collection
)

func New(mongo *mongo.Client) Models {
	client = mongo
	loggerDb = client.Database("logger")
	loggerColl = loggerDb.Collection("logger")
	return Models{
		LogEntry{},
	}
}

func NewDatabase(name string, connection string) *mongo.Database {
	opt := options.Client().ApplyURI(connection)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, opt)
	if err != nil {
		log.Fatal(err.Error())
	}

	return client.Database(name)
}

type Models struct {
	LogEntry
}

type LogEntry struct {
	ID       string    `bson:"_id,omitempty" json:"id,omitempty"`
	Name     string    `bson:"name,omitempty" json:"name,omitempty"`
	Data     string    `bson:"data,omitempty" json:"data,omitempty"`
	CreateAt time.Time `bson:"creatAt" json:"creatAt"`
	UpdateAt time.Time `bson:"updateAt" json:"updateAt"`
}

func (l *LogEntry) Insert(entry LogEntry) error {
	collection := loggerColl

	now := utilities.TimeLocalNow()

	_, err := collection.InsertOne(context.Background(), LogEntry{
		Name:     entry.Name,
		Data:     entry.Data,
		CreateAt: now,
		UpdateAt: now,
	})
	if err != nil {
		return err
	}
	return nil
}

func (l *LogEntry) All(entry LogEntry) (logs []*LogEntry, missing []string, err error) {
	collection := loggerColl

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	otps := options.Find()
	otps.SetSort(bson.D{
		{Key: "createAt", Value: -1},
	})

	// No filter find
	cur, err := collection.Find(context.TODO(), bson.D{}, otps)
	if err != nil {
		fmt.Println("Find document err")
		return nil, nil, err
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		l := &LogEntry{}
		err := cur.Decode(l)
		if err != nil {
			log.Println("error decode", cur.Current.String())
			missing = append(missing, cur.Current.String())
			continue
		}
		logs = append(logs, l)
	}

	return logs, missing, nil

}

func (l *LogEntry) CreateIndex(collectionName string, field string, unique bool) bool {

	// 1. Lets define the keys for the index we want to create
	mod := mongo.IndexModel{
		Keys:    bson.M{field: 1}, // index in ascending order or -1 for descending order
		Options: options.Index().SetUnique(unique),
	}

	// 2. Create the context for this operation
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 3. Connect to the database and access the collection
	collection := loggerColl

	// 4. Create a single index
	_, err := collection.Indexes().CreateOne(ctx, mod)
	if err != nil {
		// 5. Something went wrong, we log it and return false
		fmt.Println(err.Error())
		return false
	}

	// 6. All went well, we return true
	return true
}
