package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func main() {
	// Specify the file containing the list of file paths
	listFile := "file_list.txt"  // adjust this to your file with paths
	outputFile := "combined.txt" // the output file

	// Open the list file
	file, err := os.Open(listFile)
	if err != nil {
		log.Fatalf("Failed to open file list: %v", err)
	}
	defer file.Close()

	// Create or open the output file
	outFile, err := os.Create(outputFile)
	if err != nil {
		log.Fatalf("Failed to create output file: %v", err)
	}
	defer outFile.Close()

	// Use a scanner to read the file paths from the list file
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		filePath := scanner.Text()
		log.Printf(filePath)
		// Read the contents of each file
		data, err := ioutil.ReadFile(filepath.Clean(filePath))
		if err != nil {
			log.Printf("Failed to read file %s: %v", filePath, err)
			continue
		}

		// Write the contents to the output file
		if _, err := outFile.Write(data); err != nil {
			log.Printf("Failed to write to output file: %v", err)
		}

		// Optionally add a separator between file contents
		outFile.WriteString("\n\n=== End of " + filePath + " ===\n\n")
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Failed to read from list file: %v", err)
	}

	fmt.Println("Files successfully concatenated into", outputFile)
}
