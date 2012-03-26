package beego

import (
        "regexp"
        "strings"
)

/*
        Route
*/

// Represents a single route mapping
type Route struct {
        pattern       string
        parameterKeys ParameterKeyMap
        extension     string
        Path          string
        Controller    Controller
        MatcherFuncs  []RouteMatcherFunc
}

func (r *Route) String() string {
        return "{Route:'" + r.Path + "'}"
}

// Makes a new route from the given path
func makeRouteFromPath(path string) *Route {

        // get the path segments
        segments := getPathSegments(path)
        regexSegments := make([]string, len(segments))

        // prepare the parameter key map
        var paramKeys ParameterKeyMap = make(ParameterKeyMap)

        var extension string

        // pull out any dynamic segments
        for index, _ := range segments {

                if isDynamicSegment(segments[index]) {

                        // e.g. {id}

                        paramKeys[strings.Trim(segments[index], "{}")] = index
                        regexSegments[index] = ROUTE_REGEX_PLACEHOLDER

                } else if isExtensionSegment(segments[index]) {

                        // e.g. .json
                        extension = segments[index]

                        // trim off the last space (we don't need it)
                        regexSegments = regexSegments[0 : len(regexSegments)-1]

                } else {

                        // e.g. "groups"

                        regexSegments[index] = segments[index]

                }

        }

        patternString := "/" + strings.Join(regexSegments, "/")

        // return a new route
        var route *Route = new(Route)
        route.pattern = patternString
        route.extension = extension
        route.parameterKeys = paramKeys
        route.Path = path
        route.Controller = nil
        return route

}

// Gets the parameter values for the route from the specified path
func (route *Route) getParameterValueMap(path string) ParameterValueMap {
        return getParameterValueMap(route.parameterKeys, path)
}

// Checks whether a path matches a route or not
func (route *Route) DoesMatchPath(path string) bool {

        match, error := regexp.MatchString(route.pattern, path)

        if error == nil {
                if match {

                        if len(route.extension) > 0 {

                                // make sure the extensions match too
                                return strings.HasSuffix(strings.ToLower(path), strings.ToLower(route.extension))

                        } else {
                                return match
                        }

                } else {

                        return false

                }

        }

        // error :-(
        return false

}

// Checks whether the context for this request matches the route
func (route *Route) DoesMatchContext(c *Context) bool {

        // by default, we match
        var match bool = true

        if len(route.MatcherFuncs) > 0 {

                // there are some matcher functions, so don't automatically
                // match by default - let the matchers decide
                match = false

                // loop through the matcher functions
                for _, f := range route.MatcherFuncs {

                        // modify 'match' based on the result of the matcher function
                        switch f(c) {
                        case NoMatch:
                                match = false
                        case Match:
                                match = true
                        }

                }

        }

        // return the result
        return match

}