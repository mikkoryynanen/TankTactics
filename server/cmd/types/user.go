package types

import "github.com/google/uuid"

type User struct {
	Id 			uuid.UUID
	Username 	string
	// CurrentRoom	*Room
}
