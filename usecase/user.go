package usecase

import (
	"context"

	"presensee_project/model/payload"
)

type UserService interface {
	SignUpUser(ctx context.Context, user *payload.UserSignUpRequest) error
	LogInUser(ctx context.Context, user *payload.UserLoginRequest) (string, error)
	GetSingleUser(ctx context.Context, userID uint) (*payload.GetSingleUserResponse, error)
	GetBriefUsers(ctx context.Context, page int, limit int) (*payload.BriefUsersResponse, error)
	UpdateUser(ctx context.Context, userID uint, request *payload.UserUpdateRequest) error
	DeleteUser(ctx context.Context, userID uint) error
}
