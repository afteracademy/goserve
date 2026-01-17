package network

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/mock"
)

type MockDto struct {
	Field string `json:"field" form:"field" uri:"field" binding:"required" validate:"required"`
}

func (d *MockDto) GetValue() *MockDto {
	return d
}

func (b *MockDto) ValidateErrors(errs validator.ValidationErrors) ([]string, error) {
	var msgs []string
	for _, err := range errs {
		switch err.Tag() {
		case "required":
			msgs = append(msgs, fmt.Sprintf("%s is required", err.Field()))
		default:
			msgs = append(msgs, fmt.Sprintf("%s is invalid", err.Field()))
		}
	}
	return msgs, nil
}

func MockSuccessMsgHandler(msg string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		NewResponseSender().Send(ctx).SuccessMsgResponse(msg)
	}
}

func MockSuccessDataHandler(msg string, data any) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		NewResponseSender().Send(ctx).SuccessDataResponse(msg, data)
	}
}

func MockTestHandler(
	t *testing.T, httpMethod, path, url, body string,
	handler gin.HandlerFunc,
	headers map[string]string,
) *httptest.ResponseRecorder {
	gin.SetMode(gin.TestMode)

	rr := httptest.NewRecorder()
	ctx, r := gin.CreateTestContext(rr)
	r.Handle(httpMethod, path, handler)

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(CustomTagNameFunc())
	}

	req, err := http.NewRequest(httpMethod, url, bytes.NewBuffer([]byte(body)))
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	ctx.Request = req

	r.ServeHTTP(rr, req)

	return rr
}

func MockTestRootMiddleware(
	t *testing.T,
	middleware RootMiddleware,
	handler gin.HandlerFunc,
	headers map[string]string,
) *httptest.ResponseRecorder {
	gin.SetMode(gin.TestMode)

	rr := httptest.NewRecorder()
	ctx, r := gin.CreateTestContext(rr)
	middleware.Attach(r)
	r.GET("/", handler)

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(CustomTagNameFunc())
	}

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	ctx.Request = req

	r.ServeHTTP(rr, req)

	return rr
}

func MockTestRootMiddlewareWithUrl(
	t *testing.T, path, url string,
	middleware RootMiddleware,
	handler gin.HandlerFunc,
	headers map[string]string,
) *httptest.ResponseRecorder {
	gin.SetMode(gin.TestMode)

	rr := httptest.NewRecorder()
	ctx, r := gin.CreateTestContext(rr)
	middleware.Attach(r)
	r.GET(path, handler)

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(CustomTagNameFunc())
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	ctx.Request = req

	r.ServeHTTP(rr, req)

	return rr
}

func MockTestAuthenticationProvider(
	t *testing.T,
	auth AuthenticationProvider,
	handler gin.HandlerFunc,
	headers map[string]string,
) *httptest.ResponseRecorder {
	gin.SetMode(gin.TestMode)

	rr := httptest.NewRecorder()
	ctx, r := gin.CreateTestContext(rr)
	r.Use(auth.Middleware())
	r.GET("/", handler)

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(CustomTagNameFunc())
	}

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	ctx.Request = req

	r.ServeHTTP(rr, req)

	return rr
}

func MockTestAuthorizationProvider(
	t *testing.T,
	role string,
	auth AuthenticationProvider,
	authz AuthorizationProvider,
	handler gin.HandlerFunc,
	headers map[string]string,
) *httptest.ResponseRecorder {
	gin.SetMode(gin.TestMode)

	rr := httptest.NewRecorder()
	ctx, r := gin.CreateTestContext(rr)
	r.Use(auth.Middleware())
	if len(role) == 0 {
		r.Use(authz.Middleware())
	} else {
		r.Use(authz.Middleware(role))
	}
	r.GET("/", handler)

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(CustomTagNameFunc())
	}

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	ctx.Request = req

	r.ServeHTTP(rr, req)

	return rr
}

func MockTestController(
	t *testing.T, httpMethod, url, body string,
	controller Controller,
) *httptest.ResponseRecorder {
	gin.SetMode(gin.TestMode)

	rr := httptest.NewRecorder()
	ctx, r := gin.CreateTestContext(rr)

	controller.MountRoutes(r.Group(controller.Path()))

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(CustomTagNameFunc())
	}

	req, err := http.NewRequest(httpMethod, url, bytes.NewBuffer([]byte(body)))
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	ctx.Request = req

	r.ServeHTTP(rr, req)

	return rr
}

type MockAuthenticationProvider struct {
	mock.Mock
}

func (m *MockAuthenticationProvider) Debug() bool {
	return true
}

func (m *MockAuthenticationProvider) Middleware() gin.HandlerFunc {
	args := m.Called()
	return args.Get(0).(gin.HandlerFunc)
}

func (m *MockAuthenticationProvider) Send(ctx *gin.Context) SendResponse {
	args := m.Called(ctx)
	return args.Get(0).(SendResponse)
}

type MockAuthorizationProvider struct {
	mock.Mock
}

func (m *MockAuthorizationProvider) Debug() bool {
	return true
}

func (m *MockAuthorizationProvider) Middleware(params ...string) gin.HandlerFunc {
	args := m.Called(params)
	return args.Get(0).(gin.HandlerFunc)
}

func (m *MockAuthorizationProvider) Send(ctx *gin.Context) SendResponse {
	args := m.Called(ctx)
	return args.Get(0).(SendResponse)
}
