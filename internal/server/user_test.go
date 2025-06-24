package server

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	usersDomain "github.com/Dorrrke/notes-g2/internal/domain/users"
	"github.com/Dorrrke/notes-g2/internal/server/mocks"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLogin(t *testing.T) {
	var srv NotesAPI

	testRouter := gin.New()

	testRouter.Use(gin.Recovery())

	testRouter.POST("/login", srv.login)

	httpTest := httptest.NewServer(testRouter)
	defer httpTest.Close()

	type want struct {
		resultMsg string
		status    int
	}

	type test struct {
		name    string
		request string
		method  string
		uReq    usersDomain.UserRequest
		dbUser  usersDomain.User
		repoErr error
		want    want
	}

	tests := []test{
		{
			name: "test 1: success call",
			uReq: usersDomain.UserRequest{
				Email:    "email",
				Password: "password",
			},
			dbUser: usersDomain.User{
				UID:      "uuid-1234-55rr",
				Name:     "John Doe",
				Email:    "email",
				Password: "password",
			},
			repoErr: nil,
			request: "/login",
			method:  http.MethodPost,
			want: want{
				resultMsg: "user logined: uuid-1234-55rr",
				status:    200,
			},
		},
		{
			name: "test 2: invalid creds call",
			uReq: usersDomain.UserRequest{
				Email:    "email",
				Password: "password",
			},
			dbUser: usersDomain.User{
				UID:      "uuid-1234-55rr",
				Name:     "John Doe",
				Email:    "email",
				Password: "1234567",
			},
			repoErr: usersDomain.ErrInvalidUserCreds,
			request: "/login",
			method:  http.MethodPost,
			want: want{
				resultMsg: `{"error":"invalid creds"}`,
				status:    401,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := mocks.NewRepository(t)
			mockRepo.On("GetUser", tc.uReq.Email).Return(tc.dbUser, tc.repoErr)
			srv.repo = mockRepo

			req := resty.New().R()
			req.Method = tc.method
			req.URL = httpTest.URL + tc.request // http://localhost:8080/login

			body, err := json.Marshal(tc.uReq)
			assert.NoError(t, err)
			req.Body = body

			resp, err := req.Send()
			assert.NoError(t, err)

			respBody := string(resp.Body())

			assert.Equal(t, tc.want.status, resp.StatusCode())
			assert.Equal(t, tc.want.resultMsg, respBody)
		})
	}
}

func BenchmarkLogin(b *testing.B) {
	var srv NotesAPI

	gin.DefaultWriter = io.Discard
	gin.DisableConsoleColor()
	testRouter := gin.New()

	testRouter.Use(gin.Recovery())

	testRouter.POST("/login", srv.login)

	httpTest := httptest.NewServer(testRouter)
	defer httpTest.Close()

	uReq := usersDomain.UserRequest{
		Email:    "email",
		Password: "password",
	}
	dbUser := usersDomain.User{
		UID:      "uuid-1234-55rr",
		Name:     "John Doe",
		Email:    "email",
		Password: "password",
	}

	mockRepo := mocks.NewRepository(b)
	mockRepo.On("GetUser", uReq.Email).Return(dbUser, nil)
	srv.repo = mockRepo

	req := resty.New().R()
	req.Method = http.MethodPost
	req.URL = httpTest.URL + "/login" // http://localhost:8080/login

	body, err := json.Marshal(uReq)
	assert.NoError(b, err)
	req.Body = body

	// b.ResetTimer()
	for range b.N {
		req.Send()
	}
}

func TestReqister(t *testing.T) {
	var srv NotesAPI

	testRouter := gin.New()

	testRouter.Use(gin.Recovery())

	testRouter.POST("/register", srv.register)

	httpTest := httptest.NewServer(testRouter)
	defer httpTest.Close()

	type want struct {
		resultMsg string
		status    int
	}

	type test struct {
		name    string
		request string
		method  string
		uReq    usersDomain.User
		repoErr error
		want    want
	}

	tests := []test{
		{
			name:    "test 1: success call",
			request: "/register",
			method:  http.MethodPost,
			uReq: usersDomain.User{
				UID:      "uuid-1234-55rr",
				Name:     "John Doe",
				Email:    "email",
				Password: "password",
			},
			repoErr: nil,
			want: want{
				resultMsg: "user registered:",
				status:    200,
			},
		},

		{
			name:    "test 2: conflict call",
			request: "/register",
			method:  http.MethodPost,
			uReq: usersDomain.User{
				UID:      "uuid-1234-55rr",
				Name:     "John Doe",
				Email:    "email",
				Password: "password",
			},
			repoErr: usersDomain.ErrUserAlredyExists,
			want: want{
				resultMsg: `{"error":"user alredy exists"}`,
				status:    409,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := mocks.NewRepository(t)
			mockRepo.On("SaveUser", mock.MatchedBy(func(user usersDomain.User) bool {
				return user.Name == tc.uReq.Name &&
					user.Email == tc.uReq.Email &&
					user.Password == tc.uReq.Password
			})).Return(tc.repoErr)
			srv.repo = mockRepo

			req := resty.New().R()
			req.Method = tc.method
			req.URL = httpTest.URL + tc.request // http://localhost:8080/register

			body, err := json.Marshal(tc.uReq)
			assert.NoError(t, err)
			req.Body = body

			resp, err := req.Send()
			assert.NoError(t, err)

			respBody := string(resp.Body())

			assert.Equal(t, tc.want.status, resp.StatusCode())
			assert.Contains(t, respBody, tc.want.resultMsg)
		})
	}
}
