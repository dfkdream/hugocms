package internal

// JSON formatted HTTP error messages
const (
	JsonStatusBadRequest          = `{"code":400,"message":"Bad Request"}`
	JsonStatusForbidden           = `{"code":403,"message":"Forbidden"}`
	JsonStatusNotFound            = `{"code":404,"message":"Not Found"}`
	JsonStatusMethodNotAllowed    = `{"code":405,"message":"Method Not Allowed"}`
	JsonStatusConflict            = `{"code":409,"message":"Conflict"}`
	JsonStatusInternalServerError = `{"code":500,"message":"Internal Server Error"}`
)
