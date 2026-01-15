package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Helper function to capitalize the first letter of a string
func capitalizeFirstLetter(str string) string {
	if len(str) == 0 {
		return str
	}
	return strings.ToUpper(string(str[0])) + str[1:]
}

func generateFeature(featureTemplate string) error {
	if featureTemplate == "" {
		return errors.New("api name should be a non-empty string")
	}

	featureName := strings.ToLower(featureTemplate)
	featureDir := filepath.Join("api", featureName)
	if _, err := os.Stat(featureDir); err == nil {
		fmt.Println(featureName, "already exists")
		return nil
	}

	// Create api directory
	if err := os.MkdirAll(featureDir, os.ModePerm); err != nil {
		return err
	}

	if err := generateDto(featureDir, featureName); err != nil {
		return err
	}
	if err := generateModel(featureDir, featureName); err != nil {
		return err
	}
	if err := generateService(featureDir, featureName); err != nil {
		return err
	}
	if err := generateController(featureDir, featureName); err != nil {
		return err
	}
	return nil
}

func generateService(featureDir, featureName string) error {
	featureLower := strings.ToLower(featureName)
	featureCaps := capitalizeFirstLetter(featureName)
	servicePath := filepath.Join(featureDir, fmt.Sprintf("%sservice.go", ""))

	template := fmt.Sprintf(`package %s

import (
  "github.com/afteracademy/goserve/api/%s/dto"
	"github.com/afteracademy/goserve/api/%s/model"
	"github.com/afteracademy/goserve/arch/mongo"
	"github.com/afteracademy/goserve/arch/network"
	"github.com/afteracademy/goserve/arch/redis"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service interface {
	Find%s(id primitive.ObjectID) (*model.%s, error)
}

type service struct {
	network.BaseService
	%sQueryBuilder mongo.QueryBuilder[model.%s]
	info%sCache    redis.Cache[dto.Info%s]
}

func NewService(db mongo.Database, store redis.Store) Service {
	return &service{
		BaseService:  network.NewBaseService(),
		%sQueryBuilder: mongo.NewQueryBuilder[model.%s](db, model.CollectionName),
		info%sCache: redis.NewCache[dto.Info%s](store),
	}
}

func (s *service) Find%s(id primitive.ObjectID) (*model.%s, error) {
	filter := bson.M{"_id": id}

	msg, err := s.%sQueryBuilder.SingleQuery().FindOne(filter, nil)
	if err != nil {
		return nil, err
	}

	return msg, nil
}
`, featureLower, featureLower, featureLower, featureCaps, featureCaps, featureLower, featureCaps, featureCaps, featureCaps, featureLower, featureCaps, featureCaps, featureCaps, featureCaps, featureCaps, featureLower)

	return os.WriteFile(servicePath, []byte(template), os.ModePerm)
}

func generateController(featureDir, featureName string) error {
	featureLower := strings.ToLower(featureName)
	featureCaps := capitalizeFirstLetter(featureName)
	controllerPath := filepath.Join(featureDir, fmt.Sprintf("%scontroller.go", ""))

	template := fmt.Sprintf(`package %s

import (
	"github.com/gin-gonic/gin"
	"github.com/afteracademy/goserve/api/%s/dto"
	"github.com/afteracademy/goserve/common"
	coredto "github.com/afteracademy/goserve/arch/dto"
	"github.com/afteracademy/goserve/arch/network"
	"github.com/afteracademy/goserve/utils"
)

type controller struct {
	network.BaseController
	common.ContextPayload
	service Service
}

func NewController(
	authMFunc network.AuthenticationProvider,
	authorizeMFunc network.AuthorizationProvider,
	service Service,
) network.Controller {
	return &controller{
		BaseController: network.NewBaseController("/%s", authMFunc, authorizeMFunc),
		ContextPayload: common.NewContextPayload(),
		service:  service,
	}
}

func (c *controller) MountRoutes(group *gin.RouterGroup) {
	group.GET("/id/:id", c.get%sHandler)
}

func (c *controller) get%sHandler(ctx *gin.Context) {
	mongoId, err := network.ReqParams(ctx, coredto.EmptyMongoId())
	if err != nil {
		c.Send(ctx).BadRequestError(err.Error(), err)
		return
	}

	%s, err := c.service.Find%s(mongoId.ID)
	if err != nil {
		c.Send(ctx).NotFoundError("%s not found", err)
		return
	}

	data, err := utils.MapTo[dto.Info%s](%s)
	if err != nil {
		c.Send(ctx).InternalServerError("something went wrong", err)
		return
	}

	c.Send(ctx).SuccessDataResponse("success", data)
}
`, featureLower, featureLower, featureLower, featureCaps, featureCaps, featureLower, featureCaps, featureLower, featureCaps, featureLower)

	return os.WriteFile(controllerPath, []byte(template), os.ModePerm)
}

func generateModel(featureDir, featureName string) error {
	modelDirPath := filepath.Join(featureDir, "model")
	if err := os.MkdirAll(modelDirPath, os.ModePerm); err != nil {
		return err
	}

	featureLower := strings.ToLower(featureName)
	featureCaps := capitalizeFirstLetter(featureName)
	modelPath := filepath.Join(featureDir, fmt.Sprintf("model/%s.go", featureLower))

	tStr := `package model

import (
	"context"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/afteracademy/goserve/arch/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongod "go.mongodb.org/mongo-driver/mongo"
)

const CollectionName = "%ss"

type %s struct {
	ID        primitive.ObjectID ` + "`" + `bson:"_id,omitempty" validate:"-"` + "`" + `
	Field     string             ` + "`" + `bson:"field" validate:"required"` + "`" + `
	Status    bool               ` + "`" + `bson:"status" validate:"required"` + "`" + `
	CreatedAt time.Time          ` + "`" + `bson:"createdAt" validate:"required"` + "`" + `
	UpdatedAt time.Time          ` + "`" + `bson:"updatedAt" validate:"required"` + "`" + `
}` + `

func New%s(field string) (*%s, error) {
	time := time.Now()
	doc := %s{
		Field:     field,
		Status:    true,
		CreatedAt: time,
		UpdatedAt: time,
	}
	if err := doc.Validate(); err != nil {
		return nil, err
	}
	return &doc, nil
}

func (doc *%s) GetValue() *%s {
	return doc
}

func (doc *%s) Validate() error {
	validate := validator.New()
	return validate.Struct(doc)
}

func (*%s) EnsureIndexes(db mongo.Database) {
	indexes := []mongod.IndexModel{
		{
			Keys: bson.D{
				{Key: "_id", Value: 1},
				{Key: "status", Value: 1},
			},
		},
	}
	
	mongo.NewQueryBuilder[%s](db, CollectionName).Query(context.Background()).CreateIndexes(indexes)
}

`
	template := fmt.Sprintf(tStr, featureLower, featureCaps, featureCaps, featureCaps, featureCaps, featureCaps, featureCaps, featureCaps, featureCaps, featureCaps)

	return os.WriteFile(modelPath, []byte(template), os.ModePerm)
}

func generateDto(featureDir, featureName string) error {
	dtoDirPath := filepath.Join(featureDir, "dto")
	if err := os.MkdirAll(dtoDirPath, os.ModePerm); err != nil {
		return err
	}

	featureLower := strings.ToLower(featureName)
	featureCaps := capitalizeFirstLetter(featureName)
	dtoPath := filepath.Join(featureDir, fmt.Sprintf("dto/create_%s.go", featureLower))

	tStr := `package dto

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Info%s struct {
	ID        primitive.ObjectID ` + "`" + `json:"_id" binding:"required"` + "`" + `
	Field     string             ` + "`" + `json:"field" binding:"required"` + "`" + `
	CreatedAt time.Time          ` + "`" + `json:"createdAt" binding:"required"` + "`" + `
}

func EmptyInfo%s() *Info%s {
	return &Info%s{}
}

func (d *Info%s) GetValue() *Info%s {
	return d
}

func (d *Info%s) ValidateErrors(errs validator.ValidationErrors) ([]string, error) {
	var msgs []string
	for _, err := range errs {
		switch err.Tag() {
		case "required":
			msgs = append(msgs, fmt.Sprintf("%%s is required", err.Field()))
		case "min":
			msgs = append(msgs, fmt.Sprintf("%%s must be min %%s", err.Field(), err.Param()))
		case "max":
			msgs = append(msgs, fmt.Sprintf("%%s must be max %%s", err.Field(), err.Param()))
		default:
			msgs = append(msgs, fmt.Sprintf("%%s is invalid", err.Field()))
		}
	}
	return msgs, nil
}
`
	template := fmt.Sprintf(tStr, featureCaps, featureCaps, featureCaps, featureCaps, featureCaps, featureCaps, featureCaps)

	return os.WriteFile(dtoPath, []byte(template), os.ModePerm)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("api name should be non-empty string")
		return
	}

	featureName := os.Args[1]
	if err := generateFeature(featureName); err != nil {
		fmt.Println("Error:", err)
	}
}
