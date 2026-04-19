// Package routing provides port-to-service routing hints by mapping
// observed port changes to known network service patterns.
package routing

import "github.com/user/portwatch/internal/state"

// Direction indicates whether a port is inbound or outbound.
type Direction string

const (
	Inbound  Direction = "inbound"
	Outbound Direction = "outbound"
)

// Route describes a routing hint for a port.
type Route struct {
	Port      int
	Direction Direction
	Protocol  string
	Note      string
}

// Router maps ports to routing hints.
type Router struct {
	overrides map[int]Route
}

// New creates a Router with optional user-defined overrides.
func New(overrides map[int]Route) *Router {
	if overrides == nil {
		overrides = make(map[int]Route)
	}
	return &Router{overrides: overrides}
}

// Lookup returns a Route for the given port, consulting overrides first,
// then falling back to built-in well-known mappings.
func (r *Router) Lookup(port int) Route {
	if rt, ok := r.overrides[port]; ok {
		return rt
	}
	if rt, ok := builtins[port]; ok {
		return rt
	}
	return Route{Port: port, Direction: Inbound, Protocol: "tcp", Note: ""}
}

// Annotate returns a map of port -> Route for all ports in the given slice.
func (r *Router) Annotate(ports []state.Port) map[int]Route {
	out := make(map[int]Route, len(ports))
	for _, p := range ports {
		out[p.Number] = r.Lookup(p.Number)
	}
	return out
}

// builtins contains well-known port routing hints.
var builtins = map[int]Route{
	22:   {Port: 22, Direction: Inbound, Protocol: "tcp", Note: "SSH"},
	80:   {Port: 80, Direction: Inbound, Protocol: "tcp", Note: "HTTP"},
	443:  {Port: 443, Direction: Inbound, Protocol: "tcp", Note: "HTTPS"},
	3306: {Port: 3306, Direction: Inbound, Protocol: "tcp", Note: "MySQL"},
	5432: {Port: 5432, Direction: Inbound, Protocol: "tcp", Note: "PostgreSQL"},
	6379: {Port: 6379, Direction: Inbound, Protocol: "tcp", Note: "Redis"},
	8080: {Port: 8080, Direction: Inbound, Protocol: "tcp", Note: "HTTP-alt"},
}
