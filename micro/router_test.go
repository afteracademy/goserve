package micro

import (
	"testing"
	"time"

	"github.com/afteracademy/goserve/v2/network"
	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/micro"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockController is a mock for Controller interface
type MockController struct {
	mock.Mock
}

func (m *MockController) Path() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockController) Authentication() gin.HandlerFunc {
	args := m.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(gin.HandlerFunc)
}

func (m *MockController) Authorization(role string) gin.HandlerFunc {
	args := m.Called(role)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(gin.HandlerFunc)
}

func (m *MockController) Debug() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockController) SendNats(req NatsRequest) SendMessage {
	args := m.Called(req)
	return args.Get(0).(SendMessage)
}

func (m *MockController) MountNats(group NatsGroup) {
	m.Called(group)
}

func (m *MockController) MountRoutes(group *gin.RouterGroup) {
	m.Called(group)
}

func TestNewRouter(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("should create a new router with test mode", func(t *testing.T) {
		s := RunNatsServerOnPort(t, -1)
		defer s.Shutdown()

		nc, err := nats.Connect(s.ClientURL())
		assert.NoError(t, err)
		defer nc.Close()

		service, err := micro.AddService(nc, micro.Config{
			Name:    "test-service",
			Version: "1.0.0",
		})
		assert.NoError(t, err)

		mockNatsClient := &natsClient{
			Conn:    nc,
			Service: service,
			Timeout: 2 * time.Second,
		}

		router := NewRouter(gin.TestMode, mockNatsClient)
		assert.NotNil(t, router)
		assert.NotNil(t, router.GetEngine())
		assert.Equal(t, mockNatsClient, router.NatsClient())
	})

	t.Run("should create router with release mode", func(t *testing.T) {
		s := RunNatsServerOnPort(t, -1)
		defer s.Shutdown()

		nc, err := nats.Connect(s.ClientURL())
		assert.NoError(t, err)
		defer nc.Close()

		service, err := micro.AddService(nc, micro.Config{
			Name:    "test-service",
			Version: "1.0.0",
		})
		assert.NoError(t, err)

		mockNatsClient := &natsClient{
			Conn:    nc,
			Service: service,
			Timeout: 2 * time.Second,
		}

		router := NewRouter(gin.ReleaseMode, mockNatsClient)
		assert.NotNil(t, router)
		assert.NotNil(t, router.GetEngine())
	})
}

func TestRouter_GetEngine(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("should return gin engine", func(t *testing.T) {
		s := RunNatsServerOnPort(t, -1)
		defer s.Shutdown()

		nc, err := nats.Connect(s.ClientURL())
		assert.NoError(t, err)
		defer nc.Close()

		service, err := micro.AddService(nc, micro.Config{
			Name:    "test-service",
			Version: "1.0.0",
		})
		assert.NoError(t, err)

		mockNatsClient := &natsClient{
			Conn:    nc,
			Service: service,
			Timeout: 2 * time.Second,
		}

		router := NewRouter(gin.TestMode, mockNatsClient)
		engine := router.GetEngine()
		assert.NotNil(t, engine)
		assert.IsType(t, &gin.Engine{}, engine)
	})
}

func TestRouter_NatsClient(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("should return nats client", func(t *testing.T) {
		s := RunNatsServerOnPort(t, -1)
		defer s.Shutdown()

		nc, err := nats.Connect(s.ClientURL())
		assert.NoError(t, err)
		defer nc.Close()

		service, err := micro.AddService(nc, micro.Config{
			Name:    "test-service",
			Version: "1.0.0",
		})
		assert.NoError(t, err)

		mockNatsClient := &natsClient{
			Conn:    nc,
			Service: service,
			Timeout: 2 * time.Second,
		}

		router := NewRouter(gin.TestMode, mockNatsClient)
		client := router.NatsClient()
		assert.NotNil(t, client)
		assert.Equal(t, mockNatsClient, client)
	})
}

func TestRouter_LoadRootMiddlewares(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("should load root middlewares", func(t *testing.T) {
		s := RunNatsServerOnPort(t, -1)
		defer s.Shutdown()

		nc, err := nats.Connect(s.ClientURL())
		assert.NoError(t, err)
		defer nc.Close()

		service, err := micro.AddService(nc, micro.Config{
			Name:    "test-service",
			Version: "1.0.0",
		})
		assert.NoError(t, err)

		mockNatsClient := &natsClient{
			Conn:    nc,
			Service: service,
			Timeout: 2 * time.Second,
		}

		router := NewRouter(gin.TestMode, mockNatsClient)

		middlewares := []network.RootMiddleware{}

		assert.NotPanics(t, func() {
			router.LoadRootMiddlewares(middlewares)
		})
	})
}

func TestRouter_LoadControllers(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("should load HTTP and NATS controllers", func(t *testing.T) {
		s := RunNatsServerOnPort(t, -1)
		defer s.Shutdown()

		nc, err := nats.Connect(s.ClientURL())
		assert.NoError(t, err)
		defer nc.Close()

		service, err := micro.AddService(nc, micro.Config{
			Name:    "test-service",
			Version: "1.0.0",
		})
		assert.NoError(t, err)

		mockNatsClient := &natsClient{
			Conn:    nc,
			Service: service,
			Timeout: 2 * time.Second,
		}

		router := NewRouter(gin.TestMode, mockNatsClient)

		mockController := new(MockController)
		mockController.On("Path").Return("/api/test")
		mockController.On("MountRoutes", mock.AnythingOfType("*gin.RouterGroup")).Once()
		mockController.On("MountNats", mock.Anything).Once()

		controllers := []Controller{mockController}
		router.LoadControllers(controllers)

		mockController.AssertExpectations(t)
	})

	t.Run("should load multiple controllers", func(t *testing.T) {
		s := RunNatsServerOnPort(t, -1)
		defer s.Shutdown()

		nc, err := nats.Connect(s.ClientURL())
		assert.NoError(t, err)
		defer nc.Close()

		service, err := micro.AddService(nc, micro.Config{
			Name:    "test-service",
			Version: "1.0.0",
		})
		assert.NoError(t, err)

		mockNatsClient := &natsClient{
			Conn:    nc,
			Service: service,
			Timeout: 2 * time.Second,
		}

		router := NewRouter(gin.TestMode, mockNatsClient)

		mockController1 := new(MockController)
		mockController1.On("Path").Return("/api/users")
		mockController1.On("MountRoutes", mock.AnythingOfType("*gin.RouterGroup")).Once()
		mockController1.On("MountNats", mock.Anything).Once()

		mockController2 := new(MockController)
		mockController2.On("Path").Return("/api/posts")
		mockController2.On("MountRoutes", mock.AnythingOfType("*gin.RouterGroup")).Once()
		mockController2.On("MountNats", mock.Anything).Once()

		controllers := []Controller{mockController1, mockController2}
		router.LoadControllers(controllers)

		mockController1.AssertExpectations(t)
		mockController2.AssertExpectations(t)
	})

	t.Run("should handle root path controller", func(t *testing.T) {
		s := RunNatsServerOnPort(t, -1)
		defer s.Shutdown()

		nc, err := nats.Connect(s.ClientURL())
		assert.NoError(t, err)
		defer nc.Close()

		service, err := micro.AddService(nc, micro.Config{
			Name:    "root-service",
			Version: "1.0.0",
		})
		assert.NoError(t, err)

		mockNatsClient := &natsClient{
			Conn:    nc,
			Service: service,
			Timeout: 2 * time.Second,
		}

		router := NewRouter(gin.TestMode, mockNatsClient)

		mockController := new(MockController)
		mockController.On("Path").Return("/")
		mockController.On("MountRoutes", mock.AnythingOfType("*gin.RouterGroup")).Once()
		mockController.On("MountNats", mock.Anything).Once()

		controllers := []Controller{mockController}
		router.LoadControllers(controllers)

		mockController.AssertExpectations(t)
	})
}
