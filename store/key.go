package store

import "encoding/binary"

func PrefixBlockHeader(blockNum uint64, blockHash []byte) []byte {
	key := make([]byte, 0, len(blockHeaderPrefix)+8+len(blockHash))
	key = append(key, blockHeaderPrefix...)
	key = append(key, encodeBlockNumber(blockNum)...)
	key = append(key, blockHash...)
	return key
}

func PrefixBlockBody(blockNum uint64, blockHash []byte) []byte {
	key := make([]byte, 0, len(blockBodyPrefix)+8+len(blockHash))
	key = append(key, blockBodyPrefix...)
	key = append(key, encodeBlockNumber(blockNum)...)
	key = append(key, blockHash...)
	return key
}

func encodeBlockNumber(number uint64) []byte {
	enc := make([]byte, 8)
	binary.BigEndian.PutUint64(enc, number)
	return enc
}
