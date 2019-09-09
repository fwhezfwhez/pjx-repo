package db

import (
	"testing"
)

func TestReTry(t *testing.T) {
	if e:=DB.DB().Ping();e!=nil {
		panic(e)
	}
	select{}
}
