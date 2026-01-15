package startup

import (
	"context"

	"github.com/afteracademy/goserve/api/auth"
	authMW "github.com/afteracademy/goserve/api/auth/middleware"
	"github.com/afteracademy/goserve/api/blog"
	"github.com/afteracademy/goserve/api/blog/author"
	"github.com/afteracademy/goserve/api/blog/editor"
	"github.com/afteracademy/goserve/api/blogs"
	"github.com/afteracademy/goserve/api/contact"
	"github.com/afteracademy/goserve/api/user"
	coreMW "github.com/afteracademy/goserve/arch/middleware"
	"github.com/afteracademy/goserve/arch/mongo"
	"github.com/afteracademy/goserve/arch/network"
	"github.com/afteracademy/goserve/arch/redis"
	"github.com/afteracademy/goserve/config"
)

type Module network.Module[module]

type module struct {
	Context     context.Context
	Env         *config.Env
	DB          mongo.Database
	Store       redis.Store
	UserService user.Service
	AuthService auth.Service
	BlogService blog.Service
}

func (m *module) GetInstance() *module {
	return m
}

func (m *module) Controllers() []network.Controller {
	return []network.Controller{
		auth.NewController(m.AuthenticationProvider(), m.AuthorizationProvider(), m.AuthService),
		user.NewController(m.AuthenticationProvider(), m.AuthorizationProvider(), m.UserService),
		blog.NewController(m.AuthenticationProvider(), m.AuthorizationProvider(), m.BlogService),
		author.NewController(m.AuthenticationProvider(), m.AuthorizationProvider(), author.NewService(m.DB, m.BlogService)),
		editor.NewController(m.AuthenticationProvider(), m.AuthorizationProvider(), editor.NewService(m.DB, m.UserService)),
		blogs.NewController(m.AuthenticationProvider(), m.AuthorizationProvider(), blogs.NewService(m.DB, m.Store)),
		contact.NewController(m.AuthenticationProvider(), m.AuthorizationProvider(), contact.NewService(m.DB)),
	}
}

func (m *module) RootMiddlewares() []network.RootMiddleware {
	return []network.RootMiddleware{
		coreMW.NewErrorCatcher(), // NOTE: this should be the first handler to be mounted
		authMW.NewKeyProtection(m.AuthService),
		coreMW.NewNotFound(),
	}
}

func (m *module) AuthenticationProvider() network.AuthenticationProvider {
	return authMW.NewAuthenticationProvider(m.AuthService, m.UserService)
}

func (m *module) AuthorizationProvider() network.AuthorizationProvider {
	return authMW.NewAuthorizationProvider()
}

func NewModule(context context.Context, env *config.Env, db mongo.Database, store redis.Store) Module {
	userService := user.NewService(db)
	authService := auth.NewService(db, env, userService)
	blogService := blog.NewService(db, store, userService)

	return &module{
		Context:     context,
		Env:         env,
		DB:          db,
		Store:       store,
		UserService: userService,
		AuthService: authService,
		BlogService: blogService,
	}
}
