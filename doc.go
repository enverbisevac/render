// Package render helps manage HTTP request / response payloads.
//
// Every well-designed, robust and maintainable Web Service / REST API also needs
// well-_defined_ request and response payloads. Together with the endpoint handlers,
// the request and response payloads make up the contract between your server and the
// clients calling on it.
//
// This is where `render` comes in - offering a few simple helpers to provide a simple
// pattern for managing payload encoding and decoding.

package render
