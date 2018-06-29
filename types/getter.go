package types

import (
	"fmt"
)

// Getter struct for prefixed stores
type PrefixStoreGetter struct {
	key    StoreKey
	prefix []byte
}

func NewPrefixStoreGetter(key StoreKey, prefix []byte) PrefixStoreGetter {
	return PrefixStoreGetter{key, prefix}
}

// Implements sdk.KVStoreGetter
func (getter PrefixStoreGetter) KVStore(ctx Context) KVStore {
	return ctx.KVStore(getter.key).Prefix(getter.prefix)
}

// TransientStoreKey is used for accessing transient substores.
// Each TransientStoreKeys are related with the underlying StoreKey
// Which is registered as a substore to the multistore
type TransientStoreKey struct {
	key StoreKey
}

func NewTransientStoreKey(key StoreKey) TransientStoreKey {
	return TransientStoreKey{
		key: key,
	}
}

func (key TransientStoreKey) Name() string {
	return "!" + key.key.Name()
}

func (key TransientStoreKey) String() string {
	return fmt.Sprintf("TransientStoreKey{%p, %s}", key.key, key.key.Name())
}
