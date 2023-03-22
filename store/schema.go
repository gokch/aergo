package store

import "github.com/aergoio/aergo/types"

// this fields define database key by types

var (
	ChainDBName       = "chain"
	GenesisKey        = []byte(ChainDBName + ".genesisInfo")
	GenesisBalanceKey = []byte(ChainDBName + ".genesisBalance")
	LatestKey         = []byte(ChainDBName + ".latest")

	hardforkKey   = []byte("hardfork")
	lastHeaderKey = []byte("LastHeader")
	lastBlockKey  = []byte("LastBlock")

	reOrgKey = []byte("_reorg_marker_")
)

var (
	blockHeaderPrefix    = []byte("h") // headerPrefix + num (uint64 big endian) + hash -> header
	blockNumByHashPrefix = []byte("n") // blockNumberPrefix + hash -> num (uint64 big endian)
	blockBodyPrefix      = []byte("b") // blockBodyPrefix + num (uint64 big endian) + hash -> block body
	blockReceiptPrefix   = []byte("r") // blockReceiptsPrefix + num (uint64 big endian) + hash -> block receipts

	txLookupPrefix = []byte("t") // txLookupPrefix + hash -> transaction/receipt lookup metadata

	codePrefix        = []byte("c") // CodePrefix + code hash -> account code
	trieAccountPrefix = []byte("A")
	trieStoragePrefix = []byte("O")
)

func keyBlockHeader(blockNum uint64, blockHash []byte) []byte {
	return append(keyBlockHeaderNumber(blockNum), blockHash...)
}

func keyBlockHeaderNumber(blockNum uint64) []byte {
	return append(blockHeaderPrefix, types.BlockNoToBytes(blockNum)...)
}

func keyBlockNumByHash(blockHash []byte) []byte {
	return append(blockNumByHashPrefix, blockHash...)
}

func keyBlockBody(blockNum uint64, blockHash []byte) []byte {
	return append(keyBlockBodyNumber(blockNum), blockHash...)
}

func keyBlockBodyNumber(blockNum uint64) []byte {
	return append(blockBodyPrefix, types.BlockNoToBytes(blockNum)...)
}

func keyBlockReceipt(blockNum uint64) []byte {
	return append(blockReceiptPrefix, types.BlockNoToBytes(blockNum)...)
}

func keyTxLookup(txhash []byte) []byte {
	return append(txLookupPrefix, txhash...)
}

func keyCode(code []byte) []byte {
	return append(codePrefix, code...)
}

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
)
