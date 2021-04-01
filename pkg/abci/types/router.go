package types

import (
	"regexp"
)
var (
	// IsAlphaNumeric defines a regular expression for matching against alpha-numeric
	// values.
	IsAlphaNumeric = regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString

	// IsAlphaLower defines regular expression to check if the string has lowercase
	// alphabetic characters only.
	IsAlphaLower = regexp.MustCompile(`^[a-z]+$`).MatchString

	// IsAlphaUpper defines regular expression to check if the string has uppercase
	// alphabetic characters only.
	IsAlphaUpper = regexp.MustCompile(`^[A-Z]+$`).MatchString

	// IsAlpha defines regular expression to check if the string has alphabetic
	// characters only.
	IsAlpha = regexp.MustCompile(`^[a-zA-Z]+$`).MatchString

	// IsNumeric defines regular expression to check if the string has numeric
	// characters only.
	IsNumeric = regexp.MustCompile(`^[0-9]+$`).MatchString
)

// Router provides handlers for each transfer types.
type Router interface {
	AddRoute(r string, h Handler) (rtr Router)
	Route(path string) (h Handler)
}

// map a transfer types to a handler and an initgenesis function
type route struct {
	r string
	h Handler
}

type router struct {
	routes []route
}

// nolint
// NewRouter - create new router
// TODO either make Function unexported or make return types (router) Exported
func NewRouter() *router {
	return &router{
		routes: make([]route, 0),
	}
}


// AddRoute - TODO add description
func (rtr *router) AddRoute(r string, h Handler) Router {
	if !IsAlphaNumeric(r) {
		panic("route expressions can only contain alphanumeric characters")
	}
	rtr.routes = append(rtr.routes, route{r, h})

	return rtr
}

// Route - TODO add description
// TODO handle expressive matches.
func (rtr *router) Route(path string) (h Handler) {
	for _, route := range rtr.routes {
		if route.r == path {
			return route.h
		}
	}
	return nil
}
