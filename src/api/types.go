// This is the api module which includes service wrappers that interact with service providers
// and return a structured SearchResult that will be consumed by many different clients.
package api
// All shared api types and interfaces go here

// This is a high level interface for every api wrapper. If you make a module
// and it implements a method called Query and returns a SearchResult it can
// be considered a valid ApiWrapper
type ApiWrapper interface {
	Query(string) (SearchResult)
}

// SearchResult is the top-level result shared by all services.
// A service module is not valid if it doesn't return this.
type SearchResult struct {
	Success bool
	Media []Media
	RawData []byte
	Service string
}

// Media is a structure currently representing movies and television shows,
// but this could be extended to describe other things such as music. It also
// has the job of abstracting media represented by a myraid of service providers
// that have different payment models including free, free*, rental, and purchase
// with *free being 'SubscriberOnly' (ex. Netflix)
type Media struct { 
	Title string
	Price float64
	Link string
	PurchaseType string // Human-readable: Free, Rent, Buy, Subscribers Only
	HasHD bool
}
