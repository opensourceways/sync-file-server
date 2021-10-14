package models

import (
	"github.com/panjf2000/ants/v2"
)

var pool *ants.Pool

func NewPool(size int, log ants.Logger) error {
	p, err := ants.NewPool(size, ants.WithOptions(ants.Options{
		Logger: log,
	}))
	pool = p
	return err
}

func Stop() {
	if pool != nil {
		pool.Release()
	}
}
