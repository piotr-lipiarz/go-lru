package main

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPutGet(t *testing.T) {
	lru := Constructor[int, int](3)
	lru.Put(1, 11)

	val, err := lru.Get(1)
	require.Equal(t, 11, val)
	require.NoError(t, err)
}

func TestEvictMechanism(t *testing.T) {
	t.Run("should evict LRU item", func(t *testing.T) {
		lru := Constructor[int, int](3)
		lru.Put(1, 11)
		lru.Put(2, 12)
		lru.Put(3, 13)
		// cache capacity reached - Put bellow should evict LRU item
		lru.Put(4, 14)

		val, err := lru.Get(1)
		require.Equal(t, 0, val)
		require.EqualError(t, err, ErrNotFound.Error())
	})

	t.Run("should retain new item", func(t *testing.T) {
		lru := Constructor[int, int](3)
		lru.Put(1, 11)
		lru.Put(2, 12)
		lru.Put(3, 13)
		// cache capacity reached - Put bellow should add new item
		lru.Put(4, 14)

		val, err := lru.Get(4)
		require.Equal(t, 14, val)
		require.NoError(t, err)
	})

	t.Run("should keep items up to capacity amount", func(t *testing.T) {
		lru := Constructor[int, int](3)
		lru.Put(1, 11)
		lru.Put(2, 12)
		lru.Put(3, 13)
		// cache capacity reached - Put bellow should add new item
		lru.Put(4, 14)

		val, err := lru.Get(2)
		require.NoError(t, err)
		require.Equal(t, 12, val)

		val, err = lru.Get(3)
		require.NoError(t, err)
		require.Equal(t, 13, val)
	})
}
