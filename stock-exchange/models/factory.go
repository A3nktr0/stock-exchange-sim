package models

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Factory struct {
	Stock      Stock
	TabProcess TabProcess
	Optimise   Optimise
}

var rgx = regexp.MustCompile(`\((.*?)\)`)

func NewFactory() Factory {
	var factory Factory
	findopt := false
	factory.TabProcess = make(TabProcess)
	factory.Stock = make(Stock)

	if len(os.Args[1:]) != 2 {
		log.Fatalln("Error : go run . <file> <waiting_time>\nExemple: go run . examples/simple 10")
	}
	args := os.Args[1:]
	pathFile := args[0]

	waitingTime, errWaitingTime := strconv.ParseFloat(args[1], 64)
	if errWaitingTime != nil {
		log.Fatalln("Error : waiting time Bad Format")
	}
	factory.Optimise.WaitingTime = waitingTime

	data, errFile := os.ReadFile(pathFile)
	if errFile != nil {
		log.Fatalln("Error : file is not find")
	}

	lineTab := strings.Split(string(data), "\n")
	for _, line := range lineTab {
		if line == "" {
			continue
		}
		if line[0] == '#' {
			continue
		}
		checkProcess := len(rgx.FindAllString(line, 2)) == 2
		checkRessource := len(strings.Split(line, ":")) == 2
		checkOptimise, _ := regexp.Match("optimize:", []byte(line))

		switch true {
		case checkOptimise && !findopt:
			opt := string(rgx.Find([]byte(line)))
			tab := strings.Split(opt[1:len(opt)-1], ";")
			for _, data := range tab {
				if data == "time" {
					factory.Optimise.Time = true
					continue
				}
				factory.Optimise.Production = append(factory.Optimise.Production, data)
			}
			findopt = true
		case checkProcess:
			err := factory.NewProcess(line)
			if err.Message != "" {
				log.Fatalln(err.Message)
			}
		case checkRessource && !checkOptimise:
			err := factory.NewRessource(line)
			if err.Message != "" {
				log.Fatalln(err.Message)
			}
		default:
			fmt.Printf("Error while parsing `%v`\nExiting...\n", line)
			os.Exit(0)
		}
	}
	if len(factory.TabProcess) == 0 {
		fmt.Println("Missing processes\nExiting...")
		os.Exit(0)
	}
	// factory.Display()
	return factory
}

func (factory Factory) Display() {
	fmt.Println("TabProcess : ")
	factory.TabProcess.Display()
	fmt.Println("Stock : ")
	factory.Stock.Display()
	fmt.Println("Optimise : ", factory.Optimise)
}
