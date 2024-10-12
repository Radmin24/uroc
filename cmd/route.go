package cmd

import "github.com/labstack/echo/v4"

func RoutV1(group *echo.Group) {
	group.GET("/", startV1)
	group.POST("/items", new_item)
	group.GET("/items", get_items)
}

func RoutUI(group *echo.Group) {
	group.GET("/", startUI)

}

func RoutUser(group *echo.Group) {
	group.GET("/", Get_user)
	// group.POST("/init", get_items)
}
