package api

import (
	"net/http/httptest"
	"os"
	"testing"

	mockdb "github.com/gentcod/DummyBank/internal/database/mock"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
)

type TestServer struct {
	server *Server
	recorder *httptest.ResponseRecorder
	mockStore *mockdb.MockStore
}

var testServer TestServer

//testServerInit initializes the mockstore, http reponse recorder and the test server
func testServerInit(t *testing.T, testServer *TestServer) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	testServer.mockStore = mockdb.NewMockStore(ctrl)

	//Start test server and send requests
	testServer.server = NewServer(testServer.mockStore)
	testServer.recorder = httptest.NewRecorder()
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	//Initialize connection test, terminate test if error occurs
	os.Exit(m.Run()) 
}