package main

import "time"

const tempFolder = "/temp"

func main() {
	config, err := LoadConfig()
	if err != nil {
		panic(err)
	}

	now := time.Now()
	archiveName := now.Format("2006-01-02T15-04-05") + ".zip"
	err = CompressFolder(config.HomeAssistantPath, tempFolder+archiveName)
	if err != nil {
		panic(err)
	}
}
