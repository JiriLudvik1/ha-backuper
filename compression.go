package main

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

func CompressFolder(folderPath, destinationZipPath string) error {
	zipFile, err := os.Create(destinationZipPath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	err = filepath.Walk(folderPath, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relativePath := filepath.ToSlash(filepath.Join(".", filePath[len(folderPath):]))
		if info.IsDir() {
			relativePath += "/"
		}
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		header.Name = relativePath
		header.Method = zip.Deflate

		if info.IsDir() {
			_, err = zipWriter.CreateHeader(header)
			if err != nil {
				return err
			}
		} else {
			writer, err := zipWriter.CreateHeader(header)
			if err != nil {
				return err
			}

			file, err := os.Open(filePath)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(writer, file)
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
