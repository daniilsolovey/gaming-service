package database

const (
	SQL_CREATE_TABLE_BANK_GROUP = `
	CREATE TABLE IF NOT EXISTS
	bank_group(
		id VARCHAR(50) UNIQUE NOT NULL PRIMARY KEY,
		currency VARCHAR(50)
	);
`

	SQL_CREATE_TABLE_PLAYER = `
	CREATE TABLE IF NOT EXISTS
	player(
		player_id VARCHAR(50) UNIQUE NOT NULL PRIMARY KEY,
		nick_name VARCHAR(50),
		bank_group_id VARCHAR(50),
		balance DECIMAL
);
`

	SQL_CREATE_TABLE_SESSION = `
	CREATE TABLE IF NOT EXISTS
	session(
	id serial PRIMARY KEY,
	player_id VARCHAR(50),
	game_id VARCHAR(50),
);
`

	SQL_INSERT_PLAYER = `
	INSERT INTO
	player(
		player_id,
		balance
	)
	VALUES($1, $2);
	`

	SQL_SELECT_PLAYER_BY_ID = `
	SELECT * FROM player
	WHERE player_id = $1;
`

	SQL_UPDATE_PLAYER_BALANCE = `
	UPDATE player
		SET
			balance = $1
	WHERE player.player_id = $2;
	`
)
