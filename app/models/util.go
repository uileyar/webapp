package models

import (
	"github.com/golang/glog"
	"github.com/satori/go.uuid"
)

func CreateGUID() string {
	u1 := uuid.NewV4()
	glog.Infof("UUID: %s\n", u1)
	return u1.String()
}
