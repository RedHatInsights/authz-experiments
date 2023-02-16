package main

import (
	"seatdemo/handler"
)

//test for quay trigger
func main() {
	// Echo instance
	e := handler.GetEcho()

	// Start server
	e.Logger.Fatal(e.Start(":3000"))
}
