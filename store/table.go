package store

import "github.com/aergoio/aergo-lib/db"

// Table is a wrapper around a database that prees each []byte access with a pre-
// configured string.
type Table struct {
	db     db.DB
	prefix []byte
}

var _ db.DB = (*Table)(nil)

func NewTable(db db.DB, prefix []byte) *Table {
	return &Table{
		db:     db,
		prefix: prefix,
	}
}

func (t *Table) Type() string {
	return "table"
}

// Get retrieves the given preed []byte if it's present in the database.
func (t *Table) Set(key, value []byte) {
	t.db.Set(append([]byte(t.prefix), key...), value)
}

// Delete removes the given preed []byte from the database.
func (t *Table) Delete(key []byte) {
	t.db.Delete(append([]byte(t.prefix), key...))
}

// Get retrieves the given preed []byte if it's present in the database.
func (t *Table) Get(key []byte) []byte {
	return t.db.Get(append([]byte(t.prefix), key...))
}

// Has retrieves if a preed version of a []byte is present in the database.
func (t *Table) Exist(key []byte) bool {
	return t.db.Exist(append([]byte(t.prefix), key...))
}

// NewIterator creates a binary-alphabetical iterator over the entire preed []bytespace contained within the database.
func (t *Table) Iterator(start, end []byte) db.Iterator {
	return &TableIterator{
		iterator: t.db.Iterator(append([]byte(t.prefix), start...), append([]byte(t.prefix), end...)),
		prefix:   t.prefix,
	}
}

// NewTx creates a write-only transaction with the given database as its read source.
func (t *Table) NewTx() db.Transaction {
	return &TableTransaction{
		transaction: t.db.NewTx(),
		prefix:      t.prefix,
	}
}

// NewBulk creates a write-only bulk with the given database as its read source.
func (t *Table) NewBulk() db.Bulk {
	return &TableBulk{
		bulk:   t.db.NewBulk(),
		prefix: t.prefix,
	}
}

// Close is a noop to implement the Database interface.
func (t *Table) Close() {
	return
}

// TableIterator is a wrapper around a database iterator that prees each []byte access with a pre-
// configured string.
type TableIterator struct {
	iterator db.Iterator
	prefix   []byte
}

var _ db.Iterator = (*TableIterator)(nil)

func NewTableIterator(iterator db.Iterator, prefix []byte) *TableIterator {
	return &TableIterator{
		iterator: iterator,
		prefix:   prefix,
	}
}

func (t *TableIterator) Next() {
	t.iterator.Next()
}

func (t *TableIterator) Valid() bool {
	return t.iterator.Valid()
}

func (t *TableIterator) Key() []byte {
	return t.iterator.Key()[len(t.prefix):]
}

func (t *TableIterator) Value() []byte {
	return t.iterator.Value()
}

// TableTransaction is a wrapper around a database transaction that prees each []byte access with a pre-
// configured string.
type TableTransaction struct {
	transaction db.Transaction
	prefix      []byte
}

func NewTableTransaction(transaction db.Transaction, prefix []byte) *TableTransaction {
	return &TableTransaction{
		transaction: transaction,
		prefix:      prefix,
	}
}

var _ db.Transaction = (*TableTransaction)(nil)

func (t *TableTransaction) Set(key, value []byte) {
	t.transaction.Set(append([]byte(t.prefix), key...), value)
}

func (t *TableTransaction) Delete(key []byte) {
	t.transaction.Delete(append([]byte(t.prefix), key...))
}

func (t *TableTransaction) Commit() {
	t.transaction.Commit()
}

func (t *TableTransaction) Discard() {
	t.transaction.Discard()
}

// TableBulk is a wrapper around a database bulk that prees each []byte access with a pre-
// configured string.
type TableBulk struct {
	bulk   db.Bulk
	prefix []byte
}

var _ db.Bulk = (*TableBulk)(nil)

func NewTableBulk(bulk db.Bulk, prefix []byte) *TableBulk {
	return &TableBulk{
		bulk:   bulk,
		prefix: prefix,
	}
}

func (t *TableBulk) Set(key, value []byte) {
	t.bulk.Set(append([]byte(t.prefix), key...), value)
}

func (t *TableBulk) Delete(key []byte) {
	t.bulk.Delete(append([]byte(t.prefix), key...))
}

func (t *TableBulk) Flush() {
	t.bulk.Flush()
}

func (t *TableBulk) DiscardLast() {
	t.bulk.DiscardLast()
}
