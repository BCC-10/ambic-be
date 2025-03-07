package mysql

import (
	"errors"
	"github.com/go-sql-driver/mysql"
	gormMysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

func New(dsn string) (*gorm.DB, error) {
	return gorm.Open(gormMysql.Open(dsn), &gorm.Config{})
}

const (
	ErrDuplicateEntry      = 1062 // Duplicate entry for key
	ErrForeignKeyViolation = 1452 // Cannot add or update a child row: a foreign key constraint fails
	ErrDataTooLong         = 1406 // Data too long for column
	ErrLockWaitTimeout     = 1205 // Lock wait timeout exceeded
	ErrDeadlock            = 1213 // Deadlock found when trying to get lock
	ErrUnknownColumn       = 1054 // Unknown column in field list
	ErrTableNotExists      = 1146 // Table doesn't exist
)

func CheckError(err error, target interface{}) bool {
	if err == nil {
		return false
	}

	switch t := target.(type) {
	case int:
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == uint16(t) {
			log.Printf("MySQL Error [%d]: %s\n", mysqlErr.Number, mysqlErr.Message)
			return true
		}

	case error:
		if errors.Is(err, t) {
			log.Println("GORM Error:", err.Error())
			return true
		}
	}

	return false
}
