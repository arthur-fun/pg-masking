package main

import (
	"math/rand"
	"time"
)

type replaceAllConverter struct {
	placeHolder string
}

func (rac replaceAllConverter) mask(src interface{}) interface{} {
	if _, ok := src.(string); ok {
		return rac.placeHolder
	}
	return src
}

type replaceConverter struct {
	placeHolder string
	start       int
	maskLength  int
}

func (rc replaceConverter) mask(src interface{}) interface{} {
	if v, ok := src.(string); ok {
		srcLength := len(v)
		s := (rc.start + srcLength) % srcLength
		if s < 0 || s >= srcLength {
			return v
		}
		if s+rc.maskLength >= srcLength {
			return v[0:s] + rc.placeHolder
		}
		return v[0:s] + rc.placeHolder + v[s+rc.maskLength:]
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
