package migrations

import (
	"reflect"

	logging "github.com/payly-solucoes-de-pagamentos/golang-logging"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

type IMigration interface {
	Migrate(client *gorm.DB, migration interface{})
	RollbackToId(client *gorm.DB, migration interface{}, id string)
}

type RepositoryMigrations struct {
	logger *logging.Logger
}

func GetMigrationsFromType(migration interface{}) []*gormigrate.Migration {

	migrations := make([]*gormigrate.Migration, 0)

	migrationType := reflect.TypeOf(migration)
	migrationValue := reflect.ValueOf(migration)

	for i := 0; i < migrationType.NumMethod(); i++ {
		method := migrationType.Method(i)
		ret := method.Func.Call([]reflect.Value{migrationValue})
		returnVal := ret[0].Interface()
		migrations = append(migrations, returnVal.(*gormigrate.Migration))
	}

	return migrations

}

func (repositoryMigrations RepositoryMigrations) Migrate(client *gorm.DB, migration interface{}) {

	migrations := GetMigrationsFromType(migration)

	migrationTransaction := gormigrate.New(client, gormigrate.DefaultOptions, migrations)

	if err := migrationTransaction.Migrate(); err != nil {
		repositoryMigrations.logger.Standard.Error().Msgf("Could not migrate: %v", err)
	}
}

func (repositoryMigrations RepositoryMigrations) RollbackToId(client *gorm.DB, migration interface{}, id string) {
	migrations := GetMigrationsFromType(migration)

	migrationTransaction := gormigrate.New(client, gormigrate.DefaultOptions, migrations)

	if err := migrationTransaction.RollbackTo(id); err != nil {
		repositoryMigrations.logger.Standard.Error().Msgf("Could not rollback to ID: %v", err)
	}
}

func NewRepositoryMigrations(logger *logging.Logger) IMigration {
	repositoryMigrations := RepositoryMigrations{}
	repositoryMigrations.logger = logger
	return &repositoryMigrations
}
