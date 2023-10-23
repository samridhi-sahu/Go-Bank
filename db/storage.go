package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/samridhi-sahu/gobank/types"
)

// We created Storage interface so that if in future we want to change our db from postgress to any other than we should not
// write all methods again.
// In this case we just need to implement Storage interface with the new database by providing implementation of the storage methods
type Storage interface {
	CreateAccount(*types.Account) (int, error)
	DeleteAccount(int) error
	GetAllAccounts() ([]*types.Account, error)
	GetAccountByFilter(string, int) (*types.Account, error)
	UpdateBalance(int, int) (*types.Account, error)
}

// PostgressStore will implement the Storage Interface by providing implementation for all the storage methods
type PostgresStore struct {
	db *sql.DB
}

// NewPostgresStore will create connection with the postgres database
func NewPostgresStore() (*PostgresStore, error) {
	connStr := os.Getenv("POSTGRES_URL")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	postgres := PostgresStore{db: db}
	return &postgres, nil
}

// Init will create table in database if it does not exit already
func (s *PostgresStore) Init() error {
	return s.createAccountTable()
}

func (s *PostgresStore) createAccountTable() error {
	query := `create table if not exists account (
		id serial primary key,
		first_name varchar(50),
		last_name varchar(50),
		hashed_password varchar(100),
		number serial,
		balance serial,
		created_at timestamp
	)`

	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) CreateAccount(acc *types.Account) (int, error) {
	query := `insert into account
	(first_name, last_name, hashed_password, number, balance, created_at)
	values($1,$2,$3,$4,$5,$6)
	returning id`

	var id *int
	err := s.db.QueryRow(
		query,
		acc.FirstName,
		acc.LastName,
		acc.HashedPassword,
		acc.Number,
		acc.Balance,
		acc.CreatedAt,
	).Scan(&id)

	if err != nil {
		return -1, err
	}

	return *id, nil
}

func (s *PostgresStore) DeleteAccount(id int) error {
	_, err := s.db.Query("delete from account where id = $1", id)
	return err
}

func (s *PostgresStore) GetAccountByFilter(filter string, val int) (*types.Account, error) {
	query := `select * from account where id = $1`

	if filter == "number" {
		query = `select * from account where number = $1`
	}

	rows, err := s.db.Query(query, val)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoAccount(rows)
	}

	return nil, fmt.Errorf("account with [%s] %d not found", filter, val)
}

func (s *PostgresStore) GetAllAccounts() ([]*types.Account, error) {
	rows, err := s.db.Query("select * from account")
	if err != nil {
		return nil, err
	}

	accounts := []*types.Account{}
	for rows.Next() {
		account, err := scanIntoAccount(rows)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}

	return accounts, nil
}

// scanIntoAccount will return the single row of db by converting it into Account struct
func scanIntoAccount(rows *sql.Rows) (*types.Account, error) {
	account := types.Account{}
	err := rows.Scan(
		&account.ID,
		&account.FirstName,
		&account.LastName,
		&account.HashedPassword,
		&account.Number,
		&account.Balance,
		&account.CreatedAt,
	)
	return &account, err
}

func (s *PostgresStore) UpdateBalance(accNumber, balance int) (*types.Account, error) {
	query := "update account set balance = $1 where number = $2 returning *"
	account := types.Account{}
	err := s.db.QueryRow(query, balance, accNumber).Scan(
		&account.ID,
		&account.FirstName,
		&account.LastName,
		&account.HashedPassword,
		&account.Number,
		&account.Balance,
		&account.CreatedAt,
	)
	if err != nil {
		error := fmt.Errorf("error while update balance: %s", err.Error())
		return nil, error
	}
	return &account, nil
}
