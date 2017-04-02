package service

import "log"

type ServiceList map[string]*Service

func (sl ServiceList) Add(sName string) {
	if _, ok := sl[sName]; ok {
		log.Fatalf("Service with name '%s' already added", sName)
	} else {
		sl[sName] = NewService(sName)
	}
}

type Service struct {
	Name string
}

func NewService(n string) *Service {
	return &Service{n}
}
