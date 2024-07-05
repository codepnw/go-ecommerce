package appinfoRepositories

import "database/sql"

type IAppinfoRepository interface {

}

type appinfoRepository struct {
	db *sql.DB
}

func AppinfoRepository(db *sql.DB) IAppinfoRepository {
	return &appinfoRepository{db: db}
}