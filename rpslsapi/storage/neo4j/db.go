package neo4j

import (
	"fmt"
	"io"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/neo4j"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/rs/zerolog/log"
	"rpsls/rpslsapi"
)

type DbClient struct {
	driver       neo4j.Driver
	databaseName string
}

func NewDbClient() (DbClient, func()) {
	db, err := neo4j.NewDriver(rpslsapi.Config.DB.Uri,
		neo4j.BasicAuth(rpslsapi.Config.DB.Username, rpslsapi.Config.DB.Password, rpslsapi.Config.DB.Realm))
	if err != nil {
		panic(fmt.Errorf("failed to connect to the database: %v", err))
	}

	if err = migrateDB(db); err != nil && err != migrate.ErrNoChange {
		panic(fmt.Errorf("failed to migrate database: %v", err))
	}

	cleanup := func() {
		CloseDBResource(db)
	}

	return DbClient{
		driver:       db,
		databaseName: rpslsapi.Config.DB.Database,
	}, cleanup
}

func CloseDBResource(closer io.Closer) {
	if err := closer.Close(); err != nil {
		log.Fatal().Err(fmt.Errorf("could not close resource: %w", err))
	}
}

func migrateDB(driver neo4j.Driver) error {
	target := driver.Target()
	url := fmt.Sprintf("%s://%s:%s@%s/%s",
		target.Scheme,
		rpslsapi.Config.DB.Username,
		rpslsapi.Config.DB.Password,
		target.Host,
		rpslsapi.Config.DB.Database,
	)
	m, err := migrate.New("file://db/migrations", url)
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil {
		return err
	}

	return nil
}
