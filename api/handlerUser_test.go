package api

import (
	"fmt"
	"net/http"
	"testing"

	db "github.com/gentcod/DummyBank/internal/database"
	"github.com/gentcod/DummyBank/util"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestGetUserByIdAPI(t *testing.T) {
	testServer := testServerInit(t)
	user := randomUser()

	testServer.mockStore.EXPECT().GetUserById(gomock.Any(), gomock.Eq(user.ID)).Return(user, nil)

	url := fmt.Sprintf("/users/%v", user.ID.String())
	request, err := http.NewRequest(http.MethodGet, url, nil)
	require.NoError(t, err)

	testServer.server.router.ServeHTTP(testServer.recorder, request)
	require.Equal(t, testServer.recorder.Code, http.StatusOK)
}

//randomUser generates a random account
func randomUser() db.User {
	return db.User{
		ID: uuid.New(),
		HarshedPassword: util.RandomStr(9),
		FullName: util.RandomOwner(),
		Email: util.RandomEmail(9),
	}
}