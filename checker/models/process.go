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

func (f *Factory) Checker() {
	for _, logLine := range f.TabLog {
		if logLine[0] == '#' {
			continue
		}
		fmt.Printf("Evaluating: %v\n", logLine)
		name := strings.Split(logLine, ":")[1]
		process := f.TabProcess[name]
		if !f.TestProcess(process) {
			fmt.Printf("at %v stock insufficient\nExiting...\n", logLine)
			os.Exit(0)
		}
	}
	fmt.Println("Trace completed, no error detected.")
}

func (f *Factory) TestProcess(process Process) bool {
	for _, ing := range process.Craft {
		if ing.Quantity > f.Stock[ing.Drawer.Ressource.Name].Quantity {
			return false
		}
	}
	process.Run(f)
	return true
}

func (process Process) Run(f *Factory) {
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
