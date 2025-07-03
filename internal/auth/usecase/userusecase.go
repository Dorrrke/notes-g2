package usecase

import (
	"errors"

	"github.com/Dorrrke/notes-g2/internal/auth/models"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var ErrIvalidPassword = errors.New("invalid password")

type Repository interface {
	SaveUser(models.User) error
	GetUser(string) (models.User, error)
}

type UserUsecase struct {
	repo Repository
}

func NewUserUsecase(repo Repository) *UserUsecase {
	return &UserUsecase{repo: repo}
}

func (u *UserUsecase) Register(user models.User) (string, error) {
	user.UID = uuid.New().String()

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	user.Password = string(hash)

	err = u.repo.SaveUser(user)
	if err != nil {
		return "", err
	}

	return user.UID, nil
}

func (u *UserUsecase) Login(user models.UserRequest) (string, error) {
	dbUser, err := u.repo.GetUser(user.Email)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))
	if err != nil {
		return "", ErrIvalidPassword
	}

	return dbUser.UID, nil
}
