package configs

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func InitMigrations(env *EnviConfig) error {

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&TimeZone=UTC",
		env.DbUsername,
		env.DbPassword,
		env.DbHost,
		env.DbPort,
		env.DbName,
	)

	// Open database connection
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return err
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	defer driver.Close()

	m, err := migrate.NewWithDatabaseInstance(
		"file://entities/migrations",
		"postgres", driver)
	if err != nil {
		return err
	}

	// Run the migrations
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	log.Println("Run db migrations successfully")

	return nil
}
