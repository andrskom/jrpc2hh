package testservice

import (
	jModels "github.com/andrskom/jrpc2hh/models"
	"models"
	anotherModel "some/models"
	"net"
)

// jrpc2hh:service
type Test2 struct {
	db string // Test data
}

func NewTest2(db string) *Test2 {
	return &Test2{db: db}
}

type Test2NilArgsResult struct {
	SomeData string `json:"some_data"`
}

// jrpc2hh:method
func (s *Test2) NilArgs(args jModels.NilArgs, res *Test2NilArgsResult) error {
	return nil
}

type Test2NilResultArgs struct {
	RequiredParam string `json:"required_param"`
	OptionalParam *int   `json:"optional_param"`
}

// jrpc2hh:method
func (s *Test2) NilResult(args models.Test2NilResultArgs, res *jModels.NilResult) error {
	return nil
}

// jrpc2hh:method
func (s *Test2) AnotherPackageResult(args anotherModel.NilArgs, res *models.SomeModel) error {
	return nil
}

// jrpc2hh:method
func (s *Test2) DoubleStarAnotherResult(args jModels.NilArgs, res **models.SomeModel) error {
	return nil
}

// jrpc2hh:method
func (s *Test2) DoubleStarResult(args jModels.NilArgs, res **Test2NilArgsResult) error {
	return nil
}
