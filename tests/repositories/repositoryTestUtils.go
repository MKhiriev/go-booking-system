package repositories

import (
	"database/sql/driver"
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"reflect"
	"testing"
	"time"
)

func setupTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}

	dialector := postgres.New(postgres.Config{
		Conn:                 db,
		PreferSimpleProtocol: true,
	})

	gormDB, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("failed to open gorm db: %v", err)
	}

	cleanup := func() {
		db.Close()
	}

	return gormDB, mock, cleanup
}

type TimeArgument struct {
	checkIfNull bool
}

func AnyTimeArg() TimeArgument {
	return TimeArgument{}
}

func NotNullTimeArg() TimeArgument {
	return TimeArgument{checkIfNull: true}
}

func (a TimeArgument) Match(value driver.Value) bool {
	typeOfValueIsTime := a.typeOfTime(value)
	if a.checkIfNull == true && typeOfValueIsTime && a.isTimeNull(value) {
		return true
	}
	if typeOfValueIsTime {
		return true
	}

	return false
}

func (a TimeArgument) typeOfTime(value driver.Value) bool {
	if reflect.TypeOf(value) == reflect.TypeOf(time.Time{}) {
		return true
	}

	return false
}

func (a TimeArgument) isTimeNull(value driver.Value) bool {
	timeValue := value.(time.Time)
	if !timeValue.IsZero() {
		return true
	}

	return false
}
