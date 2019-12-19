package main

import (
	"math/rand"
	"strconv"
	"time"
)

type replaceAllConverter struct {
	placeHolder string
}

func (rac replaceAllConverter) mask(src interface{}) interface{} {
	if ptr, ok := src.(*interface{}); ok {
		if _, ok := (*ptr).(string); ok {
			var newValue interface{}
			newValue = rac.placeHolder
			return &newValue
		}
	}
	return src
}

type replaceConverter struct {
	placeHolder string
	start       int
}

func (rc replaceConverter) mask(src interface{}) interface{} {
	if ptr, ok := src.(*interface{}); ok {
		if v, ok := (*ptr).(string); ok {
			return rc.replaceString(v)
		} else if v, ok := (*ptr).([]uint8); ok {
			return rc.replaceString(string(v))
		}
	}
	return src
}

func (rc replaceConverter) replaceString(v string) *interface{} {
	var newValue interface{}
	srcLength := len(v)
	s := (rc.start + srcLength) % srcLength
	if s < 0 || s >= srcLength {
		newValue = v
	}
	if s+len(rc.placeHolder) >= srcLength {
		newValue = v[0:s] + rc.placeHolder[0:(srcLength-s)]
	} else {
		newValue = v[0:s] + rc.placeHolder + v[s+len(rc.placeHolder):]
	}

	return &newValue
}

type randomConverter struct{}

func (rc randomConverter) mask(src interface{}) interface{} {
	if ptr, ok := src.(*interface{}); ok {
		if v, ok := (*ptr).(int); ok {
			var newValue interface{}
			r := rand.New(rand.NewSource(time.Now().UnixNano()))
			newValue = int(float64(v) * (1.0 + (float64)(r.Intn(100)-50)/50.0))
			return &newValue
		} else if v, ok := (*ptr).(int8); ok {
			var newValue interface{}
			r := rand.New(rand.NewSource(time.Now().UnixNano()))
			newValue = int(float64(v) * (1.0 + (float64)(r.Intn(100)-50)/50.0))
			return &newValue
		} else if v, ok := (*ptr).(int16); ok {
			var newValue interface{}
			r := rand.New(rand.NewSource(time.Now().UnixNano()))
			newValue = int(float64(v) * (1.0 + (float64)(r.Intn(100)-50)/50.0))
			return &newValue
		} else if v, ok := (*ptr).(int32); ok {
			var newValue interface{}
			r := rand.New(rand.NewSource(time.Now().UnixNano()))
			newValue = int(float64(v) * (1.0 + (float64)(r.Intn(100)-50)/50.0))
			return &newValue
		} else if v, ok := (*ptr).(int64); ok {
			var newValue interface{}
			r := rand.New(rand.NewSource(time.Now().UnixNano()))
			newValue = int(float64(v) * (1.0 + (float64)(r.Intn(100)-50)/50.0))
			return &newValue
		} else if v, ok := (*ptr).(uint); ok {
			var newValue interface{}
			r := rand.New(rand.NewSource(time.Now().UnixNano()))
			newValue = int(float64(v) * (1.0 + (float64)(r.Intn(100)-50)/50.0))
			return &newValue
		} else if v, ok := (*ptr).(uint8); ok {
			var newValue interface{}
			r := rand.New(rand.NewSource(time.Now().UnixNano()))
			newValue = int(float64(v) * (1.0 + (float64)(r.Intn(100)-50)/50.0))
			return &newValue
		} else if v, ok := (*ptr).(uint16); ok {
			var newValue interface{}
			r := rand.New(rand.NewSource(time.Now().UnixNano()))
			newValue = int(float64(v) * (1.0 + (float64)(r.Intn(100)-50)/50.0))
			return &newValue
		} else if v, ok := (*ptr).(uint32); ok {
			var newValue interface{}
			r := rand.New(rand.NewSource(time.Now().UnixNano()))
			newValue = int(float64(v) * (1.0 + (float64)(r.Intn(100)-50)/50.0))
			return &newValue
		} else if v, ok := (*ptr).(uint64); ok {
			var newValue interface{}
			r := rand.New(rand.NewSource(time.Now().UnixNano()))
			newValue = int(float64(v) * (1.0 + (float64)(r.Intn(100)-50)/50.0))
			return &newValue
		} else if v, ok := (*ptr).(float64); ok {
			var newValue interface{}
			r := rand.New(rand.NewSource(time.Now().UnixNano()))
			newValue = v * (1.0 + (float64)(r.Intn(100)-50)/50.0)
			return &newValue
		} else if v, ok := (*ptr).(float32); ok {
			var newValue interface{}
			r := rand.New(rand.NewSource(time.Now().UnixNano()))
			newValue = v * (1.0 + (float32)(r.Intn(100)-50)/50.0)
			return &newValue
		} else if v, ok := (*ptr).([]uint8); ok {
			var newValue interface{}
			fv, err := strconv.ParseFloat(string(v), 64)
			checkErr(err)
			r := rand.New(rand.NewSource(time.Now().UnixNano()))
			newValue = (fv) * (1.0 + (float64)(r.Intn(100)-50)/50.0)
			return &newValue
		}
	}

	return src
}
