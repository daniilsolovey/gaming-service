package database

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"

	"github.com/daniilsolovey/gaming-task/internal/requester"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/reconquest/karma-go"
	"github.com/reconquest/pkg/log"
)

type Database struct {
	name     string
	host     string
	port     string
	user     string
	password string
	client   *pgxpool.Pool
}

func NewDatabase(
	name, host, port, user, password string,
) *Database {
	database := &Database{
		name:     name,
		host:     host,
		user:     user,
		password: password,
	}

	connection, err := database.connect()
	if err != nil {
		log.Fatal(err)
	}

	database.client = connection

	return database
}

type Player struct {
	ID          string
	NickName    sql.NullString
	BankGroupID sql.NullString
	Balance     float64
}

func (database *Database) connect() (*pgxpool.Pool, error) {
	databaseURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		database.user,
		database.password,
		database.host,
		database.port,
		database.name,
	)
	connection, err := pgxpool.Connect(context.Background(), databaseURL)
	if err != nil {
		return nil, karma.Format(
			err,
			"unable to connect to the database: %s",
			database.name,
		)
	}

	return connection, nil
}

func (database *Database) Close() error {
	database.client.Close()
	return nil
}

func (database *Database) CreateTables() error {
	log.Infof(
		karma.Describe("database", database.name),
		"create tables in database",
	)

	log.Info("creating bank_group table")
	_, err := database.client.Query(
		context.Background(),
		SQL_CREATE_TABLE_BANK_GROUP,
	)
	if err != nil {
		return karma.Format(
			err,
			"unable to create bank_group table in the database: %s",
			database.name,
		)
	}

	log.Info("bank_group table created")

	log.Info("creating player table")
	_, err = database.client.Query(
		context.Background(),
		SQL_CREATE_TABLE_PLAYER,
	)
	if err != nil {
		return karma.Format(
			err,
			"unable to create player table in the database: %s",
			database.name,
		)
	}

	log.Info("player table created")
	return nil
}

func (database *Database) InsertPlayer(player requester.ResponseCreatePlayer, balance float64) error {
	log.Infof(nil, "inserting player to database: %v", player)

	rows, err := database.client.Query(
		context.Background(),
		SQL_INSERT_PLAYER,
		strconv.Itoa(player.ID),
		balance,
	)
	defer func() {
		rows.Close()
	}()
	if err != nil {
		return karma.Format(
			err,
			"unable to add player to the database,"+
				" player_id: %d, balance: %f",
			player.ID, balance,
		)
	}

	log.Info("player successfully added")

	return nil
}

func (database *Database) GetPlayerByID(playerID string) (*Player, error) {
	log.Infof(nil, "receiving player from the database: %s", playerID)
	row := database.client.QueryRow(
		context.Background(),
		SQL_SELECT_PLAYER_BY_ID,
		playerID,
	)

	var (
		player Player
	)

	err := row.Scan(
		&player.ID,
		&player.NickName,
		&player.BankGroupID,
		&player.Balance,
	)

	if err != nil {
		return nil, karma.Format(
			err,
			"error during scaning player result from database rows",
		)
	}

	log.Info("player successfully received")
	return &player, nil
}

func (database *Database) UpdatePlayerBalance(playerID string, balance float64) error {
	log.Infof(nil, "updating player_balance, player_id: %s, new_balance: %f", playerID, balance)

	rows, err := database.client.Query(
		context.Background(),
		SQL_UPDATE_PLAYER_BALANCE,
		balance,
		playerID,
	)

	defer func() {
		rows.Close()
	}()

	if err != nil {
		return karma.Format(
			err,
			"unable to update player_balance in the database,"+
				" player_id: %s",
			playerID,
		)
	}

	log.Infof(nil, "player_balance successfully updated, player_id: %s", playerID)
	return nil
}
