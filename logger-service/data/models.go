package data

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

const (
	DBName = "logs"
	CollectionName = "logs"
	SecondTimeout = 15
)


type Models struct{
	LogEntry LogEntry
}

type LogEntry struct{
	ID string `bson:"_id,omitempty" json:"id,omitempty"`
	Name string `bson:"name" json:"name"`
	Data string `bson:"data" json:"data"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdateAt time.Time `bson:"updated_at" json:"updated_at"`
}

var client *mongo.Client

func New(mongo *mongo.Client) Models{
	client = mongo

	return Models{
		LogEntry: LogEntry{},
	}
}

func (l* LogEntry) Insert(entry LogEntry) error{
	collection := client.Database(DBName).Collection(CollectionName)

	_, err := collection.InsertOne(context.TODO(), entry);
	if err != nil{
		log.Println("Error Inserting into Logs", err)
		return err
	}
	return nil
}

func (l* LogEntry) Insert1() error{
	collection := client.Database(DBName).Collection(CollectionName)

	_, err := collection.InsertOne(context.TODO(), l);
	if err != nil{
		log.Println("Error Inserting into Logs", err)
		return err
	}
	return nil
}

func (l* LogEntry) All()([]LogEntry, error){
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database(DBName).Collection(CollectionName)
	findAllOptions := options.Find()
	findAllOptions.SetSort(bson.D{{"created_at", -1}})

	cursor, err := collection.Find(context.TODO(), bson.D{}, findAllOptions)
	if err!=nil{
		log.Println("Error While Finding All Logs", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var logSlice []LogEntry
	for cursor.Next(ctx){
		var item LogEntry
		err = cursor.Decode(item)
		if err != nil{
			log.Println("Error While Decoding")
			return nil, err;
		}
		logSlice = append(logSlice, item)
	}
	return logSlice, nil
}

func (l* LogEntry) GetOne(id string)(LogEntry, error){
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database(DBName).Collection(CollectionName)
	docId, err := primitive.ObjectIDFromHex(id)
	if err != nil{
		return LogEntry{}, err
	}

	var entry LogEntry
	err = collection.FindOne(ctx, bson.M{"_id":docId}).Decode(&entry)
	if err != nil{
		return LogEntry{}, err
	}

	return entry, nil
}

func (l* LogEntry) DropCollection() error{
	ctx, cancel := context.WithTimeout(context.Background(), SecondTimeout*time.Second)
	defer cancel()

	collection := client.Database(DBName).Collection(CollectionName)
	if err := collection.Drop(ctx); err != nil{
		return err
	}
	return nil
}

func (l* LogEntry) Update()(*mongo.UpdateResult, error){
	ctx, cancel := context.WithTimeout(context.Background(), SecondTimeout*time.Second)
	defer cancel()

	collection := client.Database(DBName).Collection(CollectionName)

	docId, err := primitive.ObjectIDFromHex(l.ID)
	if err != nil{
		return nil, err
	}

	result, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": docId},
		bson.D{
			{"$set", bson.D{
				{"name", l.Name},
				{"data", l.Data},
				{"updated_at", time.Now()},
			}},
		},
	)

	if err!= nil{
		return nil, err
	}
	return result, nil
}