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
	arg	db.CreateUserParams
	password string
}

func (e eqCreateUserParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.CreateUserParams)
	if !ok {
		return false
	}

	err := util.CheckPassword(e.password, arg.HarshedPassword)
	if err != nil {
		return false
	}

	e.arg.HarshedPassword = arg.HarshedPassword
	e.arg.ID = arg.ID
	return reflect.DeepEqual(e.arg, arg)
}

func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

func EqCreateUserParams(arg db.CreateUserParams, password string) gomock.Matcher { return eqCreateUserParamsMatcher{arg, password} }


func TestCreateUserAPI(t *testing.T) {
	testServer := testServerInit(t)

	user, password := randomUserAndPassword(t)

	requestBody := gin.H{
		"username": user.Username,
		"full_Name": user.FullName,
		"email": user.Email,
		"password": password,
	}

	arg := db.CreateUserParams{
		ID: uuid.New(),
		Username: user.Username,
		FullName: user.FullName,
		Email: user.Email,
		HarshedPassword: user.HarshedPassword,
	}

	jsonBody, err := json.Marshal(requestBody)
	require.NoError(t, err)

	testServer.mockStore.EXPECT().CreateUser(gomock.Any(), EqCreateUserParams(arg, password)).Times(1).Return(user, nil)
	url := "/api/v1/users/signup"
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(jsonBody))
	require.NoError(t, err)

	testServer.server.router.ServeHTTP(testServer.recorder, request)
	require.Equal(t, testServer.recorder.Code, http.StatusOK)
	requireBodyMatchCreaterUser(t, testServer.recorder.Body, user)
}

//randomUserAndPassword generates a random account
func randomUserAndPassword(t *testing.T) (user db.User, password string) {
	password = util.RandomStr(9)
	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)

	return db.User{
		ID: uuid.New(),
		HarshedPassword: hashedPassword,
		Username: util.RandomOwner(),
		FullName: util.RandomOwner(),
		Email: util.RandomEmail(9),
	}, password
}

//requireBodyMatchCreateruser checks if the server recorder body for creatUser matches the user object
func requireBodyMatchCreaterUser(t *testing.T, body *bytes.Buffer, user db.User) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var getUser createUserRequest
	err = json.Unmarshal(data, &getUser)
	require.NoError(t, err)
	require.Equal(t, user.FullName, getUser.FullName)
	require.Equal(t, user.Email, getUser.Email)
}
