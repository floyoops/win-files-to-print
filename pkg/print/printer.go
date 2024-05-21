package print

import (
	"fmt"
	"os"
	"os/exec"
)

const DefaultSumatraPath = "C:\\Users\\Administrateur\\AppData\\Local\\SumatraPDF\\SumatraPDF.exe"

type Printer struct {
	printerName string
}

func NewPrinter(printerName string) *Printer {
	return &Printer{printerName: printerName}
}

func (p *Printer) Print(pdfPath string) error {
	cmd := exec.Command(DefaultSumatraPath, "-print-to", p.printerName, pdfPath)

	fmt.Println("cmd", cmd)

	if err := cmd.Run(); err != nil {
		fmt.Println(err)
		return fmt.Errorf("failed to execute print command: %v", err)
	}

	return nil
}

func CheckLibPdfExist() bool {
	_, err := os.Stat(DefaultSumatraPath)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}
