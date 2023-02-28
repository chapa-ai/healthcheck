package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/DavidHuie/gomigrate"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"healthcheck/pkg/models"
	"os"
	"sync"
)

type Config struct {
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
}

var (
	db         *sql.DB
	dbConnOnce sync.Once
)

func GetDB() (*sql.DB, error) {

	var err error

	dbConnOnce.Do(func() {
		host := os.Getenv("DBHost")
		user := os.Getenv("DBUser")
		password := os.Getenv("DBPassword")
		dbName := os.Getenv("DBName")

		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", host, user, password, dbName)

		d, _err := sql.Open("postgres", dsn)
		if _err != nil {
			logrus.Printf("sql.Open failed: %s\n", _err)
			err = _err
			return
		}
		db = d

	})

	return db, err
}

func MigrateDb(path string) error {
	db, err := GetDB()
	if err != nil {
		logrus.Printf("failed GetDB(): %s", err)
		return err
	}
	migrator, err := gomigrate.NewMigrator(db, gomigrate.Postgres{}, path)
	if err != nil {
		logrus.Printf("failed implement migrations: %s", err)
		return err
	}

	return migrator.Migrate()
}

func SelectResultsInDB(ctx context.Context, d *sql.DB, results *models.Response) (*models.Response, error) {
	res := &models.Response{}

	sqlStmt := `SELECT url, statuscode, text FROM results WHERE "url"=$1`
	err := d.QueryRowContext(ctx, sqlStmt, results.Url).Scan(&res.Url, &res.StatusCode, &res.Text)
	if err != nil {
		logrus.Infof("queryRow: %s", err)
	}
	return res, nil
}

func SendResultsOfChecksToDb(ctx context.Context, connDb *sql.DB, results *models.Response) (string, string, error) {
	res, err := SelectResultsInDB(ctx, connDb, results)
	if err != nil {
		logrus.Errorf("failed SelectResultsInDB: %s", err)
		return "", "", err
	}

	stmt := `INSERT INTO results ("url", "statuscode", "text") VALUES ($1, $2, $3) ON CONFLICT(url)
		 DO UPDATE SET "url" = $1, "statuscode" = $2, "text" = $3 RETURNING "url", "statuscode", "text"`

	output := &models.Response{}
	err = connDb.QueryRowContext(ctx, stmt, results.Url, results.StatusCode, results.Text).Scan(&output.Url, &output.StatusCode, &output.Text)
	if err != nil {
		logrus.Errorf("failed QueryRow(): %s", err)
		return "", "", err
	}

	return res.Text, output.Text, nil
}

func GetFailedChecks(ctx context.Context) ([]*models.Response, error) {
	sqlStmt := `SELECT url, statuscode, text FROM results WHERE "statuscode" != 200 OR "text" = 'failed'`
	rows, err := db.QueryContext(ctx, sqlStmt)
	if err != nil {
		logrus.Errorf("failed db.Query: %s", err)
		return nil, err
	}

	list := []*models.Response{}
	for rows.Next() {
		obj := models.Response{}
		if err := rows.Scan(&obj.Url, &obj.StatusCode, &obj.Text); err != nil {
			logrus.Errorf("failed rows.Scan: %s", err)
			return nil, err
		}

		list = append(list, &obj)
	}
	if err = rows.Err(); err != nil {
		logrus.Errorf("failed rows.Err: %s", err)
		return nil, err
	}

	return list, nil
}
