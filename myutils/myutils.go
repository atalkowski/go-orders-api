package myutils

import (
	"errors"
	"net/http"
	"strconv"
)

var ErrNotExist = errors.New("item does not exist")
var ErrClientIsNull = errors.New("client is null")

func StrToUInt(s string) (uint64, error) {
	const decimal = 10
	const bitSize = 64
	return strconv.ParseUint(s, decimal, bitSize)
}

func GetUIntParam(r *http.Request, name string, dft uint64) (uint64, error) {
	res := r.URL.Query().Get(name)
	if res == "" {
		return dft, nil
	}
	const decimal = 10
	const bitSize = 64
	return strconv.ParseUint(res, decimal, bitSize)
}

func GetQueryParam(r *http.Request, name string, dft string) string {
	res := r.URL.Query().Get(name)
	if res == "" {
		return dft
	}
	return res
}
