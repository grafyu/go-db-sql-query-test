package main

import (
	"database/sql"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func Test_SelectClient_WhenOk(t *testing.T) {
	db, err := sql.Open("sqlite", "demo.db")
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()

	clientID := 1

	resSelect, errSelect := selectClient(db, clientID) // (Client, error)
	assert.NoError(t, errSelect)

	assert.Equal(t, resSelect.ID, 1)

	// id, fio, login, birthday, email
	assert.NotEmpty(t, resSelect.FIO)
	assert.NotEmpty(t, resSelect.Login)
	assert.NotEmpty(t, resSelect.Birthday)
	assert.NotEmpty(t, resSelect.Email)
}

func Test_SelectClient_WhenNoClient(t *testing.T) {
	db, err := sql.Open("sqlite", "demo.db")
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()

	clientID := -1

	resSelect, errSelect := selectClient(db, clientID) // (Client, error)
	assert.Equal(t, errSelect, sql.ErrNoRows)

	assert.Empty(t, resSelect.FIO)
	assert.Empty(t, resSelect.Login)
	assert.Empty(t, resSelect.Birthday)
	assert.Empty(t, resSelect.Email)
}

func Test_InsertClient_ThenSelectAndCheck(t *testing.T) {
	db, err := sql.Open("sqlite", "demo.db")
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()

	cl := Client{
		FIO:      "Test",
		Login:    "Test",
		Birthday: "19700101",
		Email:    "mail@mail.com",
	}

	resInsert, errInsert := insertClient(db, cl) // (int, error)
	cl.ID = resInsert

	require.NotEmpty(t, resInsert)
	require.Empty(t, errInsert)

	resSelect, errSelect := selectClient(db, cl.ID) // (Client, error)
	require.Empty(t, errSelect)
	assert.Equal(t, resSelect, cl)
}

func Test_InsertClient_DeleteClient_ThenCheck(t *testing.T) {
	db, err := sql.Open("sqlite", "demo.db")
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()

	cl := Client{
		FIO:      "Test",
		Login:    "Test",
		Birthday: "19700101",
		Email:    "mail@mail.com",
	}

	resInsert, errInsert := insertClient(db, cl) // (int, error)

	require.NotEmpty(t, resInsert)
	require.Empty(t, errInsert)

	_, errSelect := selectClient(db, resInsert) // (Client, error)
	require.NoError(t, errSelect)

	errDelete := deleteClient(db, resInsert)
	require.NoError(t, errDelete)

	_, errSelect = selectClient(db, resInsert)
	require.Equal(t, errSelect, sql.ErrNoRows)
}
