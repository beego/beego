package codes

type Code uint32

const (
	SessionSessionStartError Code = 5001001
)

var CodeToStr = map[Code]string{
	SessionSessionStartError : `"SESSION_MODULE_SESSION_START_ERROR"`,
}
