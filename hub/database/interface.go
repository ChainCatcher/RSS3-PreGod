package database

import (
	"context"
	"time"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/database/logger"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/database/model"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var (
	Instance Database
)

type Database interface {
	DB(ctx context.Context) *gorm.DB
	Tx(ctx context.Context) *gorm.DB

	QueryAccount(db *gorm.DB, id string, platformID int) (*model.Account, error)

	QueryAccountPlatforms(db *gorm.DB, id string, platformID int) ([]model.AccountPlatform, error)

	QueryLinks(db *gorm.DB, _type int, identity string, suffixID, pageIndex int) ([]model.Link, error)
	QueryLinksByTarget(db *gorm.DB, _type int, targetIdentity string, targetSuffixID, limit int, instance, lastInstance string) ([]model.Link, error)
	QueryLinkWithMaxPageIndex(db *gorm.DB, _type int, identity string, suffixID int) (int, error)

	QuerySignature(db *gorm.DB, fileURI string) (*model.Signature, error)
}

func Setup() error {
	db, err := gorm.Open(postgres.New(postgres.Config{
		// TODO Refactor config package
		DSN: config.Config.Postgres.DSN,
	}), &gorm.Config{
		SkipDefaultTransaction:                   true,
		NamingStrategy:                           schema.NamingStrategy{SingularTable: true},
		NowFunc:                                  func() time.Time { return time.Now().UTC() },
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger:                                   logger.New(),
	})

	if err != nil {
		return err
	}

	// Install uuid extension for postgres
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";").Error; err != nil {
		return err
	}

	if err := db.AutoMigrate(
		&model.Account{},
		&model.AccountPlatform{},
		&model.Instance{},
		&model.LinkList{},
		&model.Link{},
		&model.Signature{},
		&model.Asset{},
		&model.Note{},
	); err != nil {
		return err
	}

	Instance = &database{
		db: db,
	}

	return nil
}
