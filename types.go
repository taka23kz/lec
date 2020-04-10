package main

import (
	"github.com/go-gorp/gorp"
	ut "github.com/go-playground/universal-translator"
	"gopkg.in/go-playground/validator.v9"
)

// Controller is a controller for this application.
type Controller struct {
	dbmap *gorp.DbMap
}

// Error indicate response erorr
type Error struct {
	Error string `json:"error"`
}

// Validator is implementation of validation of rquest values.
type Validator struct {
	trans     ut.Translator
	validator *validator.Validate
}
