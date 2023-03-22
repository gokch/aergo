package store

import "github.com/aergoio/aergo/types"

func ReadGenesis(r Reader) []byte {
	return r.Get(GenesisKey)
}

func WriteGenesis(w Writer, genesis []byte) {
	w.Set(GenesisKey, genesis)
}

func ReadLatest(r Reader) uint64 {
	return types.BytesToUint64(r.Get(LatestKey))
}

func WriteLatest(w Writer, BlockNo uint64) {
	w.Set(LatestKey, types.Uint64ToBytes(BlockNo))
}

func ReadHardfork(r Reader) []byte {
	return r.Get(hardforkKey)
}

func WriteHardfork(w Writer, hardfork []byte) {
	w.Set(hardforkKey, hardfork)
}

func ReadReorgMarker(r Reader) []byte {
	return r.Get(reOrgKey)
}

func WriteReorgMarker(w Writer, reorg []byte) {
	w.Set(reOrgKey, reorg)
}

func DropReorgMarker(w Writer) {
	w.Delete(reOrgKey)
}

// config
