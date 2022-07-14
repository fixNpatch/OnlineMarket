package Controllers

import "fmt"

type Controller interface {
	InitController()
}

type MainController struct {}

func (c *MainController) InitController()  {
	fmt.Println("main controller initialized")
}