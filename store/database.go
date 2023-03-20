package store

import "github.com/aergoio/aergo-lib/db"

// Writer represents to store data regardless of db and tx
type Writer interface {
	Set(key, value []byte)
	Delete(key []byte)
}

// Reader represents to read data
type Reader interface {
	Exist(key []byte) bool
	Get(key []byte) []byte
}

// Inspector traverses the database and prints out the key/value pairs.
func InspectDatabase(db db.DB, keyPrefix, keyStart []byte) error {

	return nil
}
