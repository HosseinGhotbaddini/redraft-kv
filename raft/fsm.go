package raft

import (
	"encoding/json"

	"github.com/HosseinGhotbaddini/redraft-kv/store"
)

type FSM interface {
	Apply([]byte) error
}

type FSMImpl struct {
	Store *store.Store
}

func (f *FSMImpl) Apply(data []byte) error {
	var cmd store.Command
	if err := json.Unmarshal(data, &cmd); err != nil {
		return err
	}

	switch cmd.Op {
	case "set":
		f.Store.Set(cmd.Key, cmd.Value)
	case "delete":
		f.Store.Delete(cmd.Key)
	default:
		return nil
	}
	return nil
}
