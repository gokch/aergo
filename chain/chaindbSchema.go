package chain

const (
	chainDBName       = "chain"
	genesisKey        = chainDBName + ".genesisInfo"
	genesisBalanceKey = chainDBName + ".genesisBalance"
)

// this fields define database prefix by types
type prefix []byte

var (
	latestKey   = []byte(chainDBName + ".latest")
	hardforkKey = []byte("hardfork")

	blockPrefix       = []byte("b") // blockBodyPrefix + num (uint64 big endian) + hash -> block body
	transactionPrefix = []byte("t") // blockBodyPrefix + num (uint64 big endian) + hash -> block body
	receiptsPrefix    = []byte("r")

	blockBodyPrefix     = []byte("b") // blockBodyPrefix + num (uint64 big endian) + hash -> block body
	blockReceiptsPrefix = []byte("r") // blockReceiptsPrefix + num (uint64 big endian) + hash -> block receipts
	codePrefix          = []byte("c") // CodePrefix + code hash -> account code
	trieAccountPrefix   = []byte("A")
	trieStoragePrefix   = []byte("O")

	raftIdentityKey              = []byte("r_identity")
	raftStateKey                 = []byte("r_state")
	raftSnapKey                  = []byte("r_snap")
	raftEntryLastIdxKey          = []byte("r_last")
	raftEntryPrefix              = []byte("r_entry.")
	raftEntryInvertPrefix        = []byte("r_inv.")
	raftConfChangeProgressPrefix = []byte("r_ccstatus.")
)
