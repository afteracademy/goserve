package model

import (
	"context"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/afteracademy/goserve/arch/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongod "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const UserCollectionName = "users"

type User struct {
	ID            primitive.ObjectID   `bson:"_id,omitempty"`
	Name          string               `bson:"name" validate:"required,max=200"`
	Email         string               `bson:"email" validate:"required,email"`
	Password      *string              `bson:"password" validate:"required,min=6,max=100"`
	ProfilePicURL *string              `bson:"profilePicUrl,omitempty" validate:"omitempty,max=500"`
	Roles         []primitive.ObjectID `bson:"roles,omitempty" validate:"required"`
	Verified      bool                 `bson:"verified" validate:"-"`
	Status        bool                 `bson:"status" validate:"-"`
	CreatedAt     time.Time            `bson:"createdAt" validate:"required"`
	UpdatedAt     time.Time            `bson:"updatedAt" validate:"required"`

	// docs
	RoleDocs []*Role `bson:"-" validate:"-"`
}

func NewUser(email string, pwdHash string, name string, profilePicUrl *string, roles []*Role) (*User, error) {
	roleIds := make([]primitive.ObjectID, len(roles))
	for i, role := range roles {
		roleIds[i] = role.ID
	}

	now := time.Now()
	u := User{
		Email:         email,
		Password:      &pwdHash,
		Name:          name,
		ProfilePicURL: profilePicUrl,
		Roles:         roleIds,
		Verified:      false,
		Status:        true,
		CreatedAt:     now,
		UpdatedAt:     now,
		RoleDocs:      roles,
	}
	if err := u.Validate(); err != nil {
		return nil, err
	}
	return &u, nil
}

func (user *User) GetValue() *User {
	return user
}

func (user *User) Validate() error {
	validate := validator.New()
	return validate.Struct(user)
}

func (*User) EnsureIndexes(db mongo.Database) {
	indexes := []mongod.IndexModel{
		{
			Keys: bson.D{
				{Key: "_id", Value: 1},
				{Key: "status", Value: 1},
			},
		},
		{
			Keys: bson.D{
				{Key: "email", Value: 1},
				{Key: "status", Value: 1},
			},
		},
		{
			Keys: bson.D{
				{Key: "email", Value: 1},
			},
			Options: options.Index().SetUnique(true),
		},
	}
	mongo.NewQueryBuilder[User](db, UserCollectionName).Query(context.Background()).CreateIndexes(indexes)
}
