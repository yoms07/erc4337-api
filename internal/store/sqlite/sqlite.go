package sqlite_store

import (
	"database/sql"
	"web3-account-abstraction-api/internal/model"
	"web3-account-abstraction-api/internal/store"

	_ "github.com/mattn/go-sqlite3"
)

type sqliteStore struct {
	db *sql.DB
}

func (s sqliteStore) CountWallet() (int, error) {
	row := s.db.QueryRow(`
		SELECT count(address) FROM wallet
	`)
	var count int
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (s sqliteStore) CreateWallet(wallet model.UserWallet) error {
	stmt, err := s.db.Prepare(`
		INSERT INTO wallet(address) VALUES (?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(wallet.Sender)
	return err
}

func (s sqliteStore) GetAllWallet() ([]model.UserWallet, error) {
	rows, err := s.db.Query(`
		SELECT address FROM wallet
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := []model.UserWallet{}
	for rows.Next() {
		var address string
		err = rows.Scan(&address)
		if err != nil {
			return nil, err
		}
		result = append(result, model.UserWallet{Sender: address})
	}
	return result, nil
}

func (s sqliteStore) GetWallet(sender string) (model.UserWallet, error) {
	stmt, err := s.db.Prepare(`
		SELECT address FROM wallet
		WHERE address = ?
	`)
	if err != nil {
		return model.UserWallet{}, err
	}
	defer stmt.Close()
	var addr string
	err = stmt.QueryRow(sender).Scan(&addr)
	if err != nil {
		return model.UserWallet{}, err
	}

	return model.UserWallet{
		Sender: addr,
	}, nil
}

func NewStore(db *sql.DB) store.Store {
	return sqliteStore{db: db}
}
