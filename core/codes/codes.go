package codes

type Code uint32

const (
	SessionSessionStartError Code = 5001001
)

var strToCode = map[string]Code{
	`"SESSION_MODULE_SESSION_START_ERROR"`:   SessionSessionStartError,
}
