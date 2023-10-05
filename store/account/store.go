package account

import (
	"database/sql"
	cr "errors"
	"fmt"

	"github.com/shivasaicharanruthala/webapp/errors"
	"github.com/shivasaicharanruthala/webapp/model"
	"github.com/shivasaicharanruthala/webapp/store"
)

type accountStore struct {
	DB *sql.DB
}

func New(db *sql.DB) store.Account {
	return &accountStore{DB: db}
}

const ACCOUNTS_TABLE_NAME string = "accounts"

func (acc *accountStore) Insert(account *model.Account) (*model.Account, error) {
	return nil, nil
}

func (acc *accountStore) BulkInsert(cols []string, rows [][]string) error {
	insertQuery := generateBulkInsertQuery(ACCOUNTS_TABLE_NAME, cols, len(rows))

	// Insert the record into the database (replace with your SQL)
	values := flattenBulkUserRecord(rows)
	_, err := acc.DB.Exec(insertQuery, values...)
	if err != nil {
		return errors.NewCustomError(err)
	}

	return nil
}

func (acc *accountStore) IsAccountExists(email string) (*model.Account, error) {
	var account model.Account

	row := acc.DB.QueryRow("SELECT id, email, password FROM accounts WHERE email=$1", email)
	if row.Err() != nil {
		return nil, errors.NewCustomError(row.Err())
	}

	err := row.Scan(&account.ID, &account.Email, &account.Password)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, errors.NewCustomError(cr.New("Username does not exists"))
		}

		return nil, errors.NewCustomError(err)
	}

	return &account, nil
}

func (acc *accountStore) FlushData(tableName string) error {
	deleteQuery := fmt.Sprintf("DELETE FROM %s", tableName)

	// Execute the delete query
	_, err := acc.DB.Exec(deleteQuery)
	if err != nil {
		return errors.NewCustomError(err)
	}

	fmt.Println("All data deleted from the table.")
	return nil
}
