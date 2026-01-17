package user

import (
	"github.com/afteracademy/goserve/api/user/dto"
	"github.com/afteracademy/goserve/api/user/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockService struct {
	mock.Mock
}

func (m *MockService) FetchUserPrivateProfile(user *model.User) (*dto.UserPrivate, error) {
	args := m.Called(user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.UserPrivate), args.Error(1)
}

func (m *MockService) FetchUserPublicProfile(userId uuid.UUID) (*dto.UserPublic, error) {
	args := m.Called(userId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.UserPublic), args.Error(1)
}

func (m *MockService) FetchUserById(id uuid.UUID) (*model.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockService) FetchUserByEmail(email string) (*model.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockService) CreateUser(
	email string, password string, name string, profilePicURL *string, roles []*model.Role,
) (*model.User, error) {
	args := m.Called(email, password, name, profilePicURL, roles)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockService) RemoveUserByEmail(email string) (bool, error) {
	args := m.Called(email)
	return args.Bool(0), args.Error(1)
}

func (m *MockService) FetchRoleByCode(code model.RoleCode) (*model.Role, error) {
	args := m.Called(code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Role), args.Error(1)
}

func (m *MockService) IsEmailExists(email string) (bool, error) {
	args := m.Called(email)
	return args.Bool(0), args.Error(1)
}
