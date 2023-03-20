package store

import "github.com/aergoio/aergo/internal/common"

// Code
func ReadCode(r Reader, codeHash []byte) []byte {
	return r.Get(keyCode(codeHash))
}

func WriteCode(w Writer, code []byte) {
	codeHash := common.Hasher(code)
	w.Set(keyCode(codeHash), code)
}

func DropCode(w Writer, blockHash []byte) {
	w.Delete(keyBlockNumByHash(blockHash))
}

// contract

// address

// ~~
