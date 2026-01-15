package startup

import (
	auth "github.com/afteracademy/goserve/api/auth/model"
	blog "github.com/afteracademy/goserve/api/blog/model"
	contact "github.com/afteracademy/goserve/api/contact/model"
	user "github.com/afteracademy/goserve/api/user/model"
	"github.com/afteracademy/goserve/arch/mongo"
)

func EnsureDbIndexes(db mongo.Database) {
	go mongo.Document[auth.Keystore](&auth.Keystore{}).EnsureIndexes(db)
	go mongo.Document[auth.ApiKey](&auth.ApiKey{}).EnsureIndexes(db)
	go mongo.Document[user.User](&user.User{}).EnsureIndexes(db)
	go mongo.Document[user.Role](&user.Role{}).EnsureIndexes(db)
	go mongo.Document[blog.Blog](&blog.Blog{}).EnsureIndexes(db)
	go mongo.Document[contact.Message](&contact.Message{}).EnsureIndexes(db)
}
