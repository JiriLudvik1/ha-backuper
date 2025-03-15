package main

import (
	"archive/zip"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// A constant list of paths to ignore during compression
var ignoredPaths = []string{
	"auth/",
	"deps/",
	"home-assistant_v2.db",
	"log/",
	"tts/",
}

func shouldIgnore(filePath string) bool {
	for _, ignored := range ignoredPaths {
		if strings.HasPrefix(filepath.ToSlash(filePath), ignored) {
			return true
		}
	}
	return false
}

func CompressFolder(folderPath, destinationZipPath string) error {
	zipFile, err := os.Create(destinationZipPath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	err = filepath.WalkDir(folderPath, func(filePath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relativePath := filepath.ToSlash(filepath.Join(".", filePath[len(folderPath):]))

		if shouldIgnore(relativePath) {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		info, err := d.Info()
		if err != nil {
			return err
		}

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
