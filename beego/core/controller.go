package beego

// Interface for controller types that handle requests
type Controller interface {

        // When implemented, handles the request
        HandleRequest(c *Context)
}

// The ControllerFunc type is an adapter to allow the use of
// ordinary functions as goweb handlers.  If f is a function
// with the appropriate signature, ControllerFunc(f) is a
// Controller object that calls f.
type ControllerFunc func(*Context)

// HandleRequest calls f(c).
func (f ControllerFunc) HandleRequest(c *Context) {
        f(c)
}