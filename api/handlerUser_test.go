package api

import (
	db "github.com/gentcod/DummyBank/internal/database"
	"github.com/gentcod/DummyBank/util"
	"github.com/google/uuid"
)

//randomUser generates a random account
func randomUser() db.User {
	return db.User{
		ID: uuid.New(),
		HarshedPassword: util.RandomStr(9),
		FullName: util.RandomOwner(),
		Email: util.RandomEmail(9),
	}
}
