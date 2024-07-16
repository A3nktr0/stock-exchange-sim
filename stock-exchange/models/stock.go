package models

import (
	"fmt"
	"strconv"
	"strings"
)

// type Stock struct {
// tab map[string]Ressource
// }

type Stock map[string]Drawer

type Drawer struct {
	Ressource Ressource
	Quantity  int
}

func (s Stock) NewRessource(ressource Ressource, quantity int) {
	current, ok := s[ressource.Name]
	if ok {
		if current.Ressource.CraftName != ressource.CraftName && current.Ressource.CraftName == "None" && ressource.CraftName != "" {
			current.Ressource.CraftName = ressource.CraftName
			s[ressource.Name] = current
		}
		return
	} else {
		s[ressource.Name] = Drawer{
			Ressource: ressource,
			Quantity:  quantity,
		}
	}
}

func (f Factory) NewRessource(cmd string) ErrorCustom {
	err := ErrorCustom{Message: "Bad Format Ressource"}
	temp := strings.Split(cmd, ":")
	newRessource := Ressource{
		Name:      temp[0],
		CraftName: "None",
	}
	quantity, errAtoi := strconv.Atoi(temp[1])
	if errAtoi != nil {
		return err
	}
	f.Stock.NewRessource(newRessource, quantity)
	return ErrorCustom{}
}

func (s Stock) Display() {
	for _, drawer := range s {
		fmt.Println("	=============")
		fmt.Println("	Ressource name : ", drawer.Ressource.Name)
		fmt.Println("	Craft Name : ", drawer.Ressource.CraftName)
		fmt.Println("	Quantity : ", drawer.Quantity)
	}
	fmt.Println("	=============")
}
