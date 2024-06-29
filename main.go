package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"os"

	"github.com/comame/sshproxy/models/sshKeys"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	devEnv := os.Getenv("DEV")
	if devEnv == "" {
		panic("まだ DEV が必要")
	}

	db, err := sql.Open("mysql", "root@unix(./.testdb/mysql.sock)/sshproxy")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	ctx := context.Background()
	con, err := db.Conn(ctx)
	if err != nil {
		panic(err)
	}
	defer con.Close()

	result, err := createAuthorizedKeys(ctx, con)
	if err != nil {
		panic(err)
	}

	log.Println(result)
}

func createAuthorizedKeys(ctx context.Context, con *sql.Conn) (string, error) {
	rows, err := con.QueryContext(ctx, "SELECT authorized_key_id, authenticated_user_id, username, options, pubkey FROM authorized_key")
	if err != nil {
		return "", err
	}

	var ret string

	for rows.Next() {
		var authorizedKeyID int
		var authenticatedUserID, username, optionsJson, pubkey string

		if err := rows.Scan(&authorizedKeyID, &authenticatedUserID, &username, &optionsJson, &pubkey); err != nil {
			return "", err
		}

		var options sshKeys.Option
		if err := json.Unmarshal([]byte(optionsJson), &options); err != nil {
			return "", err
		}

		if !sshKeys.IsValidSSHPublicKey(pubkey) {
			return "", errors.New("invalid format SSH pubkey")
		}

		ret += sshKeys.CreateAuthorizedKeyLine(pubkey, options) + "\n"
	}

	return ret, nil
}
