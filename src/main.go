package main

import (
	"telvina/APG2_SmartCalc/internal/view"
	"telvina/APG2_SmartCalc/pkg/configurator"
)

func main() {
	v := view.New(configurator.New())

	v.Run()
}
