package models

import "strings"

type AppInfo struct {
	APP      string
	PROJ     string
	ENV      string
	RUNNTIME string
	TIME     string
	ID       string
	DEBPACK  string
	DEBUG    bool
}

func (a *AppInfo) DebSplit(debpkg string) string {
	return strings.Join(strings.Split(debpkg, ","), " ")
}