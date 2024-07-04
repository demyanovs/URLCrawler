package utils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

var storeData = map[string]any{
	"key1": 1,
	"key2": 2,
	"key3": 3,
	"key4": PageData{
		URL:        "path1",
		StatusCode: 200,
		Title:      "title1",
		Desc:       "desc1",
		Keywords:   "key1, key2",
	},
}

func TestGet_NoSuchKeyError(t *testing.T) {
	store := NewStore()
	_, err := store.Get("")
	require.ErrorIs(t, ErrorNoSuchKey, err)
}

func TestGet_Success(t *testing.T) {
	store := createAndFillStore()
	for tk, tv := range storeData {
		v, err := store.Get(tk)
		require.NoError(t, err)
		require.Equal(t, v, tv, fmt.Sprintf("key %s, value: %#v", tk, tv))
	}
}

func TestDelete_Success(t *testing.T) {
	store := createAndFillStore()

	res, err := store.Get("key2")
	require.NoError(t, err)
	require.Equal(t, 2, res)
	require.Len(t, store.m, 4)

	store.Delete("key2")
	require.Len(t, store.m, 3)
	_, err = store.Get("key2")
	require.ErrorIs(t, ErrorNoSuchKey, err)

	store.Delete("unknown key")
	require.Len(t, store.m, 3)
}

func TestList_EmptyStoreSuccess(t *testing.T) {
	store := NewStore()
	list := store.List()
	require.Empty(t, list)
}

func TestList_Success(t *testing.T) {
	store := createAndFillStore()
	list := store.List()
	require.Equal(t, storeData, list)
}

func TestKeys_EmptyStoreSuccess(t *testing.T) {
	store := NewStore()
	list := store.Keys()
	require.Empty(t, list)
}

func TestKeys_Success(t *testing.T) {
	store := createAndFillStore()
	list := store.Keys()
	require.ElementsMatch(t, []string{"key1", "key2", "key3", "key4"}, list)
}

func TestValues_Success(t *testing.T) {
	store := createAndFillStore()
	list := store.Values()
	require.ElementsMatch(t, []any{1, 2, 3, PageData{
		URL:        "path1",
		StatusCode: 200,
		Title:      "title1",
		Desc:       "desc1",
		Keywords:   "key1, key2",
	}}, list)
}

func TestClear_Success(t *testing.T) {
	store := createAndFillStore()
	require.Equal(t, 4, store.Len())

	store.Clear()
	require.Equal(t, 0, store.Len())
}

func TestLen_Success(t *testing.T) {
	store := createAndFillStore()
	require.Equal(t, 4, store.Len())

	store.Add("key5", 5)
	require.Equal(t, 5, store.Len())

	store.Delete("key1")
	require.Equal(t, 4, store.Len())
}

func createAndFillStore() *Store {
	store := NewStore()

	for k, v := range storeData {
		store.Add(k, v)
	}

	return store
}
