package ui

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"win-files-to-print/pkg/print"
)

type CLI struct {
	config *print.ConfigPrinter
}

func NewCLI(config *print.ConfigPrinter) *CLI {
	return &CLI{config: config}
}

func (c *CLI) ChoosePrinter(printerList *print.PrinterList) {
	for {
		printerList.Render()

		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter the index of the item to select: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		index, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("Invalid input. Please enter a valid index.")
			continue
		}

		printName, found := printerList.GetByIndex(index)
		if found == false {
			fmt.Println("Invalid print.")
			continue
		}

		c.config.SetPrintName(printName)

		fmt.Printf("Selected: %s\n", printName)
		return
	}

}

func (c *CLI) ChooseFolder() {
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Veuillez entrer le chemin du dossier : ")
		input, _ := reader.ReadString('\n')
		dirname := strings.TrimSpace(input)

		if CheckDirExists(dirname) {
			c.config.SetFolder(dirname)
			return
		}

		fmt.Printf("Le dossier '%s' n'existe pas ou n'est pas un dossier.\n", dirname)
		continue
	}
}

func CheckDirExists(dirname string) bool {
	info, err := os.Stat(dirname)
	if err != nil {
		return false
	}
	return info.IsDir()
}
