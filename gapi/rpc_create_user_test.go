package gapi

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	db "github.com/gentcod/DummyBank/internal/database"
	"github.com/gentcod/DummyBank/pb"
	"github.com/gentcod/DummyBank/util"
	"github.com/gentcod/DummyBank/worker"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

type eqCreateUserParamsMatcher struct {
	arg      db.CreateUserTxParams
	password string
	user     db.User
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
	if !reflect.DeepEqual(e.arg.CreateUserParams, arg.CreateUserParams) {
		return false
	}

	err = arg.AfterCreate(e.user)
	return err == nil
}

func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

func EqCreateUserParams(arg db.CreateUserTxParams, password string, user db.User) gomock.Matcher {
	return eqCreateUserParamsMatcher{arg, password, user}
}

func TestPbCreateUserAPI(t *testing.T) {
	testServer := testServerInit(t)

	user, password := randomUserAndPassword(t)
	userTxResult := db.CreateUserTxResult{User: user}

	requestBody := &pb.CreateUserRequest{
		Username: user.Username,
		FullName: user.FullName,
		Email:    user.Email,
		Password: password,
	}

	arg := db.CreateUserTxParams{
		CreateUserParams: db.CreateUserParams{
			ID:              user.ID,
			Username:        user.Username,
			FullName:        user.FullName,
			Email:           user.Email,
			HarshedPassword: user.HarshedPassword,
		},
	}

	testServer.mockStore.EXPECT().
		CreateUserTx(gomock.Any(), EqCreateUserParams(arg, password, user)).
		Times(1).
		Return(userTxResult, nil)

	taskPayload := &worker.PayloadSendVerifyEmail{
		Username: user.Username,
	}
	testServer.mockwk.EXPECT().
		DistributeTaskSendVerifyEmail(gomock.Any(), taskPayload, gomock.Any()).
		Times(1).
		Return(nil)

	resp, err := testServer.server.CreateUser(context.Background(), requestBody)
	require.NoError(t, err)
	require.NotNil(t, resp)
	createdUser := resp.GetUser()
	require.Equal(t, user.Username, createdUser.Username)
	require.Equal(t, user.FullName, createdUser.FullName)
	require.Equal(t, user.Email, createdUser.Email)
	require.WithinDuration(t, user.PasswordChangedAt, createdUser.PasswordChangedAt.AsTime(), time.Second)
	require.WithinDuration(t, user.CreatedAt, createdUser.CreatedAt.AsTime(), time.Second)
	require.WithinDuration(t, user.UpdatedAt, createdUser.UpdatedAt.AsTime(), time.Second)
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
