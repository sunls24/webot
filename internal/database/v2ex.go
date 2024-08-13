package database

import (
	"webot/internal/entity"
)

func (db *DB) V2exExist(id string) bool {
	return db.Limit(1).Find(&entity.V2ex{}, id).RowsAffected == 1
}
