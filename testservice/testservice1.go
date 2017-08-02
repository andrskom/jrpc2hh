package testservice

import (
	jModels "github.com/andrskom/jrpc2hh/models"
	"models"
	"net"
)

// jrpc2hh:service
type Test1 struct {
	db string // Test data
}

func NewTest1(db string) *Test1 {
	return &Test1{db: db}
}

type Test1NilArgsResult struct {
	SomeData string `json:"some_data"`
}

// jrpc2hh:method
func (s *Test1) NilArgs(args jModels.NilArgs, res *Test1NilArgsResult) error {
	return nil
}

type Test1NilResultArgs struct {
	RequiredParam string `json:"required_param"`
	OptionalParam *int   `json:"optional_param"`
}

// jrpc2hh:method
func (s *Test1) NilResult(args Test1NilResultArgs, res *jModels.NilResult) error {
	return nil
}

// jrpc2hh:method:withContext
func (s *Test1) AnotherPackageResult(args jModels.NilArgs, res *models.SomeModel) error {
	return nil
}

// jrpc2hh:method
func (s *Test1) DoubleStarAnotherResult(args jModels.NilArgs, res **models.SomeModel) error {
	return nil
}

// jrpc2hh:method
func (s *Test1) DoubleStarResult(args jModels.NilArgs, res **Test1NilArgsResult) error {
	return nil
}
