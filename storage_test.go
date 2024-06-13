package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStore(t *testing.T) {
	s := newStore()
	defer teardown(t, s)

	for i := 0; i < 50; i++ {
		key := fmt.Sprintf("foo_%d",i)
		data := []byte("some jpg bytes")
		if err := s.writeStream(key, bytes.NewReader(data)); err != nil {
			t.Error(err)
		}

		if ok := s.Has(key); !ok {
			t.Errorf("disk doesn't have file")
		}

		r, err := s.Read(key)
		if err != nil {
			t.Error(err)
		}
		b, _ := ioutil.ReadAll(r)
		if string(b) != string(data) {
			t.Error("wrong content in file")
		}

		if err := s.Delete(key); err != nil {
			t.Error(err)
		}
	}

}

func TestPathTransform(t *testing.T) {
	key := "mombestpicture"
	pathName := CASPathTransformFunc(key)
	fmt.Println(pathName)
}



func TestDelete(t *testing.T) {

	s := NewStore(StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	})
	key := "specialpicture"
	data := bytes.NewReader([]byte("some jpg bytes"))
	if err := s.writeStream("specialpicture", data); err != nil {
		t.Error(err)
	}

	err := s.Delete(key)
	assert.Nil(t, err, err)
}

func newStore() *Store {
	return NewStore(StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	})
}

func teardown(t *testing.T, s *Store) {
	s.Clear()
}
