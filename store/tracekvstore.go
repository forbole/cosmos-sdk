package store

import (
	"encoding/json"
	"io"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	writeOp operation = "W"
	readOp  operation = "R"
)

type (
	// TraceKVStore implements the KVStore interface with tracing enabled.
	// Operations are traced on each core KVStore call and written to the
	// underlying io.writer.
	//
	// TODO: Should we use a buffered writer and implement Commit on
	// TraceKVStore?
	TraceKVStore struct {
		parent sdk.KVStore
		writer io.Writer
	}

	// operation represents an IO operation
	operation string

	// traceOperation implements a traced KVStore operation
	traceOperation struct {
		Timestamp int       `json:"timestamp,omitempty"`
		Operation operation `json:"operation,omitempty"`
		Key       string    `json:"key,omitempty"`
		Value     string    `json:"value,omitempty"`
	}
)

// NewTraceKVStore returns a reference to a new traceKVStore given a parent
// KVStore implementation and a buffered writer.
func NewTraceKVStore(parent sdk.KVStore, writer io.Writer) *TraceKVStore {
	return &TraceKVStore{parent: parent, writer: writer}
}

// Get implements the KVStore interface. It traces a read operation and
// delegates a Get call to the parent KVStore.
func (tkv *TraceKVStore) Get(key []byte) []byte {
	value := tkv.parent.Get(key)

	tkv.writeOperation(readOp, key, value)
	return value
}

// Set implements the KVStore interface. It traces a write operation and
// delegates the Set call to the parent KVStore.
func (tkv *TraceKVStore) Set(key []byte, value []byte) {
	tkv.writeOperation(writeOp, key, value)
	tkv.parent.Set(key, value)
}

// Delete implements the KVStore interface. It traces a write operation and
// delegates the Delete call to the parent KVStore.
func (tkv *TraceKVStore) Delete(key []byte) {
	tkv.writeOperation(writeOp, key, nil)
	tkv.parent.Delete(key)
}

// Has implements the KVStore interface. It delegates the Has call to the
// parent KVStore.
func (tkv *TraceKVStore) Has(key []byte) bool {
	return tkv.parent.Has(key)
}

// Prefix implements the KVStore interface.
func (tkv *TraceKVStore) Prefix(prefix []byte) KVStore {
	return prefixStore{tkv, prefix}
}

// Iterator implements the KVStore interface. It delegates the Iterator call
// the to the parent KVStore.
func (tkv *TraceKVStore) Iterator(start, end []byte) sdk.Iterator {
	return tkv.iterator(start, end, true)
}

// ReverseIterator implements the KVStore interface. It delegates the
// ReverseIterator call the to the parent KVStore.
func (tkv *TraceKVStore) ReverseIterator(start, end []byte) sdk.Iterator {
	return tkv.iterator(start, end, false)
}

// iterator facilitates iteration over a KVStore. It delegates the necessary
// calls to it's parent KVStore.
func (tkv *TraceKVStore) iterator(start, end []byte, ascending bool) sdk.Iterator {
	if ascending {
		return tkv.parent.Iterator(start, end)
	}

	return tkv.parent.ReverseIterator(start, end)
}

// GetStoreType implements the KVStore interface. It returns the underlying
// KVStore type.
func (tkv *TraceKVStore) GetStoreType() sdk.StoreType {
	return tkv.parent.GetStoreType()
}

// CacheWrap implements the KVStore interface. It panics as a TraceKVStore
// cannot be cache wrapped.
func (tkv *TraceKVStore) CacheWrap() sdk.CacheWrap {
	panic("you cannot CacheWrap a TraceKVStore")
}

// writeOperation writes a KVStore operation to the underlying io.Writer as
// JSON-encoded data.
func (tkv *TraceKVStore) writeOperation(op operation, key, value []byte) {
	raw, err := json.Marshal(traceOperation{
		Timestamp: time.Now().UTC().Nanosecond(),
		Operation: op,
		Key:       string(key),
		Value:     string(value),
	})

	if err == nil {
		tkv.writer.Write(raw)
		io.WriteString(tkv.writer, "\n")
	}
}
