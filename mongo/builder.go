package mongo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type QueryBuilder[T any] interface {
	GetCollection() *mongo.Collection
	SingleQuery() Query[T]
	Query(context context.Context) Query[T]
}

type queryBuilder[T any] struct {
	collection *mongo.Collection
	timeout    time.Duration
}

func (c *queryBuilder[T]) GetCollection() *mongo.Collection {
	return c.collection
}

func (c *queryBuilder[T]) SingleQuery() Query[T] {
	return newSingleQuery[T](c.collection, c.timeout)
}

func (c *queryBuilder[T]) Query(context context.Context) Query[T] {
	return newQuery[T](context, c.collection)
}

func NewQueryBuilder[T any](db Database, collectionName string) QueryBuilder[T] {
	return &queryBuilder[T]{
		collection: db.GetInstance().Collection(collectionName),
		timeout:    db.GetInstance().config.Timeout,
	}
}
