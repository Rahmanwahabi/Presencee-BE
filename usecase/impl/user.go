package usecase

import (
	"context"

	"presensee_project/model/payload"
	"presensee_project/repository"
	"presensee_project/usecase"
	"presensee_project/utils"
	"presensee_project/utils/jwt_service"
	"presensee_project/utils/password"

	"github.com/google/uuid"
)

type (
	UserServiceImpl struct {
		userRepository repository.UserRepository
		passwordHash   password.PasswordFunc
		jwtService     jwt_service.JWTService
	}
)

func NewUserServiceImpl(userRepository repository.UserRepository, function password.PasswordFunc, jwt jwt_service.JWTService) usecase.UserService {
	return &UserServiceImpl{
		userRepository: userRepository,
		passwordHash:   function,
		jwtService:     jwt,
	}
}

func (u *UserServiceImpl) SignUpUser(ctx context.Context, user *payload.UserSignUpRequest) error {
	hashedPassword, err := u.passwordHash.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)

	userEntity := user.ToEntity()
	userEntity.ID = uint(uuid.New().ID())

	err = u.userRepository.CreateUser(ctx, userEntity)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserServiceImpl) LogInUser(ctx context.Context, user *payload.UserLoginRequest) (string, error) {
	userEntity, err := u.userRepository.FindByEmail(ctx, user.Email)
	if err != nil {
		if err == utils.ErrUserNotFound {
			return "", utils.ErrInvalidCredentials
		}

		return "", err
	}

	err = u.passwordHash.CompareHashAndPassword([]byte(userEntity.Password), []byte(user.Password))
	if err != nil {
		return "", utils.ErrInvalidCredentials
	}

	token, err := u.jwtService.GenerateToken(userEntity)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (d *UserServiceImpl) GetSingleUser(ctx context.Context, userID uint) (*payload.GetSingleUserResponse, error) {
	user, err := d.userRepository.GetSingleUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	var userResponse = payload.NewGetSingleUserResponse(user)

	return userResponse, nil
}

func (u *UserServiceImpl) GetBriefUsers(ctx context.Context, page int, limit int) (*payload.BriefUsersResponse, error) {
	offset := (page - 1) * limit

	users, err := u.userRepository.GetBriefUsers(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	return payload.NewBriefUsersResponse(users), nil
}

func (u *UserServiceImpl) UpdateUser(ctx context.Context, userID uint, request *payload.UserUpdateRequest) error {
	user := request.ToEntity()
	user.ID = userID

	if user.Password != "" {
		hashedPassword, err := u.passwordHash.GenerateFromPassword([]byte(user.Password), 10)
		if err != nil {
			return err
		}

		user.Password = string(hashedPassword)
	}

	return u.userRepository.UpdateUser(ctx, user)
}

func (d *UserServiceImpl) DeleteUser(ctx context.Context, userID uint) error {

	return d.userRepository.DeleteUser(ctx, userID)
}
