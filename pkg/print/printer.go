package print

import (
	"fmt"
	"os"
	"os/exec"
)

const DefaultSumatraPath = "C:\\Users\\Administrateur\\AppData\\Local\\SumatraPDF\\SumatraPDF.exe"
const DefaultAcrobatPath = "c:\\Program Files\\Adobe\\Acrobat DC\\Acrobat\\Acrobat.exe"

type Printer struct {
	printerName string
}

func NewPrinter(printerName string) *Printer {
	return &Printer{printerName: printerName}
}

func (p *Printer) Print(pdfPath string) error {
	/*cmd := exec.Command(DefaultSumatraPath, "-print-to", p.printerName, pdfPath)*/
	cmd := exec.Command(DefaultAcrobatPath, "/t", pdfPath, fmt.Sprintf("\"%s\"", p.printerName))
	// Acrobat.exe /n /s /o /h /t
	fmt.Println("cmd", cmd)

	// Démarrer la commande
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("erreur lors du démarrage de la commande: %v", err)
	}

	// Attendre que la commande se termine
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("erreur lors de l'attente de la fin de la commande: %v", err)
	}

	return nil
}

func (p *Printer) GetPrinterName() string {
	return p.printerName
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
