package models

import (
	"github.com/panjf2000/ants/v2"
)

var pool *ants.Pool

func NewPool(size int) error {
	p, err := ants.NewPool(size)
	pool = p
	return err
}
