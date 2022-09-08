package pgxuuid

import (
	"github.com/google/uuid"
	"github.com/jackc/pgtype"
	pgxUUID "github.com/vgarvardt/pgx-google-uuid/v4"
)

type UUID = pgxUUID.UUID

func New(googleUUID uuid.UUID) UUID {
	return UUID{
		UUID:   googleUUID,
		Status: pgtype.Present,
	}
}

func NewPointer(googleUUID *uuid.UUID) UUID {
	if googleUUID == nil {
		return UUID{
			Status: pgtype.Null,
		}
	}
	return New(*googleUUID)
}

func Parse(s string) (UUID, error) {
	id, err := uuid.Parse(s)
	return New(id), err
}
