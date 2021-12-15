package biz

import (
	"context"
	"errors"
	v1 "github.com/go-kratos/kratos-layout/api/helloworld/v1"
	"github.com/go-kratos/kratos/v2/log"
	"math/rand"
)

var (
	ErrUsernameInvalid = errors.New("username invalid")
	ErrPasswordInvalid = errors.New("password invalid")
)

type UserDTO struct {
	Username string
	Nickname string
	Avatar   string
}

type User struct {
	Id       int64
	Username string
	Password string
	Nickname string
	Avatar   string
}

func newUser(
	username string,
	password string,
	nickname string,
	avatar string,
) (*User, error) {
	// todo: gen id
	id := rand.Int63()
	// check username simply
	if len(username) == 0 {
		return nil, ErrUsernameInvalid
	}
	// check password simply
	if len(password) == 0 {
		return nil, ErrPasswordInvalid
	}
	// create a user
	return &User{
		Id:       id,
		Username: username,
		Password: password,
		Nickname: nickname,
		Avatar:   avatar,
	}, nil
}

type UserRepo interface {
	Find(ctx context.Context, id int64) (*User, error)
	Save(ctx context.Context, user *User) error
	FindByUsername(ctx context.Context, username string) (*User, error)
}

type UserUseCase struct {
	userRepo UserRepo
	log      *log.Helper
}

func NewUserUseCase(userRepo UserRepo) *UserUseCase {
	return &UserUseCase{
		userRepo: userRepo,
	}
}

func (u *UserUseCase) Registry(ctx context.Context, request *v1.RegisterRequest) (*v1.RegisterReply, error) {

	// is already exists?
	user, err := u.userRepo.FindByUsername(ctx, request.Username)
	if err != nil {
		u.log.Errorf("find user by username error: %v", err)
		return nil, v1.ErrorInternalError("internal error")
	}
	if user != nil {
		return nil, v1.ErrorUserAlreadyExists("username %s already exists", request.Username)
	}

	// create user
	user, err = newUser(request.Username, request.Password, request.Nickname, request.Avatar)
	if err != nil {
		return nil, v1.ErrorRegisterUserFailed("register user failed: %s", err.Error())
	}

	// save user
	err = u.userRepo.Save(ctx, user)
	if err != nil {
		u.log.Errorf("save user error: %v", err)
		return nil, v1.ErrorInternalError("internal error")
	}
	return &v1.RegisterReply{}, nil
}
