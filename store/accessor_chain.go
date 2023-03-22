package store

import (
	"github.com/aergoio/aergo/types"
	"github.com/golang/protobuf/proto"
)

// BlockNumByHash
func ReadBlockNumByHash(r Reader, blockHash []byte) uint64 {
	blockNo := r.Get(keyBlockNumByHash(blockHash))
	return types.BytesToUint64(blockNo)
}

func WriteBlockNumByHash(w Writer, blockHash []byte, blockNum uint64) {
	w.Set(keyBlockNumByHash(blockHash), types.Uint64ToBytes(blockNum))
}

func DropBlockNumByHash(w Writer, blockHash []byte) {
	w.Delete(keyBlockNumByHash(blockHash))
}

// Block header
func ReadBlockHeader(r Reader, blockNo uint64, blockHash []byte) *types.BlockHeader {
	raw := r.Get(keyBlockHeader(blockNo, blockHash))
	if raw == nil {
		return nil
	}
	blockHeader := new(types.BlockHeader)
	err := proto.Unmarshal(raw, blockHeader)
	if err != nil {
		return nil
	}

	return blockHeader
}

func WriteBlockHeader(w Writer, blockNo uint64, blockHash []byte, blockHeader *types.BlockHeader) {
	raw, err := proto.Marshal(blockHeader)
	if err != nil {
		return
	}
	w.Set(keyBlockHeader(blockNo, blockHash), raw)
}

func DropBlockHeader(w Writer, blockNo uint64, blockHash []byte) {
	w.Delete(keyBlockHeader(blockNo, blockHash))
}

// Block body ( transactions )
func ReadBlockBody(r Reader, blockNo uint64, blockHash []byte) *types.BlockBody {
	raw := r.Get(keyBlockBody(blockNo, blockHash))
	if raw == nil {
		return nil
	}
	blockBody := new(types.BlockBody)
	err := proto.Unmarshal(raw, blockBody)
	if err != nil {
		return nil
	}
	return blockBody
}

func WriteBlockBody(w Writer, blockNo uint64, blockHash []byte, blockBody *types.BlockBody) {
	raw, err := proto.Marshal(blockBody)
	if err != nil {
		return
	}
	w.Set(keyBlockBody(blockNo, blockHash), raw)
}

func DropBlockBody(w Writer, blockNo uint64, blockHash []byte) {
	w.Delete(keyBlockBody(blockNo, blockHash))
}

// Block Receipt ( TODO )
/*
func ReadBlockReceipt(r Reader, blockNo uint64, blockHash []byte) *types.Receipt {
	raw := r.Get(keyBlockBody(blockNo, blockHash))
	if raw == nil {
		return nil
	}
	receipt := new(types.Receipt)
	err := proto.Unmarshal(raw, receipt)
	if err != nil {
		return nil
	}
	return blockBody
}

func WriteBlockBody(w Writer, blockNo uint64, blockHash []byte, blockBody *types.BlockBody) {
	raw, err := proto.Marshal(blockBody)
	if err != nil {
		return
	}
	w.Set(keyBlockBody(blockNo, blockHash), raw)
}

func DropBlockBody(w Writer, blockNo uint64, blockHash []byte) {
	w.Delete(keyBlockBody(blockNo, blockHash))
}
*/

// TxLookup
func ReadTxLookup(r Reader, txHash []byte) (blockNo uint64) {
	raw := r.Get(keyTxLookup(txHash))
	return types.BytesToUint64(raw)
}

func WriteTxLookup(w Writer, txHash []byte, blockNo uint64) {
	w.Set(keyTxLookup(txHash), types.Uint64ToBytes(blockNo))
}

func DropTxLookup(w Writer, txHash []byte) {
	w.Delete(keyTxLookup(txHash))
}
