package user

import (
	usersDomain "github.com/Dorrrke/notes-g2/internal/domain/users"
	"github.com/google/uuid"
)

type Repository interface {
	SaveUser(user usersDomain.User) error
	GetUser(login string) (usersDomain.User, error)
}

type UserService struct {
	repo Repository
}

func New(repo Repository) *UserService {
	return &UserService{repo: repo}
}

func (us *UserService) RegisterUser(user usersDomain.User) (string, error) {
	user.UID = uuid.New().String()

	err := us.repo.SaveUser(user)
	if err != nil {
		return ``, err
	}
	return user.UID, nil
}

func (us *UserService) LoginUser(userCreds usersDomain.UserRequest) (string, error) {
	dbUser, err := us.repo.GetUser(userCreds.Email)
	if err != nil {
		return ``, err
	}

	if dbUser.Password != userCreds.Password {
		return ``, usersDomain.ErrInvalidUserCreds
	}

	return dbUser.UID, nil
}
