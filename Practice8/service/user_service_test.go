package service

import (
	"practice-8/repository"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)
func TestGetUserByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := repository.NewMockUserRepository(ctrl)
	userService := NewUserService(mockRepo)
	user := &repository.User{ID: 1, Name: "Bakytzhan Agai"}
	mockRepo.EXPECT().GetUserByID(1).Return(user, nil)
	result, err := userService.GetUserByID(1)
	assert.NoError(t, err)
	assert.Equal(t, user, result)
}
func TestCreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := repository.NewMockUserRepository(ctrl)
	userService := NewUserService(mockRepo)
	user := &repository.User{ID: 1, Name: "Bakytzhan Agai"}
	mockRepo.EXPECT().CreateUser(user).Return(nil)
	err := userService.CreateUser(user)
	assert.NoError(t, err)
}

func TestUserService_RegisterUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	service := NewUserService(mockRepo)
	user := &repository.User{Name: "Test"}

	
	mockRepo.EXPECT().GetByEmail("exists@test.com").Return(&repository.User{}, nil)
	err := service.RegisterUser(user, "exists@test.com")
	assert.EqualError(t, err, "user with this email already exists")

	
	mockRepo.EXPECT().GetByEmail("new@test.com").Return(nil, nil)
	mockRepo.EXPECT().CreateUser(user).Return(nil)
	err = service.RegisterUser(user, "new@test.com")
	assert.NoError(t, err)
}

func TestUserService_UpdateUserName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	service := NewUserService(mockRepo)

	
	user := &repository.User{ID: 2, Name: "OldName"}
	mockRepo.EXPECT().GetUserByID(2).Return(user, nil)
	
	mockRepo.EXPECT().UpdateUser(gomock.Any()).Do(func(u *repository.User) {
		assert.Equal(t, "NewName", u.Name)
	}).Return(nil)

	err := service.UpdateUserName(2, "NewName")
	assert.NoError(t, err)
}

func TestUserService_DeleteUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	service := NewUserService(mockRepo)

	
	err := service.DeleteUser(1)
	assert.EqualError(t, err, "it is not allowed to delete admin user")

	
	mockRepo.EXPECT().DeleteUser(2).Return(nil)
	err = service.DeleteUser(2)
	assert.NoError(t, err)
}