package main

import (
	"fmt"
	"log"
	"rs201/utils"
	"strconv"
	"strings"

	"github.com/pterm/pterm"
)

func main() {
	// load config
	config, err := utils.LoadConfig()
	if err != nil {
		log.Fatal("Error loading config:", err)
	}
	// print config variable
	fmt.Println(config)

	// Initialize an empty slice to hold the options
	var options []string

	// loop through config and add to options
	for _, relay := range config.Relays {
		options = append(options, fmt.Sprintf("%s [%s @ %s]", relay.IP, relay.Name, relay.Group))
	}

	// Use PTerm's interactive select feature to present the options to the user and capture their selection
	// The Show() method displays the options and waits for the user's input
	selectedOption, _ := pterm.DefaultInteractiveSelect.WithDefaultText("Please select relay").WithOptions(options).Show()

	// split on space to take the ip from the first part
	relayIP := strings.Split(selectedOption, " ")[0]

	// relay := config.GetRelayByIP(relayIP)

	pterm.Printfln("Selected: %s", pterm.Green(relayIP))

	// Display the selected option to the user with a green color for emphasis
	pterm.Info.Printfln("Selected option: %s", pterm.Green(relayIP))

	// in a loop get commands from user
	for {
		// Create a menu with commands get status, turn on, and turn off
		command, _ := pterm.DefaultInteractiveSelect.WithOptions([]string{"Get Status", "Raw", "Sync Reset", "Async Reset", "Exit"}).Show()

		// Display the selected command to the user with a green color for emphasis
		pterm.Info.Printfln("Selected command: %s", pterm.Green(command))

		// switch on command
		switch command {
		case "Get Status":
			// get status
			status, err := utils.GetStatus(relayIP)
			if err != nil {
				pterm.Error.Printfln("Error getting status for %s: %s", pterm.Green(relayIP), pterm.Red(err))
			}
			formatStatus(status)
		case "Raw":
			// raw
			// get command from user
			command, _ := pterm.DefaultInteractiveTextInput.WithMultiLine(false).WithDefaultText("Enter command").Show()
			// run command
			status, err := utils.Raw(relayIP, command)
			if err != nil {
				pterm.Error.Printfln("Error running command for %s: %s", pterm.Green(relayIP), pterm.Red(err))
			}
			formatStatus(status)
		case "Sync Reset":
			ch := promptChannel(true)
			delay, err := promptDelay()
			if err != nil {
				pterm.Error.Printfln("Error getting delay: %s", pterm.Red(err))
				continue
			}
			if delay <= 0 {
				pterm.Error.Printfln("Delay must be greater than 0")
				continue
			}
			pterm.Info.Printfln("Resetting CH%d for %s", ch, pterm.Green(relayIP))
			state, err := utils.Reset(relayIP, ch, delay)
			if err != nil {
				pterm.Error.Printfln("Error running command for %s: %s", pterm.Green(relayIP), pterm.Red(err))
			}
			pterm.Info.Printfln("Status: %s", pterm.Green(state))
		case "Async Reset":
			ch := promptChannel(false)
			delay, err := promptDelay()
			if err != nil {
				pterm.Error.Printfln("Error getting delay: %s", pterm.Red(err))
				continue
			}
			if delay <= 0 {
				pterm.Error.Printfln("Delay must be greater than 0")
				continue
			}
			pterm.Info.Printfln("Resetting CH%d for %s", ch, pterm.Green(relayIP))
			state, err := utils.ResetWithDelay(relayIP, ch, delay)
			if err != nil {
				pterm.Error.Printfln("Error running command for %s: %s", pterm.Green(relayIP), pterm.Red(err))
			}
			pterm.Info.Printfln("Status: %s", pterm.Green(state))
		case "Exit":
			// exit
			pterm.Info.Printfln("Exiting")
			return
		}

	}
}

func formatStatus(status string) {
	// convert the 8 character string to table format
	tableData := pterm.TableData{
		{"CH1", "CH2", "CH3", "CH4", "CH5", "CH6", "CH7", "CH8"},
		{"0", "0", "0", "0", "0", "0", "0", "0"},
	}

	for i := 0; i < 8; i++ {
		if status[i] == '1' {
			tableData[1][i] = "1"
		}
	}

	pterm.DefaultTable.WithHasHeader().WithBoxed().WithData(tableData).Render()
}

func promptChannel(showAll bool) uint8 {
	options := []string{"CH1", "CH2"}

	if showAll {
		options = append([]string{"ALL"}, options...)
	}

	// get channel from user
	channel, _ := pterm.DefaultInteractiveSelect.WithDefaultText("select channel").WithOptions(options).Show()
	// convert to int
	channelInt := uint8(0)
	switch channel {
	case "ALL":
		channelInt = 0
	case "CH1":
		channelInt = 1
	case "CH2":
		channelInt = 2
	}
	return channelInt
}

func promptDelay() (int, error) {
	// get delay from user
	delayStr, _ := pterm.DefaultInteractiveTextInput.WithDefaultText("select delay").Show()
	// convert to int
	return strconv.Atoi(delayStr)
}
