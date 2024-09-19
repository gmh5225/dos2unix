// DOS2UNIX
//
// This program converts text files from DOS/Windows format to Unix format
// by changing line endings from CRLF (\r\n) to LF (\n).
//
// Usage:
//   go run dos2unix.go <filename>
//
// Example:
//   go run dos2unix.go example.txt
//
// The program will modify the file in-place. Make sure you have a backup
// of important files before running this program.

package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run dos2unix.go <filename>")
		os.Exit(1)
	}

	filename := os.Args[1]
	err := dos2unix(filename)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("File conversion successful")
}

func dos2unix(filename string) error {
	// Open input file
	input, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer input.Close()

	// Create temporary output file in the same directory as the input file
	tempFile := filename + ".tmp"
	output, err := os.Create(tempFile)
	if err != nil {
		return err
	}
	defer output.Close()

	reader := bufio.NewReader(input)
	writer := bufio.NewWriter(output)

	// Read and convert
	for {
		b, err := reader.ReadByte()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if b == '\r' {
			nextByte, err := reader.Peek(1)
			if err == nil && nextByte[0] == '\n' {
				reader.ReadByte() // Discard \r
				writer.WriteByte('\n')
			} else {
				writer.WriteByte('\r')
			}
		} else {
			writer.WriteByte(b)
		}
	}

	writer.Flush()

	// Close files before renaming
	input.Close()
	output.Close()

	// Replace the original file
	return os.Rename(output.Name(), filename)
}
