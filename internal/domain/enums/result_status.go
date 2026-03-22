package enums

type ResultStatus string

const (
	ResultSuccess        ResultStatus = "SUCCESS"
	ResultPartialSuccess ResultStatus = "PARTIAL_SUCCESS"
	ResultFailed         ResultStatus = "FAILED"
)
