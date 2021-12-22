package bridge

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type Bridge struct {
	ID            string `json:"id"`
	IP            string `json:"internalipaddress"`
	api_key       string
	friendly_name string
}

func SearchForBridges() *[]Bridge {

	// Empty slice of Bridges
	var foundBridges []Bridge

	// Get the list of briges from Hue
	// Need to implement MDNS so we don't rely on calling out to Hue in future
	resp, err := http.Get("https://discovery.meethue.com/")
	if err != nil {
		log.Fatalln(err)
	}

	// Read the response as bytes
	respBodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	// Convert the bytes to JSON and add each entry to the Bridges slice
	json.Unmarshal(respBodyBytes, &foundBridges)

	return &foundBridges
}

func AddNewBridge() {
	reader := bufio.NewReader(os.Stdin)
	foundBridges := SearchForBridges()

	if len(*foundBridges) > 0 {
		fmt.Println("Found the following bridges:")
		for i := range *foundBridges {
			fmt.Printf("%v: IP: %v Bridge ID: %v\n", i, (*foundBridges)[i].IP, (*foundBridges)[i].ID)
		}
		fmt.Print("Please select the number of the Bridge you would like to add, or type '*' to attempt to connect to all bridges: ")
		input, _ := reader.ReadString('\n')
		selection := strings.TrimSuffix("\n", input)
		selectionAsInt, _ := strconv.Atoi(selection)

		if selection == "*" {
			fmt.Println("Attempting to connect to all listed Bridges. Please press the Link button on the Bridge you wish to activate...")
		} else if selectionAsInt >= 0 || selectionAsInt < len(*foundBridges) {
			fmt.Printf("Adding Bridge: %v\n", (*foundBridges)[selectionAsInt].IP)
		}

	} else {
		fmt.Println("No bridges were found on this network. Please confirm the Bridge is set up and has access to the internet.")
		fmt.Println("Huebot currently relies on the Hue Bridge having access to the internet, and being able to contact https://discovery.meethue.com")
	}

}
