package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"testing"

	db "github.com/gentcod/DummyBank/internal/database"
	"github.com/gentcod/DummyBank/util"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

type eqCreateUserParamsMatcher struct {
	arg      db.CreateUserTxParams
	password string
}

func (e eqCreateUserParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.CreateUserTxParams)
	if !ok {
		return false
	}

	err := util.CheckPassword(e.password, arg.CreateUserParams.HarshedPassword)
	if err != nil {
		return false
	}

	e.arg.CreateUserParams.HarshedPassword = arg.CreateUserParams.HarshedPassword
	e.arg.CreateUserParams.ID = arg.CreateUserParams.ID
	return reflect.DeepEqual(e.arg.CreateUserParams, arg.CreateUserParams)
}

func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

func EqCreateUserParams(arg db.CreateUserTxParams, password string) gomock.Matcher {
	return eqCreateUserParamsMatcher{arg, password}
}

func TestCreateUserAPI(t *testing.T) {
	testServer := testServerInit(t)

	user, password := randomUserAndPassword(t)
	userTxResult := db.CreateUserTxResult{User: user}
	userProfile := getUserProfile(user)

	requestBody := gin.H{
		"username":  user.Username,
		"full_Name": user.FullName,
		"email":     user.Email,
		"password":  password,
	}

	arg := db.CreateUserTxParams{
		CreateUserParams: db.CreateUserParams{
			ID:              user.ID,
			Username:        user.Username,
			FullName:        user.FullName,
			Email:           user.Email,
			HarshedPassword: user.HarshedPassword,
		},
		AfterCreate: func(user db.User) error { return nil },
	}

	jsonBody, err := json.Marshal(requestBody)
	require.NoError(t, err)

	testServer.mockStore.EXPECT().CreateUserTx(gomock.Any(), EqCreateUserParams(arg, password)).Times(1).Return(userTxResult, nil)
	url := "/api/v1/users/signup"
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(jsonBody))
	require.NoError(t, err)

	testServer.server.router.ServeHTTP(testServer.recorder, request)
	require.Equal(t, testServer.recorder.Code, http.StatusOK)
	requireBodyMatchUserProfile(t, testServer.recorder.Body, userProfile)
}

// randomUserAndPassword generates a random account
func randomUserAndPassword(t *testing.T) (user db.User, password string) {
	password = util.RandomStr(9)
	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)

	return db.User{
		ID:              uuid.New(),
		HarshedPassword: hashedPassword,
		Username:        util.RandomOwner(),
		FullName:        util.RandomOwner(),
		Email:           util.RandomEmail(9),
	}, password
}

// requireBodyMatchuserProfile checks if the server recorder body for created user matches the user object
func requireBodyMatchUserProfile(t *testing.T, body *bytes.Buffer, user UserProfile) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var resp ApiResponse[UserProfile]
	err = json.Unmarshal(data, &resp)
	profile := resp.Data
	require.NoError(t, err)
	require.Equal(t, user.FullName, profile.FullName)
	require.Equal(t, user.Email, profile.Email)
}
