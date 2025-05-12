package gapi

import (
	"database/sql"
	"testing"
	"time"

	db "github.com/gentcod/DummyBank/internal/database"
	"github.com/gentcod/DummyBank/pb"
	"github.com/gentcod/DummyBank/util"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestPbUpdateUserAPI(t *testing.T) {
	testServer := testServerInit(t)

	user, _ := randomUserAndPassword(t)

	newEmail := util.RandomEmail(5)
	newFullName := util.RandomOwner()

	requestBody := &pb.UpdateUserRequest{
		Id:       user.ID.String(),
		FullName: &newFullName,
		Email:    &newEmail,
	}

	arg := db.UpdateUserParams{
		ID: user.ID,
		Email: sql.NullString{
			String: newEmail,
			Valid:  true,
		},
		FullName: sql.NullString{
			String: newFullName,
			Valid:  true,
		},
	}

	updatedUser := db.User{
		ID:              user.ID,
		HarshedPassword: user.HarshedPassword,
		Username:        user.Username,
		FullName:        newFullName,
		Email:           newEmail,
	}

	testServer.mockStore.EXPECT().
		UpdateUser(gomock.Any(), gomock.Eq(arg)).
		Times(1).
		Return(updatedUser, nil)

	ctx := buildContextWithAuth(t, testServer.server.tokenGenerator, user)
	resp, err := testServer.server.UpdateUser(ctx, requestBody)
	require.NoError(t, err)
	require.NotNil(t, resp)
	updatedUserResp := resp.GetUser()
	require.Equal(t, user.Username, updatedUserResp.Username)
	require.Equal(t, updatedUser.FullName, updatedUserResp.FullName)
	require.NotEqual(t, user.FullName, updatedUserResp.FullName)
	require.Equal(t, updatedUser.Email, updatedUserResp.Email)
	require.NotEqual(t, user.Email, updatedUserResp.Email)
	require.WithinDuration(t, user.PasswordChangedAt, updatedUserResp.PasswordChangedAt.AsTime(), time.Second)
	require.WithinDuration(t, user.CreatedAt, updatedUserResp.CreatedAt.AsTime(), time.Second)
	require.WithinDuration(t, updatedUser.UpdatedAt, updatedUserResp.UpdatedAt.AsTime(), time.Second)
}
