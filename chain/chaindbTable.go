package chain

import "github.com/aergoio/aergo-lib/db"

// table is a wrapper around a database that prefixes each key access with a pre-
// configured string.
type table struct {
	db     db.DB
	prefix string
}

func NewTable(db db.DB, prefix string) *table {
	return &table{
		db:     db,
		prefix: prefix,
	}
}

func (t *table) Type(key []byte) string {
	return "table"
}

// Get retrieves the given prefixed key if it's present in the database.
func (t *table) Set(key, value []byte) {
	t.db.Set(append([]byte(t.prefix), key...), value)
}

// Delete removes the given prefixed key from the database.
func (t *table) Delete(key []byte) {
	t.db.Delete(append([]byte(t.prefix), key...))
}

// Get retrieves the given prefixed key if it's present in the database.
func (t *table) Get(key []byte) []byte {
	return t.db.Get(append([]byte(t.prefix), key...))
}

// Has retrieves if a prefixed version of a key is present in the database.
func (t *table) Exist(key []byte) bool {
	return t.db.Exist(append([]byte(t.prefix), key...))
}

// NewIterator creates a binary-alphabetical iterator over the entire prefixed keyspace contained within the database.
func (t *table) Iterator(start, end []byte) db.Iterator {
	return &tableIterator{
		iterator: t.db.Iterator(append([]byte(t.prefix), start...), append([]byte(t.prefix), end...)),
		prefix:   t.prefix,
	}
}

// NewTx creates a write-only transaction with the given database as its read source.
func (t *table) NewTx() db.Transaction {
	return &tableTransaction{
		transaction: t.db.NewTx(),
		prefix:      t.prefix,
	}
}

// NewBulk creates a write-only bulk with the given database as its read source.
func (t *table) NewBulk() db.Bulk {
	return &tableBulk{
		bulk:   t.db.NewBulk(),
		prefix: t.prefix,
	}
}

// Close is a noop to implement the Database interface.
func (t *table) Close() {
	return
}

type tableIterator struct {
	iterator db.Iterator
	prefix   string
}

func (t tableIterator) Next() {
	t.iterator.Next()
}

func (t tableIterator) Valid() bool {
	return t.iterator.Valid()
}

func (t tableIterator) Key() []byte {
	return t.iterator.Key()[len(t.prefix):]
}

func (t tableIterator) Value() []byte {
	return t.iterator.Value()
}

type tableTransaction struct {
	transaction db.Transaction
	prefix      string
}

func (t *tableTransaction) Set(key, value []byte) {
	t.transaction.Set(append([]byte(t.prefix), key...), value)
}

func (t *tableTransaction) Delete(key []byte) {
	t.transaction.Delete(append([]byte(t.prefix), key...))
}

func (t *tableTransaction) Commit() {
	t.transaction.Commit()
}

func (t *tableTransaction) Discard() {
	t.transaction.Discard()
}

type tableBulk struct {
	bulk   db.Bulk
	prefix string
}

func (t *tableBulk) Set(key, value []byte) {
	t.bulk.Set(append([]byte(t.prefix), key...), value)
}

func (t *tableBulk) Delete(key []byte) {
	t.bulk.Delete(append([]byte(t.prefix), key...))
}

func (t *tableBulk) Flush() {
	t.bulk.Flush()
}

func (t *tableBulk) DiscardLast() {
	t.bulk.DiscardLast()
}
