package errors

type Code string

const (
	CodeUnknown       Code = "UNKNOWN"
	CodeValidation    Code = "VALIDATION"
	CodeDomainFailure Code = "DOMAIN_FAILURE"
	CodeInternal      Code = "INTERNAL"
	CodeIO            Code = "IO"
)
