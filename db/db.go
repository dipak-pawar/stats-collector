package db

import (
	"database/sql"
	"log"
	"time"

	"github.com/baijum/pgmigration"
	"github.com/jpillora/backoff"
	// DB is actually initialized here
	_ "github.com/lib/pq"
	"os"
)

// Connect to database
func Connect(conf string) *sql.DB {
	var err error
	DB, err := sql.Open("postgres", conf)
	if err != nil {
		log.Fatal(err)
	}

	_, err = DB.Exec("SELECT 1") // DB.Ping() seems to be not working always

	b := &backoff.Backoff{
		Min:    7 * time.Second,
		Factor: 2,
		Max:    7 * time.Minute,
	}

	go func() {
		for {
			_, err = DB.Exec("SELECT 1") // DB.Ping() seems to be not working always
			if err != nil {
				d := b.Duration()
				log.Printf("%s (pinging failed), reconnecting in %s", err, d)
				time.Sleep(d)
				continue
			}
			b.Reset()
			time.Sleep(b.Max)
		}
	}()
	return DB
}

func MigrateDatabase(DB *sql.DB) {
	var err error
	defer handleExit()

	b := &backoff.Backoff{
		Min:    7 * time.Second,
		Factor: 2,
		Max:    7 * time.Minute,
	}

	for i := 0; i < 7; i++ {
		_, err = DB.Exec("SELECT 1") // DB.Ping() seems to be not working always
		if err != nil {
			d := b.Duration()
			log.Printf("%s (pinging failed), reconnecting in %s", err, d)
			time.Sleep(d)
			continue
		}
		b.Reset()
	}

	_, err = DB.Exec("SELECT 1") // DB.Ping() seems to be not working always
	if err != nil {
		log.Println("Migration failed.", err.Error())
		panic(Exit{1})
	}

	err = SchemaMigrate(DB)
	if err != nil {
		log.Println("Migration failed.", err.Error())
		panic(Exit{1})
	}
	log.Println("Migration completed.")
}

// Exit code for clean exit
type Exit struct {
	Code int
}

// exit code handler
func handleExit() {
	if e := recover(); e != nil {
		if exit, ok := e.(Exit); ok {
			os.Exit(exit.Code)
		}
		panic(e) // not an Exit, bubble up
	}
}

// SchemaMigrate migrate database schema
func SchemaMigrate(DB *sql.DB) error {
	return pgmigration.Migrate(DB, AssetNames, Asset, nil)
}
