package usecase

import (
	"7-solutions/model"
	"7-solutions/repository"
	"7-solutions/utils"
	"context"
	"errors"
)

type UserUsecase interface {
	Register(ctx context.Context, name, email, password string) error
	Login(ctx context.Context, email, password string) (string, error)
	GetUser(ctx context.Context, id string) (*model.User, error)
	ListUsers(ctx context.Context) ([]model.User, error)
	UpdateUser(ctx context.Context, id string, name, email string) error
	DeleteUser(ctx context.Context, id string) error
	CountUsers(ctx context.Context) (int64, error)
}

type userUsecase struct {
	repo      repository.UsersRepository
	jwtSecret string
}

func NewUserUsecase(repo repository.UsersRepository, jwtSecret string) UserUsecase {
	return &userUsecase{repo: repo, jwtSecret: jwtSecret}
}

func (u *userUsecase) Register(ctx context.Context, name, email, password string) error {
	hashed, err := utils.HashPassword(password)
	if err != nil {
		return err
	}
	user := &model.User{
		Name:     name,
		Email:    email,
		Password: hashed,
	}
	return u.repo.Create(ctx, user)
}

func (u *userUsecase) Login(ctx context.Context, email, password string) (string, error) {
	user, err := u.repo.GetByEmail(ctx, email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}
	if !utils.CheckPasswordHash(password, user.Password) {
		return "", errors.New("invalid credentials")
	}
	token, err := utils.GenerateJWT(user.ID.Hex(), u.jwtSecret)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (u *userUsecase) GetUser(ctx context.Context, id string) (*model.User, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *userUsecase) ListUsers(ctx context.Context) ([]model.User, error) {
	return u.repo.List(ctx)
}

func (u *userUsecase) UpdateUser(ctx context.Context, id string, name, email string) error {
	return u.repo.Update(ctx, id, &model.User{Name: name, Email: email})
}

func (u *userUsecase) DeleteUser(ctx context.Context, id string) error {
	return u.repo.Delete(ctx, id)
}

func (u *userUsecase) CountUsers(ctx context.Context) (int64, error) {
	return u.repo.Count(ctx)
}
