package postgres

import (
	"errors"
)

var (
	errNotFound      = errors.New("not found")
	errNotUpdated    = errors.New("not updated")
	errAlreadyExists = errors.New("already exists")
)

func (p *Postgres) IsNotFoundError(err error) bool {
	return errors.Is(err, errNotFound)
}

func (p *Postgres) IsNotUpdatedError(err error) bool {
	return errors.Is(err, errNotUpdated)
}

func (p *Postgres) IsAlreadyExistsError(err error) bool {
	return errors.Is(err, errAlreadyExists)
}
