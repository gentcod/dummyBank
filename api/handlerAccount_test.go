package api

import (
	// "bytes"
	// "encoding/json"
	"fmt"
	// "io"
	"net/http"
	"time"

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
	testServer:= testServerInit(t)
	account, user := randomAccount(t)

	//Build stubs
	testServer.mockStore.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(account, nil)

	url := fmt.Sprintf("/accounts/%v", account.ID.String())
	request, err := http.NewRequest(http.MethodGet, url, nil)
	require.NoError(t, err)

	addAuthorization(t, request, testServer.server.tokenGenerator, authorizationTypeBearer, user.Username, time.Minute)
	require.Equal(t, http.StatusOK, testServer.recorder.Code)
}

func TestGetAccountsAPI(t *testing.T) {
	testServer := testServerInit(t)

	var pageId int32 = 1
	var pageSize int32 = 10
	accounts, user := randomAccounts(int(pageSize), t)

	arg := db.GetAccountsParams{
		Owner: user.ID,
		Limit: pageSize,
		Offset: (pageId - 1) * pageSize,
	}

	testServer.mockStore.EXPECT().GetAccounts(gomock.Any(), gomock.Eq(arg)).Times(1).Return(accounts, nil)

	url := fmt.Sprintf("/accounts?page_id=%v&page_size=%v", pageId, pageSize)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	require.NoError(t, err)

	addAuthorization(t, request, testServer.server.tokenGenerator, authorizationTypeBearer, user.Username, time.Minute)
	require.Equal(t, http.StatusOK, testServer.recorder.Code)
}

//TODO: Implement code refractoring for different test cases

//randomAccount generates a random account
func randomAccount(t *testing.T) (account db.Account, user db.User) {
	user, password := randomUserAndPassword(t)
	
	if err := util.CheckPassword(password, user.HarshedPassword); err != nil {
		return
	}

	return db.Account{
		ID: uuid.New(),
		Owner: user.ID,
		Balance: util.RandomMoney(),
		Currency: util.RandomCur(),
	}, user
}

//randomAccounts generates random accounts
func randomAccounts(num int, t *testing.T) ([]db.Account, db.User) {
	var accounts []db.Account
	user, _ := randomUserAndPassword(t)

	for i := 0; i < int(num); i++ {
		account, _ := randomAccount(t)
		account.Owner = user.ID
		accounts = append(accounts, account)
	}
	return accounts, user
}

// //requireBodyMatchAccount checks if the server recorder body matches the account object
// func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account db.Account) {
// 	data, err := io.ReadAll(body)
// 	require.NoError(t, err)

// 	var getAccount db.Account
// 	err = json.Unmarshal(data, &getAccount)
// 	require.NoError(t, err)
// 	require.Equal(t, account, getAccount)
// }

// //requireBodyMatchAccounts checks if the server recorder body matches the accounts object
// func requireBodyMatchAccounts(t *testing.T, body *bytes.Buffer, accounts []db.Account, num int) {
// 	data, err := io.ReadAll(body)
// 	require.NoError(t, err)

// 	var getAccount []db.Account
// 	err = json.Unmarshal(data, &getAccount)
// 	require.NoError(t, err)

// 	for i := 0; i < num; i++ {
// 		require.Equal(t, accounts[i], getAccount[i])
// 	}
// }