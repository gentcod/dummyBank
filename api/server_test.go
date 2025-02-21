package api

import (
	"net/http/httptest"
	"os"
	"testing"
	"time"

	mockdb "github.com/gentcod/DummyBank/internal/database/mock"
	"github.com/gentcod/DummyBank/util"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

//TestServer contains all configurations to run mock db tests for the api
type TestServer struct {
	server *Server
	recorder *httptest.ResponseRecorder
	mockStore *mockdb.MockStore
}

//testServerInit initializes the mockstore, http reponse recorder and the test server.
//It returns an initialized TestServer 
func testServerInit(t *testing.T) (testServer TestServer) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	testServer.mockStore = mockdb.NewMockStore(ctrl)

	config := util.Config{
		TokenSymmetricKey: util.RandomStr(32),
		AccessTokenDuration: time.Minute,
	}

	//Start test server and send requests
	server, err := NewServer(config, testServer.mockStore, nil)
	require.NoError(t, err)

	testServer.server = server
	testServer.recorder = httptest.NewRecorder()
	return testServer
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	//Initialize connection test, terminate test if error occurs
	os.Exit(m.Run()) 
}