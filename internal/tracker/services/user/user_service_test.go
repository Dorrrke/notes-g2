package user

import (
	"testing"

	usersDomain "github.com/Dorrrke/notes-g2/internal/tracker/domain/users"
	"github.com/Dorrrke/notes-g2/internal/tracker/services/user/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLoginUser(t *testing.T) {
	type want struct {
		userID string
		err    error
	}

	type test struct {
		name    string
		user    usersDomain.User
		userReq usersDomain.UserRequest
		want    want
	}

	tests := []test{
		{
			name: "test 1: success call",
			userReq: usersDomain.UserRequest{
				Email:    "email",
				Password: "password",
			},
			user: usersDomain.User{
				UID:      "uuid",
				Name:     "John Doe",
				Email:    "email",
				Password: "password",
			},
			want: want{
				userID: "uuid",
				err:    nil,
			},
		},
		{
			name: "test 2: fail call",
			userReq: usersDomain.UserRequest{
				Email:    "email",
				Password: "password",
			},
			user: usersDomain.User{
				UID:      "uuid",
				Name:     "John Doe",
				Email:    "email",
				Password: "password1234",
			},
			want: want{
				userID: "",
				err:    usersDomain.ErrInvalidUserCreds,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repoMock := mocks.NewRepository(t)
			repoMock.On("GetUser", tc.userReq.Email).Return(tc.user, nil)

			testUserService := New(repoMock)

			userID, err := testUserService.LoginUser(tc.userReq)
			if tc.want.err != nil {
				assert.ErrorIs(t, err, tc.want.err)
				return
			}
			assert.NoError(t, err)

			assert.Equal(t, tc.want.userID, userID)
		})
	}
}

func TestRegisterUser(t *testing.T) {
	type want struct {
		err error
	}

	type test struct {
		name string
		user usersDomain.User
		want want
	}

	tests := []test{
		{
			name: "test 1: success call",
			user: usersDomain.User{
				UID:      "uuid-1234-55rr",
				Name:     "John Doe",
				Email:    "email",
				Password: "password",
			},
			want: want{
				err: nil,
			},
		},
		{
			name: "test 2: unique error case",
			user: usersDomain.User{
				UID:      "uuid-1234-55rr",
				Name:     "John Doe",
				Email:    "email",
				Password: "password",
			},
			want: want{
				err: usersDomain.ErrUserAlredyExists,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repoMock := mocks.NewRepository(t)
			repoMock.On("SaveUser", mock.MatchedBy(func(user usersDomain.User) bool {
				return user.Name == tc.user.Name &&
					user.Email == tc.user.Email &&
					user.Password == tc.user.Password
			})).Return(tc.want.err)

			testUserService := New(repoMock)

			userID, err := testUserService.RegisterUser(tc.user)
			if tc.want.err != nil {
				assert.ErrorIs(t, err, tc.want.err)
				return
			}

			assert.NoError(t, err)

			assert.NotEmpty(t, userID)
		})
	}
}
