package gapi

import (
	"context"
	"fmt"
	"testing"
	"time"

	db "github.com/gentcod/DummyBank/internal/database"
	mockdb "github.com/gentcod/DummyBank/internal/database/mock"
	"github.com/gentcod/DummyBank/token"
	mockwk "github.com/gentcod/DummyBank/worker/mock"
	"google.golang.org/grpc/metadata"

	// "github.com/gentcod/DummyBank/pb"
	"github.com/gentcod/DummyBank/util"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

// TestServer contains all configurations to run mock db tests for the api
type TestPbServer struct {
	server    *Server
	mockStore *mockdb.MockStore
	mockwk    *mockwk.MockTaskDistributor
}

// testServerInit initializes the mockstore, http reponse recorder and the test server.
// It returns an initialized TestServer
func testServerInit(t *testing.T) (testServer TestPbServer) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	testServer.mockStore = mockdb.NewMockStore(ctrl)

	taskCtrl := gomock.NewController(t)
	defer taskCtrl.Finish()
	testServer.mockwk = mockwk.NewMockTaskDistributor(taskCtrl)

	config := util.Config{
		TokenSymmetricKey:   util.RandomStr(32),
		AccessTokenDuration: time.Minute,
	}

	//Start test server and send requests
	server, err := NewServer(config, testServer.mockStore, testServer.mockwk)
	require.NoError(t, err)

	testServer.server = server
	return testServer
}

func buildContext(t *testing.T, generator token.Generator) context.Context {
	return context.Background()
}

func buildContextWithAuth(t *testing.T, generator token.Generator, user db.User) context.Context {
	accessToken, _, err := generator.CreateToken(user.Username, user.ID, time.Minute)
	bearerToken := fmt.Sprintf("%s %s", authorizationTypeBearer, accessToken)
	require.NoError(t, err)
	md := metadata.MD{
		authorizationHeader: []string{
			bearerToken,
		},
	}

	return metadata.NewIncomingContext(context.Background(), md)
}
