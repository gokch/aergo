package store

// this fields define database key by types

var (
	ChainDBName       = "chain"
	GenesisKey        = []byte(ChainDBName + ".genesisInfo")
	GenesisBalanceKey = []byte(ChainDBName + ".genesisBalance")
	LatestKey         = []byte(ChainDBName + ".latest")

	hardforkKey   = []byte("hardfork")
	lastHeaderKey = []byte("LastHeader")
	lastBlockKey  = []byte("LastBlock")
)

var (
	blockHeaderPrefix = []byte("h") // headerPrefix + num (uint64 big endian) + hash -> header
	// blockNumberPrefix = Prefix[]byte("n") // blockNumberPrefix + hash -> num (uint64 big endian)
	blockPrefix       = []byte("b") // blockBodyPrefix + num (uint64 big endian) + hash -> block body
	transactionPrefix = []byte("t") // blockBodyPrefix + num (uint64 big endian) + hash -> block body
	receiptsPrefix    = []byte("r")
	codePrefix        = []byte("c") // CodePrefix + code hash -> account code

	blockBodyPrefix     = []byte("b") // blockBodyPrefix + num (uint64 big endian) + hash -> block body
	blockReceiptsPrefix = []byte("r") // blockReceiptsPrefix + num (uint64 big endian) + hash -> block receipts
	trieAccountPrefix   = []byte("A")
	trieStoragePrefix   = []byte("O")
)

var (
	raftIdentityKey              = []byte("r_identity")
	raftStateKey                 = []byte("r_state")
	raftSnapKey                  = []byte("r_snap")
	raftEntryLastIdxKey          = []byte("r_last")
	raftEntryPrefix              = []byte("r_entry.")
	raftEntryInvertPrefix        = []byte("r_inv.")
	raftConfChangeProgressPrefix = []byte("r_ccstatus.")

	enter = []byte("keyword_enterprise")
	pref  = []byte("keyword_name")

	enterpriseAdmins = []byte("ADMINS")

	enterpriseConfPrefix     = []byte("conf\\")
	enterpriseRPCPermissions = []byte("RPCPERMISSIONS")
	enterpriseP2PWhite       = []byte("P2PWHITE")
	enterpriseP2PBlack       = []byte("P2PBLACK")
	enterpriseAccountWhite   = []byte("ACCOUNTWHITE")

	namePrefix         = []byte("name")
	systemProposal     = []byte("proposal")
	systemStaking      = []byte("staking")
	systemStakingTotal = []byte("stakingtotal")
	systemVote         = []byte("vote")
	systemTotal        = []byte("total")
	systemVoteSort     = []byte("sort")
	systemVpr          = []byte("VotingPowerBucket/")
	systemParam        = []byte("param\\")

	dposLibStatusKey = []byte("dpos.LibStatus")

	reOrgMarker = []byte("_reorg_marker_")
)
