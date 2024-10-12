package main

import (
	"rest_api_server/cmd"

	"github.com/labstack/echo/v4"
)

func main() {
	WebServer := echo.New()

	v1 := WebServer.Group("/api/v1")
	ui := WebServer.Group("/ui")
	use := WebServer.Group("/user")

	cmd.RoutV1(v1)
	cmd.RoutUI(ui)
	cmd.RoutUser(use)

	WebServer.Logger.Fatal(WebServer.Start(":8081"))
}

// http://localhost:8080/api/v1/items?id=1
// http://localhost:8080/ui/
