package main

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

func AddFileToZip(destinationFilePath string, sourceFilePath string) error {
	zipFile, err := os.Create(destinationFilePath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	sourceFile, err := os.Open(sourceFilePath)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	sourceFileInfo, err := sourceFile.Stat()
	if err != nil {
		return err
	}

	zipHeader, err := zip.FileInfoHeader(sourceFileInfo)
	if err != nil {
		return err
	}

	zipHeader.Name = sourceFilePath
	zipHeader.Method = zip.Deflate

	zipWriterEntry, err := zipWriter.CreateHeader(zipHeader)
	if err != nil {
		return err
	}

	_, err = io.Copy(zipWriterEntry, sourceFile)
	if err != nil {
		return err
	}

	return nil
}

func AddFolderToZip(zipWriter *zip.Writer, sourceFolder string) error {
	err := filepath.Walk(sourceFolder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(sourceFolder, path)
		if err != nil {
			return err
		}

		zipHeader, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		zipHeader.Name = filepath.Join(sourceFolder, relPath)
		zipHeader.Method = zip.Deflate

		zipEntry, err := zipWriter.CreateHeader(zipHeader)
		if err != nil {
			return err
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(zipEntry, file)
		if err != nil {
			return err
		}

		return nil
	})
	return err
}
