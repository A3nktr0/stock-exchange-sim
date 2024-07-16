package models

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// type TabProcess struct {
// tab map[string]Process
// }

type Process struct {
	ProcessName string
	Craft       []Ingredient
	Produit     []Ingredient
	Cycle       int
}

type Ingredient struct {
	Drawer   *Drawer
	Quantity int
}

type TabProcess map[string]Process

func (tabP TabProcess) AddProcess(process Process) bool {
	_, ok := tabP[process.ProcessName]
	if ok {
		return false
	}
	tabP[process.ProcessName] = process
	return true
}

func (factory Factory) NewProcess(line string) ErrorCustom {
	var process Process
	err := ErrorCustom{Message: "Bad Format Process"}

	data := strings.Split(line, ":")
	if len(data) < 3 {
		return err
	}
	// process.ProcessName = strings.ReplaceAll(data[0], "do_", "")
	process.ProcessName = data[0]
	if process.ProcessName == "" {
		fmt.Printf("Error while parsing `%v`\nExiting...\n", line)
		os.Exit(0)
	}
	tempString := rgx.FindAllString(line, 2)
	if len(tempString) != 2 {
		return err
	}
	listCraft := strings.Split(tempString[0][1:len(tempString[0])-1], ";")
	productionList := tempString[1][1 : len(tempString[1])-1]

	// fmt.Println(listCraft, produit)

	for _, craft := range listCraft {
		dataCraft := strings.Split(craft, ":")
		quantity, errAtoi := strconv.Atoi(dataCraft[1])
		if errAtoi != nil {
			return err
		}
		ress := Ressource{
			Name: dataCraft[0],
		}
		factory.Stock.NewRessource(ress, 0)
		drawer := factory.Stock[ress.Name]
		process.Craft = append(process.Craft, Ingredient{
			Drawer:   &drawer,
			Quantity: quantity,
		})
	}
	tabProduit := strings.Split(productionList, ";")
	for _, dataProduit := range tabProduit {
		produit := strings.Split(dataProduit, ":")
		quantity, errAtoi := strconv.Atoi(produit[1])
		if errAtoi != nil {
			return err
		}
		ress := Ressource{
			Name:      produit[0],
			CraftName: process.ProcessName,
		}
		factory.Stock.NewRessource(ress, 0)
		drawer := factory.Stock[ress.Name]

		process.Produit = append(process.Produit, Ingredient{
			Drawer:   &drawer,
			Quantity: quantity,
		})
	}
	cycle, errcycle := strconv.Atoi(strings.ReplaceAll(data[len(data)-1], ")", ""))
	if errcycle != nil {
		return err
	}
	process.Cycle = cycle
	factory.TabProcess.AddProcess(process)
	return ErrorCustom{}
}

func (tabP TabProcess) Display() {
	for _, proc := range tabP {
		fmt.Println("===================================")
		fmt.Printf("Process Name : %s\n", proc.ProcessName)
		fmt.Println("	=============")
		fmt.Println("	Craft:")
		fmt.Println("	=============")
		for index, ing := range proc.Craft {
			fmt.Printf("	index : %d\n", index)
			fmt.Println("	ressource name", ing.Drawer.Ressource.Name)
			fmt.Println("	Quantity", ing.Quantity)
			fmt.Println("")
		}
		fmt.Println("	=============")
		fmt.Println("	Production:")
		fmt.Println("	=============")
		for index, ing := range proc.Produit {
			fmt.Printf("	index : %d\n", index)
			fmt.Println("	ressource name", ing.Drawer.Ressource.Name)
			fmt.Println("	Quantity", ing.Quantity)
			fmt.Println("")
		}
		fmt.Println("	=============")
		fmt.Println("	Cycle : ", proc.Cycle)
		fmt.Print("\n\n")
	}
}

func (f *Factory) InitProcessLine() {
	Create_log_file()
	check := true
	cycle := 0
	tour := 0
	time := int(37846 * f.Optimise.WaitingTime)
	Inser_log("Main Processes :\n")
	for check && tour < time {
		tour++
		check, cycle = f.Craft(f.Optimise.Production[0], 1, cycle)
	}
	Inser_log(fmt.Sprintf("No more process doable at cycle %v\n", cycle+1))
	Inser_log("Stock :\n")
	for _, drawer := range f.Stock {
		Inser_log(fmt.Sprintf(" %v => %v\n", drawer.Ressource.Name, drawer.Quantity))
	}
}

func (f *Factory) Craft(prodName string, quantityRequire int, cycle int) (bool, int) {
	drawerProd := f.Stock[prodName]
	if quantityRequire <= drawerProd.Quantity && drawerProd.Ressource.Name != f.Optimise.Production[0] {
		return true, cycle
	} else if drawerProd.Ressource.CraftName != "None" {
		process := f.TabProcess[drawerProd.Ressource.CraftName]
		quantityProduit := 0
		for _, ingredient := range process.Produit {
			if ingredient.Drawer.Ressource.Name == prodName {
				quantityProduit = ingredient.Quantity
			}
		}
		MaxCycle := 0
		for _, ingredient := range process.Craft {
			crafting, tempCycle := f.Craft(ingredient.Drawer.Ressource.Name, ingredient.Quantity, cycle)
			if !crafting {
				return crafting, cycle
			}
			if MaxCycle < tempCycle {
				MaxCycle = tempCycle
			}
		}
		cycle = MaxCycle
		for ; 0 < quantityRequire; quantityRequire -= quantityProduit {
			process.Run(f, cycle)
		}
		return true, cycle + process.Cycle
		// return false
	}
	return false, cycle
}

func (process Process) Run(f *Factory, cycle int) {
	Inser_log(fmt.Sprintf(" %v:%v\n", cycle, process.ProcessName))
	for _, ing := range process.Craft {
		current := f.Stock[ing.Drawer.Ressource.Name]
		current.Quantity -= ing.Quantity
		f.Stock[ing.Drawer.Ressource.Name] = current
	}
	for _, ing := range process.Produit {
		current := f.Stock[ing.Drawer.Ressource.Name]
		current.Quantity += ing.Quantity
		f.Stock[ing.Drawer.Ressource.Name] = current
	}
}

// func (f Factory) Cycle(processLine ProcessLine) {
// 	check := true
// 	for check {

// 	}
// }
