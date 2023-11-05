package main

import (
	"archive/zip"
	"fmt"
	"os"
)

func BackupFolders(settings Settings) {

	if len(settings.Folders) == 0 {
		return
	}

	destinationZip := "folders.zip"

	zipFile, err := os.Create(destinationZip)
	if err != nil {
		fmt.Println(err)
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	for _, folder := range settings.Folders {
		AddFolderToZip(zipWriter, folder)
	}

	fmt.Println("Folder backup complete!")
}
