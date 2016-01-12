package models

import (
	"math"

	"github.com/golang/glog"
	"github.com/satori/go.uuid"
)

func CreateGUID() string {
	u1 := uuid.NewV4()
	glog.Infof("UUID: %s\n", u1)
	return u1.String()
}

func Round(f float32, n int) float32 {
	pow10_n := math.Pow10(n)
	return float32(math.Trunc((float64(f)+0.5/pow10_n)*pow10_n) / pow10_n)
}
