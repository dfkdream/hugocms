package main

import "errors"

var (
	errDuplicatedUser = errors.New("duplicated user found")
)

// JSON formatted HTTP error messages
const (
	jsonStatusBadRequest          = `{"code":400,"message":"Bad Request"}`
	jsonStatusForbidden           = `{"code":403,"message":"Forbidden"}`
	jsonStatusNotFound            = `{"code":404,"message":"Not Found"}`
	jsonStatusMethodNotAllowed    = `{"code":405,"message":"Method Not Allowed"}`
	jsonStatusInternalServerError = `{"code":500,"message":"Internal Server Error"}`
)
