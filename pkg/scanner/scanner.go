package scanner

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

func (s *Scanner) PrintList() {
	// Exécute la commande WMIC pour obtenir la liste des imprimantes
	cmd := exec.Command("wmic", "printer", "get", "name")

	// Récupère la sortie de la commande
	out, err := cmd.Output()
	if err != nil {
		log.Fatalf("Failed to execute command: %v", err)
	}

	// Utilise un scanner pour lire la sortie ligne par ligne
	scanner := bufio.NewScanner(strings.NewReader(string(out)))

	// Ignore la première ligne qui contient le nom de la colonne
	if scanner.Scan() {
		// Affiche un message d'entête
		fmt.Println("Available Printers:")
	}

	// Parcours des lignes restantes
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			fmt.Println(line)
		}
	}

	// Vérifie les erreurs de scanner
	if err := scanner.Err(); err != nil {
		log.Fatalf("Scanner error: %v", err)
	}
}
