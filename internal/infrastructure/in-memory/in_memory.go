package inmemory

import (
	usersDomain "github.com/Dorrrke/notes-g2/internal/domain/users"
	"github.com/Dorrrke/notes-g2/pkg/logger"
)

var emptyUser = usersDomain.User{}

type InMemory struct {
	userStorage map[string]usersDomain.User
}

func New() *InMemory {
	log := logger.Get()

	log.Debug().Msg("create in memory storage")
	return &InMemory{
		userStorage: make(map[string]usersDomain.User),
	}
}

func (im *InMemory) SaveUser(user usersDomain.User) error {
	for _, us := range im.userStorage {
		if us.Email == user.Email {
			return usersDomain.ErrUserAlredyExists
		}
	}

	im.userStorage[user.UID] = user
	return nil
}

func (im *InMemory) GetUser(login string) (usersDomain.User, error) {
	for _, us := range im.userStorage {
		if us.Email == login {
			return us, nil
		}
	}

	return emptyUser, usersDomain.ErrUserNotFound
}
