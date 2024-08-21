package configs

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"

	post "gorm.io/driver/postgres"

	"gorm.io/gorm"
)

func InitDb(env EnviConfig) (*gorm.DB, error) {
	log.Println("create pool database connection")

	dbURL := fmt.Sprintf("host=%v user=%s password=%s dbname=%s port=%v sslmode=disable TimeZone=Asia/Jakarta", env.DbHost, env.DbUsername, env.DbPassword, env.DbName, env.DbPort)
	sqlDb, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Panicln("failed to connect database", err)
		return nil, err
	}

	sqlDb.SetConnMaxIdleTime(30)
	sqlDb.SetMaxOpenConns(50)
	sqlDb.SetConnMaxLifetime(2 * time.Minute)

	log.Println("pool database connection is created")

	ormDb, err := gorm.Open(post.New(post.Config{
		Conn: sqlDb,
	}), &gorm.Config{})

	if err != nil {
		log.Println("error on creating gorm connection")
		return nil, err
	}

	log.Println("gorm connection is created")

	return ormDb, nil
}
