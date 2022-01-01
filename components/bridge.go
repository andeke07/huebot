// Functions relating to interacting with a Hue Bridge
package bridge

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type Bridge struct {
	ID            string `json:"id"`
	IP            string `json:"internalipaddress"`
	API_key       string
	Friendly_name string
}

// A scruct representing the required information to request an API key from the Bridge
type APIRequestPayload struct {
	DeviceType        string `json:"devicetype"`
	Generateclientkey bool   `json:"generateclientkey"`
}

// Search the network for Hue Bridges.
// Currently uses the discovery.meethue.com endpoint for a list of bridges on the current network.
// May not work if the User's PC is on a proxy or the Bridges do not have access to the internet
// Returns a slice of type Bridge
func SearchForBridges() []Bridge {
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

	return foundBridges
}

// Get the API Key from a specified bridge (bridgeSelection) from the slice of Bridges (bridges). If bridgeSelection == 999, attempt to
// connect to all bridges in the passed Bridge slice (the user can press the link button on the bridge they wish to use)
func GetBridgeAPIKey(bridges []Bridge, bridgeSelection int) bool {
	// awaitingAPIKey := true
	bridgeApiEndpoint := "http://" + bridges[bridgeSelection].IP + "/api"
	fmt.Println("Getting API Key from " + bridgeApiEndpoint)

	p := APIRequestPayload{"huebot", true}
	b, _ := json.Marshal(p)

	ticker := time.NewTicker(5 * time.Second)
	for range ticker.C {
		response, err := http.Post(bridgeApiEndpoint, "text/json", bytes.NewBuffer(b))
		if err != nil {
			log.Fatal(err)
		}
		bodyBytes, _ := io.ReadAll(response.Body)

		fmt.Println(string(bodyBytes))
	}

	// If we couldn't get the API Key
	return false
}
func AddNewBridge() {
	reader := bufio.NewReader(os.Stdin)
	foundBridges := SearchForBridges()

	if len(foundBridges) > 0 {
		fmt.Println("Found the following bridges:")
		for i := range foundBridges {
			fmt.Printf("%v: IP: %v Bridge ID: %v\n", i, foundBridges[i].IP, foundBridges[i].ID)
		}
		fmt.Print("Please select the number of the Bridge you would like to add, or type '*' to attempt to connect to all bridges: ")
		input, _ := reader.ReadString('\n')
		selection := strings.TrimSuffix(input, "\n")
		selectionAsInt, _ := strconv.Atoi(selection)

		if selection == "*" {
			fmt.Println("Attempting to connect to all listed Bridges. Please press the Link button on the Bridge you wish to activate...")
			GetBridgeAPIKey(foundBridges, 999)

		} else if selection != "*" && selectionAsInt >= 0 && selectionAsInt < len(foundBridges) {
			fmt.Printf("Adding Bridge: %v\n", (foundBridges)[selectionAsInt].IP)
			// Get bridge API Key
			if GetBridgeAPIKey(foundBridges, selectionAsInt) {

			} else {
				fmt.Println("There was an error connecting to the Bridge.")
			}
			// Give it a friendly name
			// Add it to the database
		}

	} else {
		fmt.Println("No bridges were found on this network. Please confirm the Bridge is set up and has access to the internet.")
		fmt.Println("Huebot currently relies on the Hue Bridge having access to the internet, and being able to contact https://discovery.meethue.com")
	}

}
