package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/mock"
	"github.com/afteracademy/goserve/api/auth/dto"
	"github.com/afteracademy/goserve/api/auth/model"
	userModel "github.com/afteracademy/goserve/api/user/model"
)

type MockService struct {
	mock.Mock
}

func (m *MockService) SignUpBasic(signUpDto *dto.SignUpBasic) (*dto.UserAuth, error) {
	args := m.Called(signUpDto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.UserAuth), args.Error(1)
}

func (m *MockService) SignInBasic(signInDto *dto.SignInBasic) (*dto.UserAuth, error) {
	args := m.Called(signInDto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.UserAuth), args.Error(1)
}

func (m *MockService) RenewToken(tokenRefreshDto *dto.TokenRefresh, accessToken string) (*dto.UserTokens, error) {
	args := m.Called(tokenRefreshDto, accessToken)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.UserTokens), args.Error(1)
}

func (m *MockService) SignOut(keystore *model.Keystore) error {
	args := m.Called(keystore)
	return args.Error(0)
}

func (m *MockService) IsEmailRegisted(email string) bool {
	args := m.Called(email)
	return args.Bool(0)
}

func (m *MockService) GenerateToken(user *userModel.User) (string, string, error) {
	args := m.Called(user)
	return args.String(0), args.String(1), args.Error(2)
}

func (m *MockService) CreateKeystore(client *userModel.User, primaryKey string, secondaryKey string) (*model.Keystore, error) {
	args := m.Called(client, primaryKey, secondaryKey)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Keystore), args.Error(1)
}

func (m *MockService) FindKeystore(client *userModel.User, primaryKey string) (*model.Keystore, error) {
	args := m.Called(client, primaryKey)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Keystore), args.Error(1)
}

func (m *MockService) FindRefreshKeystore(client *userModel.User, pKey string, sKey string) (*model.Keystore, error) {
	args := m.Called(client, pKey, sKey)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Keystore), args.Error(1)
}

func (m *MockService) VerifyToken(tokenStr string) (*jwt.RegisteredClaims, error) {
	args := m.Called(tokenStr)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*jwt.RegisteredClaims), args.Error(1)
}

func (m *MockService) DecodeToken(tokenStr string) (*jwt.RegisteredClaims, error) {
	args := m.Called(tokenStr)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*jwt.RegisteredClaims), args.Error(1)
}

func (m *MockService) SignToken(claims jwt.RegisteredClaims) (string, error) {
	args := m.Called(claims)
	return args.String(0), args.Error(1)
}

func (m *MockService) ValidateClaims(claims *jwt.RegisteredClaims) bool {
	args := m.Called(claims)
	return args.Bool(0)
}

func (m *MockService) FindApiKey(key string) (*model.ApiKey, error) {
	args := m.Called(key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.ApiKey), args.Error(1)
}

func (m *MockService) CreateApiKey(key string, version int, permissions []model.Permission, comments []string) (*model.ApiKey, error) {
	args := m.Called(key, version, permissions, comments)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.ApiKey), args.Error(1)
}

func (m *MockService) DeleteApiKey(apikey *model.ApiKey) (bool, error) {
	args := m.Called(apikey)
	return args.Bool(0), args.Error(1)
}
