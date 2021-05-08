package common

import "errors"

type ErrCfgExists struct {
	Path string
}

func (e *ErrCfgExists) Error() string {
	return "config "+e.Path+ " already exists"
}

type ErrCfgNotExists struct {
	Path string
}

func (e *ErrCfgNotExists) Error() string {
	return "config " + e.Path + " does not exist"
}

type ErrDirectoryNotEmpty struct {
	Path string
}

func (e *ErrDirectoryNotEmpty) Error() string {
	return "directory " + e.Path + " is not empty"
}

type ErrDirectoryNotExists struct {
	Path string
}

func (e *ErrDirectoryNotExists) Error() string {
	return "directory " + e.Path + " not exists"
}

var ErrEmptyRelative = errors.New("empty relative path")

type ErrNoSuchNode struct {
	Relative string
}
func (e *ErrNoSuchNode) Error() string {
	return "no such node: " + e.Relative
}