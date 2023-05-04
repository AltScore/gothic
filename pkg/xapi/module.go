package xapi

// Module is an interface that should be implemented by all modules provides API endpoints (Route).
type Module interface {
	Routes() []Route
}
