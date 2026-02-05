package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	command := os.Args[1]

	switch command {
	case "migrate":
		runMigration()
	case "import":
		if len(os.Args) < 3 {
			log.Fatal("Usage: cli import <directory>")
		}
		importDocuments(os.Args[2])
	case "status":
		showStatus()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
	}
}

func printUsage() {
	fmt.Println("Eino Study CLI")
	fmt.Println("\nUsage:")
	fmt.Println("  cli migrate           - Run database migrations")
	fmt.Println("  cli import <dir>      - Import documents from directory")
	fmt.Println("  cli status            - Show system status")
}

func runMigration() {
	fmt.Println("Running migrations...")
	// TODO: Implement migration logic
	fmt.Println("Migrations completed")
}

func importDocuments(dir string) {
	fmt.Printf("Importing documents from: %s\n", dir)
	// TODO: Implement document import logic
	fmt.Println("Import completed")
}

func showStatus() {
	fmt.Println("System Status:")
	// TODO: Implement status check logic
	fmt.Println("Database: Connected")
	fmt.Println("Documents: 0")
}