package print

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type FolderScan struct {
	Path     string               // Chemin du dossier à scanner
	FilesPdf map[string]time.Time // Fichiers PDF trouvés dans le dossier
}

// NewFolderScan crée une nouvelle instance de FolderScan avec le chemin spécifié
func NewFolderScan(path string) *FolderScan {
	return &FolderScan{Path: path}
}

// ScanPDFFiles scanne le dossier spécifié et collecte les fichiers PDF dans FilesPdf
func (fs *FolderScan) ScanPDFFiles() error {
	// Réinitialiser la liste des fichiers PDF
	fs.FilesPdf = nil

	// Ouvrir le dossier spécifié
	dir, err := os.Open(fs.Path)
	if err != nil {
		return err
	}
	defer dir.Close()

	// Lister les fichiers dans le dossier
	files, err := dir.Readdir(-1)
	if err != nil {
		return err
	}

	// Parcourir les fichiers pour collecter les fichiers PDF
	fs.FilesPdf = make(map[string]time.Time)
	for _, file := range files {
		if !file.IsDir() && isPdfFile(file.Name()) {
			fs.FilesPdf[file.Name()] = file.ModTime()
		}
	}

	return nil
}

func (fs *FolderScan) GetFilesPathScanned() []string {
	var filesPath []string
	for nameFile := range fs.FilesPdf {
		filesPath = append(filesPath, fmt.Sprintf("%s\\%s", fs.Path, nameFile))
	}

	return filesPath
}

func (fs *FolderScan) CountFilesPdfScanned() int {
	return len(fs.FilesPdf)
}

func isPdfFile(filename string) bool {
	return filepath.Ext(filename) == ".pdf"
}
