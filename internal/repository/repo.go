package repository

import "github.com/wb-go/wbf/dbpg"

type Repository struct {
	db *dbpg.DB
}

func New(db *dbpg.DB) *Repository {
	return &Repository{db: db}
}
