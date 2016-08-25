package server

import (
	"net/url"
	"strconv"
)

type Params struct {
	url.Values
}

func (p Params) GetInt64(key string) (int64, error) {
	val := p.Get(key)
	return strconv.ParseInt(val, 10, 64)
}

func (p Params) GetBool(key string) (bool, error) {
	val := p.Get(key)
	return strconv.ParseBool(val)
}

func (p Params) GetFloat64(key string) (float64, error) {
	val := p.Get(key)
	return strconv.ParseFloat(val, 64)
}