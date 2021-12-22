package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	bridge "github.com/andeke07/huebot/components"
)

func main() {
	// Set up a reader to read from stdio
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Welcome to Huebot!")
	fmt.Println("Currently no Bridges configured")
	fmt.Print("Would you like to configure a bridge now? (y/N): ")

	// Get user's input
	input, _ := reader.ReadString('\n')

	if (strings.TrimSuffix(strings.ToLower(input), "\n")) == "y" {
		fmt.Println("Searching for Bridges...")
		bridge.AddNewBridge()
	} else {
		fmt.Println("Thanks for using Huebot!")
	}

}
