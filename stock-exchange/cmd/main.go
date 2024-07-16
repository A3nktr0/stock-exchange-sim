package main

import models "stock-exchange-sim/models"

func main() {
	factory := models.NewFactory()
	factory.InitProcessLine()
}
