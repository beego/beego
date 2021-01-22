package error

// A Code is an unsigned 32-bit error code as defined in the beego spec.
type Code uint32

const (
	// SessionSessionStartError means func SessionStart error in session module.
	SessionSessionStartError Code = 5001001
)

// CodeToStr is a map about Code and Code's message
var CodeToStr = map[Code]string{
	SessionSessionStartError: `"SESSION_MODULE_SESSION_START_ERROR"`,
}