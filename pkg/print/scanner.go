package print

import (
	"bufio"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

type Scanner struct {
}

func NewScanner() *Scanner {
	return &Scanner{}
}

func (s *Scanner) PrintList() (*PrinterList, error) {
	// Exécute la commande WMIC pour obtenir la liste des imprimantes
	cmd := exec.Command("wmic", "printer", "get", "name")

	// Récupère la sortie de la commande
	out, err := cmd.Output()
	if err != nil {
		log.Fatalf("Failed to execute command: %v", err)
	}

	// Utilise un print pour lire la sortie ligne par ligne
	scanner := bufio.NewScanner(strings.NewReader(string(out)))

	// Ignore la première ligne qui contient le nom de la colonne
	if scanner.Scan() {
		// Affiche un message d'entête
		fmt.Println("Available Printers:")
	}

	printerList := &PrinterList{
		List:      []string{},
		IndexMap:  make(map[int]string),
		nextIndex: 1,
	}

	// Parcours des lignes restantes
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			printerList.Add(line)
		}
	}

	// Vérifie les erreurs de print
	if err := scanner.Err(); err != nil {
		log.Fatalf("Scanner error: %v", err)
	}

	return printerList, nil
}
