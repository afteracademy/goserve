package mongo

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Query[T any] interface {
	Close()
	CreateIndexes(indexes []mongo.IndexModel) error
	FindOne(filter bson.M, opts *options.FindOneOptions) (*T, error)
	FindAll(filter bson.M, opts *options.FindOptions) ([]*T, error)
	FindPaginated(filter bson.M, page int64, limit int64, opts *options.FindOptions) ([]*T, error)
	InsertOne(doc *T) (*primitive.ObjectID, error)
	InsertAndRetrieveOne(doc *T) (*T, error)
	InsertMany(doc []*T) ([]primitive.ObjectID, error)
	InsertAndRetrieveMany(doc []*T) ([]*T, error)
	UpdateOne(filter bson.M, update bson.M) (*mongo.UpdateResult, error)
	UpdateMany(filter bson.M, update bson.M) (*mongo.UpdateResult, error)
	DeleteOne(filter bson.M) (*mongo.DeleteResult, error)
}

type query[T any] struct {
	collection *mongo.Collection
	context    context.Context
	cancel     context.CancelFunc
}

func newSingleQuery[T any](collection *mongo.Collection, timeout time.Duration) Query[T] {
	context, cancel := context.WithTimeout(context.Background(), timeout)
	return &query[T]{
		context:    context,
		cancel:     cancel,
		collection: collection,
	}
}

func newQuery[T any](context context.Context, collection *mongo.Collection) Query[T] {
	return &query[T]{
		context:    context,
		collection: collection,
	}
}

func (q *query[T]) Close() {
	if q.cancel != nil {
		q.cancel()
	}
}

func (q *query[T]) CreateIndexes(indexes []mongo.IndexModel) error {
	defer q.Close()
	fmt.Println("database indexing for: " + q.collection.Name())
	result, err := q.collection.Indexes().CreateMany(q.context, indexes)
	fmt.Println(q.collection.Name(), result)
	return err
}

func (q *query[T]) FindOne(filter bson.M, opts *options.FindOneOptions) (*T, error) {
	defer q.Close()
	var doc T
	err := q.collection.FindOne(q.context, filter, opts).Decode(&doc)
	if err != nil {
		return nil, err
	}

	return &doc, nil
}

func (q *query[T]) FindAll(filter bson.M, opts *options.FindOptions) ([]*T, error) {
	defer q.Close()
	cursor, err := q.collection.Find(q.context, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %w", err)
	}
	defer cursor.Close(q.context)

	var docs []*T

	for cursor.Next(q.context) {
		var result T
		err := cursor.Decode(&result)
		if err != nil {
			return nil, fmt.Errorf("error decoding result: %w", err)
		}
		docs = append(docs, &result)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return docs, nil
}

func (q *query[T]) FindPaginated(filter bson.M, page int64, limit int64, opts *options.FindOptions) ([]*T, error) {
	defer q.Close()
	skip := (page - 1) * limit

	if opts == nil {
		opts = options.Find()
	}
	opts.SetSkip(skip)
	opts.SetLimit(int64(limit))

	cursor, err := q.collection.Find(q.context, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %w", err)
	}
	defer cursor.Close(q.context)

	var docs []*T

	for cursor.Next(q.context) {
		var result T
		err := cursor.Decode(&result)
		if err != nil {
			return nil, fmt.Errorf("error decoding result: %w", err)
		}
		docs = append(docs, &result)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return docs, nil
}

func (q *query[T]) InsertOne(doc *T) (*primitive.ObjectID, error) {
	defer q.Close()
	result, err := q.collection.InsertOne(q.context, doc)
	if err != nil {
		return nil, err
	}

	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, fmt.Errorf("database query error for: %s", insertedID)
	}

	return &insertedID, nil
}

func (q *query[T]) InsertAndRetrieveOne(doc *T) (*T, error) {
	defer q.Close()
	result, err := q.collection.InsertOne(q.context, doc)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": result.InsertedID}
	retrived, err := q.FindOne(filter, nil)
	if err != nil {
		return nil, err
	}

	return retrived, nil
}

func (q *query[T]) InsertMany(docs []*T) ([]primitive.ObjectID, error) {
	defer q.Close()
	var iDocs []any
	for _, doc := range docs {
		iDocs = append(iDocs, doc)
	}

	result, err := q.collection.InsertMany(q.context, iDocs)
	if err != nil {
		return nil, err
	}

	var insertedIDs []primitive.ObjectID

	for _, v := range result.InsertedIDs {
		insertedID, ok := v.(primitive.ObjectID)
		if !ok {
			return nil, fmt.Errorf("database query error for: %s", insertedID)
		}
		insertedIDs = append(insertedIDs, insertedID)
	}

	return insertedIDs, nil
}

func (q *query[T]) InsertAndRetrieveMany(docs []*T) ([]*T, error) {
	defer q.Close()
	var iDocs []any
	for _, doc := range docs {
		iDocs = append(iDocs, doc)
	}

	result, err := q.collection.InsertMany(q.context, iDocs)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": bson.M{"$in": result.InsertedIDs}}

	retrieved, err := q.FindAll(filter, nil)
	if err != nil {
		return nil, err
	}

	return retrieved, nil
}

/*
 * Example -> update := bson.M{"$set": bson.M{"field": "newValue"}}
 */
func (q *query[T]) UpdateOne(filter bson.M, update bson.M) (*mongo.UpdateResult, error) {
	defer q.Close()
	result, err := q.collection.UpdateOne(q.context, filter, update)
	if err != nil {
		return nil, err
	}
	return result, nil
}

/*
 * Example -> update := bson.M{"$set": bson.M{"field": "newValue"}}
 */
func (q *query[T]) UpdateMany(filter bson.M, update bson.M) (*mongo.UpdateResult, error) {
	defer q.Close()
	result, err := q.collection.UpdateMany(q.context, filter, update)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (q *query[T]) DeleteOne(filter bson.M) (*mongo.DeleteResult, error) {
	defer q.Close()
	result, err := q.collection.DeleteOne(q.context, filter)
	if err != nil {
		return nil, err
	}

	return result, nil
}
