package main

import (
	"miniproject/ds"
	"miniproject/gapi"
)

func main() {

	ds.NewDataSource()
	gapi.RunGRPCServer()

}
