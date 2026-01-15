package model

import (
	"context"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/afteracademy/goserve/arch/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongod "go.mongodb.org/mongo-driver/mongo"
)

const KeystoreCollectionName = "keystores"

type Keystore struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Client       primitive.ObjectID `bson:"client" validate:"required"`
	PrimaryKey   string             `bson:"pKey" validate:"required"`
	SecondaryKey string             `bson:"sKey" validate:"required"`
	Status       bool               `bson:"status" validate:"-"`
	CreatedAt    time.Time          `bson:"createdAt" validate:"required"`
	UpdatedAt    time.Time          `bson:"updatedAt" validate:"required"`
}

func NewKeystore(clientID primitive.ObjectID, primaryKey string, secondaryKey string) (*Keystore, error) {
	now := time.Now()
	k := Keystore{
		Client:       clientID,
		PrimaryKey:   primaryKey,
		SecondaryKey: secondaryKey,
		Status:       true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if err := k.Validate(); err != nil {
		return nil, err
	}
	return &k, nil
}

func (keystore *Keystore) GetValue() *Keystore {
	return keystore
}

func (keystore *Keystore) Validate() error {
	validate := validator.New()
	return validate.Struct(keystore)
}

func (*Keystore) EnsureIndexes(db mongo.Database) {
	indexes := []mongod.IndexModel{
		{
			Keys: bson.D{
				{Key: "client", Value: 1},
				{Key: "pKey", Value: 1},
				{Key: "status", Value: 1},
			},
		},
		{
			Keys: bson.D{
				{Key: "client", Value: 1},
				{Key: "pKey", Value: 1},
				{Key: "sKey", Value: 1},
				{Key: "status", Value: 1},
			},
		},
	}
	mongo.NewQueryBuilder[Keystore](db, KeystoreCollectionName).Query(context.Background()).CreateIndexes(indexes)
}
