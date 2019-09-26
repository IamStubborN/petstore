package mockdb

import "github.com/IamStubborN/petstore/config"

type Database struct {
}

func (d Database) InitDatabase(cfg config.DB) *Database {
	return &Database{}
}

func (d *Database) Close() error {
	panic("not implemented")
}
