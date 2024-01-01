package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	// "net/http/httptest"
	"testing"

	db "github.com/gentcod/DummyBank/internal/database"
	// mockdb "github.com/gentcod/DummyBank/internal/database/mock"
	"github.com/gentcod/DummyBank/util"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestGetAccountByIdAPI(t *testing.T) {
	testServerInit(t, &testServer)
	account := randomAccount()

	//Build stubs
	testServer.mockStore.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(account, nil)

	url := fmt.Sprintf("/accounts/%v", account.ID.String())
	request, err := http.NewRequest(http.MethodGet, url, nil)
	require.NoError(t, err)

	testServer.server.router.ServeHTTP(testServer.recorder, request)
	require.Equal(t, http.StatusOK, testServer.recorder.Code)
	requireBodyMatchAccount(t, testServer.recorder.Body, account)
}

func TestGetAccountsAPI(t *testing.T) {
	testServerInit(t, &testServer)

	var pageId int32 = 1
	var pageSize int32 = 10
	arg := db.GetAccountsParams{
		Limit: pageSize,
		Offset: (pageId - 1) * pageSize,
	}

	accounts := randomAccounts(int(pageSize))

	testServer.mockStore.EXPECT().GetAccounts(gomock.Any(), gomock.Eq(arg)).Times(1).Return(accounts, nil)

	url := fmt.Sprintf("/accounts?page_id=%v&page_size=%v", pageId, pageSize)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	require.NoError(t, err)

	testServer.server.router.ServeHTTP(testServer.recorder, request)
	require.Equal(t, http.StatusOK, testServer.recorder.Code)
	requireBodyMatchAccounts(t, testServer.recorder.Body, accounts, int(pageSize))
}

//TODO: Implement code refractoring for test

//randomAccount generates a random account
func randomAccount() db.Account {
	return db.Account{
		ID: uuid.New(),
		Owner: util.RandomOwner(),
		Balance: util.RandomMoney(),
		Currency: util.RandomCur(),
	}
}

//randomAccounts generates random accounts
func randomAccounts(num int) []db.Account {
	var accounts []db.Account

	for i := 0; i < int(num); i++ {
		accounts = append(accounts, randomAccount())
	}
	return accounts
}

//requireBodyMatchAccount checks if the server recorder body matches the account object
func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account db.Account) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var getAccount db.Account
	err = json.Unmarshal(data, &getAccount)
	require.NoError(t, err)
	require.Equal(t, account, getAccount)
}

//requireBodyMatchAccounts checks if the server recorder body matches the accounts object
func requireBodyMatchAccounts(t *testing.T, body *bytes.Buffer, accounts []db.Account, num int) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var getAccount []db.Account
	err = json.Unmarshal(data, &getAccount)
	require.NoError(t, err)

	for i := 0; i < num; i++ {
		require.Equal(t, accounts[i], getAccount[i])
	}
}