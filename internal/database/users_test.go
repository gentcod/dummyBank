package db

import (
	"context"
	"database/sql"
	// "database/sql"
	"testing"
	"time"

	"github.com/gentcod/DummyBank/util"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

//createRandomUser creates a random User for testing and returns a User object
func createRandomUser(t *testing.T) User {
	hashedPassword, err := util.HashPassword(util.RandomStr(10))
	require.NoError(t, err)

	arg := CreateUserParams{
		ID: uuid.New(),
		Username: util.RandomStr(8),
		FullName: util.RandomOwner(),
		Email: util.RandomEmail(9),
		HarshedPassword: hashedPassword,
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.HarshedPassword, user.HarshedPassword)

	require.NotZero(t, user.ID)
	require.True(t, user.PasswordChangedAt.IsZero()) 
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestUpdateUser(t *testing.T) {
	ranUser := createRandomUser(t)
	hashedPassword, err := util.HashPassword(ranUser.HarshedPassword)
	require.NoError(t, err)

	arg := UpdateUserParams{
		ID: ranUser.ID,
		HarshedPassword: sql.NullString{
			String: hashedPassword,
			Valid: true,
		},
		PasswordChangedAt: sql.NullTime{
			Time: time.Now(),
			Valid: true,
		},
	}

	updatedUser, err := testQueries.UpdateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, updatedUser)

	require.Equal(t, ranUser.ID, updatedUser.ID)
	require.Equal(t, ranUser.FullName, updatedUser.FullName)
	require.Equal(t, ranUser.Email, updatedUser.Email)
	require.Equal(t, arg.HarshedPassword.String, updatedUser.HarshedPassword)

	require.WithinDuration(t, ranUser.CreatedAt, updatedUser.CreatedAt, time.Second)
	require.WithinDuration(t, arg.PasswordChangedAt.Time, updatedUser.PasswordChangedAt, time.Second)
}

// func TestGetUserById(t *testing.T) {
// 	user1 := createRandomUser(t)
// 	user2, err := testQueries.GetUserById(context.Background(), user1.ID)
// 	require.NoError(t, err)
// 	require.NotEmpty(t, user2)

// 	require.Equal(t, user1.ID, user2.ID)
// 	require.Equal(t, user1.HarshedPassword, user2.HarshedPassword)
// 	require.Equal(t, user1.FullName, user2.FullName)
// 	require.Equal(t, user1.Email, user2.Email)
// 	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
// 	require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
// }

// func TestGetUsers(t *testing.T) {
// 	for i := 0; i < 10; i++ {
// 		createRandomUser(t)
// 	}
// 	arg := GetUsersParams{
// 		Limit: 5,
// 		Offset: 5,
// 	}

// 	users, err := testQueries.GetUsers(context.Background(), arg)
// 	require.NoError(t, err)
// 	require.Len(t, users, 5)

// 	for _, user := range users {
// 		require.NotEmpty(t, user)
// 	}
// }

// func TestDeleteUser(t *testing.T) {
// 	user1 := createRandomUser(t)
// 	err := testQueries.DeleteUser(context.Background(), user1.ID)
// 	require.NoError(t, err)

// 	user2, err := testQueries.GetUserById(context.Background(), user1.ID)
// 	require.Error(t, err)
// 	require.EqualError(t, err, sql.ErrNoRows.Error())
// 	require.Empty(t, user2)
// }