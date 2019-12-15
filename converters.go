package main

import (
	"math/rand"
	"time"
)

type replaceConverter struct {
	placeHolder string
}

func (rc replaceConverter) mask(src interface{}) interface{} {
	if _, ok := src.(string); ok {
		return rc.placeHolder
	}
	return src
}

type randomConverter struct{}

func (rc randomConverter) mask(src interface{}) interface{} {
	if v, ok := src.(*int); ok {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		return int(float64(*v) * (1.0 + (float64)(r.Intn(100)-50)/50.0))
	} else if v, ok := src.(*float64); ok {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		return (*v) * (1.0 + (float64)(r.Intn(100)-50)/50.0)
	}
	return src
}
