package main

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"github.com/fatih/color"
	"github.com/Dartmouth-OpenAV/microservice-framework/framework"
)

// Send a command to the device and process its response
func sendCommand(socketKey string, command []byte, expectedResponse string) ([]byte, error) {
	function := "sendCommand"
	framework.Log(function + " - Command: " + string(command) + " for: " + socketKey)

	// setLineDelimiter(defaultLineDelimiter) // Handled in microservice.go setFrameworkGlobals

	// Write our command
	line := string(command)
	if !framework.WriteLineToSocket(socketKey, line) {
		// the socket connection may have died, so we'll try once to reopen it
		errMsg := function + " - bfsd54w error writing to " + socketKey // + " closing and trying to open the socket again"
		framework.AddToErrors(socketKey, errMsg)
	}

	// Now we need to read the response
	msg := framework.ReadLineFromSocket(socketKey)
	// Shure P300 may inject random responses, so we need to read past unwanted ones
	// See if we got nothing or the expected response
	for (len(msg) != 0) && (!strings.Contains(msg, expectedResponse)){
		// We read something, but it's not what we seek, so try reading again
		keepMaxReadTries := framework.MaxReadTries
		keepReadNoSleepTries := framework.ReadNoSleepTries
		// Tell the framework to not retry much on this read because there probably isn't anything 
		// there to read anyway
		framework.MaxReadTries = 1 
		framework.ReadNoSleepTries = 1 
		msg = framework.ReadLineFromSocket(socketKey)
		framework.MaxReadTries = keepMaxReadTries
		framework.ReadNoSleepTries = keepReadNoSleepTries
		framework.Log("vewada#@ Tried reading again, got: [" + msg + "]")
	}
	framework.Log("DONE reading again")
	// msg = msg
	if len(msg) == 0 { // No response
		errMsg := function + " - asdfasf2323 error reading from " + socketKey 
		framework.AddToErrors(socketKey, errMsg)
		//closeSocketConnection(socketKey)
		return (command), errors.New(errMsg)
	} else { // Got something to process
		// Process what we read
		if len(msg) > 3 {
			// check for leading '<' - a sign of healthy communication
			if msg[0] != '<' {
				errMsg := function + " - jka9npe error didn't find a < in response from " + socketKey
				framework.AddToErrors(socketKey, errMsg)
				// return string(command), errors.New(errMsg)
				return (command), errors.New(errMsg)
			}
		}
		// framework.Log(fmt.Sprintf("Command [%s] got actual response [%s]\n", string(command), msg))
		stringMsg := string(msg)
		respVals := strings.Split(stringMsg, " ") // Shure responses are space delimited

		if strings.Contains(string(command), "SET") {
			// framework.Log(fmt.Sprintf("Set command: %s response: [%s] respVals[1]: %s\n", string(command), stringMsg, respVals[1]))
			if respVals[1] != "REP" {
				errMsg := function + " - error got non-\"REP\" response from " + socketKey + "to SET command"
				framework.AddToErrors(socketKey, errMsg)
				// return string(command), errors.New(errMsg)
				return (command), errors.New(errMsg)
			}
		}

		if strings.Contains(stringMsg, expectedResponse) {
			// Since they match this was presumably responding to our command and we'll use the most recent
			// framework.Log(fmt.Sprintf("Got expected response: [%s], msg is: [%s]\n", expectedResponse, msg))
			//needToRead = false  // success so quit
		} else { // Didn't get expected response, so we'll loop again in case we can read past extraneous responses and find the one we want
			errMsg := color.HiRedString(fmt.Sprintf(
				function+" - 3q2fsvxc Didn't get expected response. The command was: %s response was: %s expectedResponse was: %s\n\n",
				string(command), string(msg)), expectedResponse)
			framework.Log(errMsg)
			// framework.AddToErrors(socketKey, errMsg)
			// return string(command), errors.New(errMsg)
		}
	}

	return []byte(msg), nil
}

func setAudioMute(socketKey string, name string, value string) (string, error) {
	function := "setAudioMute"
	onOff := "not set"

	if value == `"true"` {
		onOff = "ON"
	} else {
		onOff = "OFF"
	}

	expectedResponse := name + " AUDIO_MUTE " + onOff 
	if len(name) == 1{
		expectedResponse = "0" + expectedResponse  // Shure adds a leading zero for single digit channel numbers
	}
	commandStr := []byte("< SET " + expectedResponse + " >")
	resp, err := doShureCommand(socketKey, commandStr, expectedResponse, 4)
	if err != nil {
		errMsg := fmt.Sprintf(function+" - Error fr564b %v", err.Error())
		framework.AddToErrors(socketKey, errMsg)
		return resp, errors.New(errMsg)
	}

	if resp != `"` + onOff + `"` {
		errMsg := fmt.Sprintf(function+" - Invalid mute state returned in response received: %s", resp)
		framework.Log(fmt.Sprintf(color.HiRedString(errMsg)))
		return resp, errors.New(errMsg)
	}

	if resp == `"ON"` {
		return `"true"`, nil
	} else if resp == `"OFF"` { 
		return `"false"`, nil
	} else { // not a legal value
		errMsg := function + " - unrecognized audio mute returned: " + resp + " is not a legal value\n"
		framework.Log(errMsg)
		framework.AddToErrors(socketKey, errMsg)
		return resp, errors.New(errMsg)
	}
}

func getAudioMute(socketKey string, name string) (string, error) {
	function := "getAudioMute"
	// Send the Shure get mute status command
	expectedResponse := name + " AUDIO_MUTE" 
	if len(name) == 1 {
		expectedResponse = "0" + expectedResponse  // Shure adds a leading zero for single digit channel numbers
	}
	commandStr := []byte("< GET " + expectedResponse + " >")
	resp, err := doShureCommand(socketKey, commandStr, expectedResponse, 4)
	if err != nil {
		errMsg := fmt.Sprintf(function+" - Error khT&k %v", err.Error())
		framework.AddToErrors(socketKey, errMsg)
		return resp, errors.New(errMsg)
	}

	if resp == `"ON"` {
		return `"true"`, nil
	} else if resp == `"OFF"` { 
		return `"false"`, nil
	} else { // not a legal value
		errMsg := function + " - unrecognized audio mute returned: " + resp + " is not a legal value\n"
		framework.Log(errMsg)
		framework.AddToErrors(socketKey, errMsg)
		return resp, errors.New(errMsg)
	}
}

func setVolume(socketKey string, name string, level string) (string, error) {
	function := "setVolume"

	levelNoQuotes := strings.Replace(level, `"`, "", -1) // Get rid of the JSON body quotes
	deviceLevel := newTransformVolume(levelNoQuotes)

	// Add leading zeros to the level to get to four digits
	if len(deviceLevel) == 1 {
		deviceLevel = "000" + deviceLevel
	} else if len(deviceLevel) == 2 {
		deviceLevel = "00" + deviceLevel
	} else if len(deviceLevel) == 3 {
		deviceLevel = "0" + deviceLevel
	}
	expectedResponse := name + " AUDIO_GAIN_HI_RES " + deviceLevel 
	if len(name) == 1{
		expectedResponse = "0" + expectedResponse  // Shure adds a leading zero for single digit channel numbers
	}
	commandStr := []byte("< SET " + expectedResponse + " >")
	respGain, err := doShureCommand(socketKey, commandStr, expectedResponse, 4)
	if err != nil {
		errMsg := fmt.Sprintf(function+" - Error kjh867@ %v", err.Error())
		framework.AddToErrors(socketKey, errMsg)
		return respGain, errors.New(errMsg)
	}

	// No errors, so untransform the volume we got from Shure
	respGain = `"` + newUnTransformVolume(strings.Replace(respGain, `"`, "", -1)) + `"`
	if !almostEqualorEqual(respGain, level) { // If what the device said it set doesn't match what we told it to set
		errMsg := fmt.Sprintf(function + " - Invalid gain state returned in response received b345df")
		framework.AddToErrors(socketKey, errMsg)
		return respGain, errors.New(errMsg)
	}

	return respGain, err
}

func getVolume(socketKey string, name string) (string, error) {
	function := "getVolume"
	// Send the Shure volume get command
	
	expectedResponse := name + " AUDIO_GAIN_HI_RES "
	if len(name) == 1{
		expectedResponse = "0" + expectedResponse  // Shure adds a leading zero for single digit channel numbers
	}
	commandStr := []byte("< GET " + expectedResponse + " >")
	respGain, err := doShureCommand(socketKey, commandStr, expectedResponse, 4)
	if err != nil {
		errMsg := fmt.Sprintf(function+" - Error khT&k %v", err.Error())
		framework.AddToErrors(socketKey, errMsg)
		return respGain, errors.New(errMsg)
	}

	// Untransform to OpenAV standard 0 - 100 range
	respGain = `"` + newUnTransformVolume(strings.Replace(respGain, `"`, "", -1)) + `"`

	return respGain, err
}


func setMatrixVolume(socketKey string, in string, out string, level string) (string, error) {
	function := "setMatrixVolume"

	levelNoQuotes := strings.Replace(level, `"`, "", -1) // Get rid of the JSON body quotes
	deviceLevel := newTransformVolume(levelNoQuotes)
	// Add leading zeros to the level to get to four digits
	if len(deviceLevel) == 1 {
		deviceLevel = "000" + deviceLevel
	} else if len(deviceLevel) == 2 {
		deviceLevel = "00" + deviceLevel
	} else if len(deviceLevel) == 3 {
		deviceLevel = "0" + deviceLevel
	}
	if len(in) == 1 {
		in = "0" + in  // Shure adds a leading zero for single digit channel numbers
	}
	if len(out) == 1 {
		out = "0" + out  
	}
	expectedResponse := in + " MATRIX_MXR_GAIN " + out + " " + deviceLevel 
	commandStr := []byte("< SET " + expectedResponse + " >")	// Send the Shure "Set Matrix Mixer Gain" command
	respGain, err := doShureCommand(socketKey, commandStr, expectedResponse, 5)
	if err != nil {
		errMsg := fmt.Sprintf(function+" - Error 342ffs %v", err.Error())
		framework.AddToErrors(socketKey, errMsg)
		return respGain, errors.New(errMsg)
	}

	// No errors, so untransform the volume we got from Shure
	respGain = `"` + newUnTransformVolume(strings.Replace(respGain, `"`, "", -1)) + `"`
	// framework.Log("respGain: [" + respGain + "]" + "level: [" + level + "]")
	if !almostEqualorEqual(respGain, level) { // If what the device said it set doesn't match what we told it to set
		errMsg := fmt.Sprintf(function + " - Invalid gain state returned in response received fewa34")
		framework.AddToErrors(socketKey, errMsg)
		return respGain, errors.New(errMsg)
	}	
	
	return respGain, nil
}

func almostEqualorEqual(val1 string, val2 string) bool {
	// We have a rounding error problem after transforming volumes to Shure values and back, so we need
	//  to allow numbers to be a little off in the return values.
	intVal1 := stringToInt(val1)
	intVal2	:= stringToInt(val2)

	diff := intVal1 - intVal2
	if diff > -2 && diff < 2 {  // The two numbers are within 2 of each other (1 would probably be good enough)
		return true
	} else {
		return false
	}
}

func stringToInt(theString string) int {
	noQuotesString := strings.Replace(theString, `"`, "", -1)
	theInt, err := strconv.Atoi(noQuotesString)
	if err != nil {
		framework.Log(fmt.Sprintf("stringToInt - Error 789k;jkfd converting string %s to int", theString))
		return -999
	} else {
		return theInt
	}
}

func getMatrixVolume(socketKey string, in string, out string) (string, error) {
	function := "setMatrixVolume"

	if len(in) == 1 {
		in = "0" + in  // Shure adds a leading zero for single digit channel numbers
	}
	if len(out) == 1 {
		out = "0" + out  
	}
	expectedResponse := in + " MATRIX_MXR_GAIN " + out
	commandStr := []byte("< GET " + expectedResponse + " >")	
	// Send the Shure "Get Matrix Mixer Gain" command
	respGain, err := doShureCommand(socketKey, commandStr, expectedResponse, 5)
	if err != nil {
		errMsg := fmt.Sprintf(function+" - Error f354fda %v", err.Error())
		framework.AddToErrors(socketKey, errMsg)
		return respGain, errors.New(errMsg)
	}

	// No errors, so untransform the volume we got from Shure
	respGain = `"` + newUnTransformVolume(strings.Replace(respGain, `"`, "", -1)) + `"`
	
	return respGain, nil
}


func setMatrixMute(socketKey string, in string, out string, mute string) (string, error) {
	function := "setMatrixMute"

	if len(in) == 1 {
		in = "0" + in  // Shure adds a leading zero for single digit channel numbers
	}
	if len(out) == 1 {
		out = "0" + out  
	}
	var onOff = "not set"
	if mute == `"` + "true" + `"` {
		onOff = "OFF"
	} else {
		onOff = "ON"
	}
	expectedResponse := in + " MATRIX_MXR_ROUTE " + out + " " + onOff 
	commandStr := []byte("< SET " + expectedResponse + " >")	

	// Send the Shure "Set Matrix Mixer Gain" command
	respMute, err := doShureCommand(socketKey, commandStr, expectedResponse, 5)
	if err != nil {
		errMsg := fmt.Sprintf(function+" - Error bfg534h %v", err.Error())
		framework.AddToErrors(socketKey, errMsg)
		return respMute, errors.New(errMsg)
	}

	retVal := "not set"
	if respMute == `"` + "ON" + `"` {
		retVal = "false" // NOT muted when connected
	} else if respMute == `"` + "OFF" + `"` {
		retVal = "true"
	} else {
		errMsg := fmt.Sprintf(function+" - Invalid mute value returned: %s", respMute)
		framework.AddToErrors(socketKey, errMsg)
		return string(respMute), errors.New(errMsg)

	}

	return `"` + retVal + `"`, nil
}

func getMatrixMute(socketKey string, in string, out string) (string, error) {
	function := "getMatrixMute"

	if len(in) == 1 {
		in = "0" + in  // Shure adds a leading zero for single digit channel numbers
	}
	if len(out) == 1 {
		out = "0" + out  
	}
	expectedResponse := in + " MATRIX_MXR_ROUTE " + out + " "
	commandStr := []byte("< GET " + expectedResponse + " >")	
	// Send the Shure "Get Matrix Mixer Gain" command
	respMute, err := doShureCommand(socketKey, commandStr, expectedResponse, 5)
	if err != nil {
		errMsg := fmt.Sprintf(function+" - Error 34tgsr %v", err.Error())
		framework.AddToErrors(socketKey, errMsg)
		return respMute, errors.New(errMsg)
	}

	retVal := "not set"
	if respMute == `"` + "ON" + `"` {
		retVal = "false" // NOT muted when connected
	} else if respMute == `"` + "OFF" + `"` {
		retVal = "true"
	} else {
		errMsg := fmt.Sprintf(function+" - Invalid mute value returned: %s", respMute)
		framework.AddToErrors(socketKey, errMsg)
		return string(respMute), errors.New(errMsg)

	}

	return `"` + retVal + `"`, nil
}

// The following two functions keep our volume in the range 0 - 100 and attempt to make volume change appear linear
func newTransformVolume(vol string) string {
	function := "transformVolume"
	// Make the volume into the range Shure uses
	intVol, err := strconv.Atoi(vol)
	if err != nil {
		errMsg := fmt.Sprintf(function+" - Error converting volume r23dfs&* %v", err.Error())
		return errMsg
	}
	xform := math.Log10(1.0+(float64(intVol)/11.0)) * 100.0
	return (strconv.Itoa(int(xform) * 14)) // Shure is 0 - 1400, we use 0 - 100
}

func newUnTransformVolume(vol string) string {
	function := "unTransformVolume"
	// Make the volume back to our 0 - 100 range
	intVol, err := strconv.Atoi(vol)
	if err != nil {
		errMsg := fmt.Sprintf(function+" - Error converting volume qwrefsd#@ %v", err.Error())
		return errMsg
	}

	xform := (math.Pow(10, float64(intVol)/1400.0) - 1.0) * 11.0 // Shure is 0 - 1400, we use 0 - 100
	intxform := int(xform)
	if intxform > 0 {
		intxform = intxform + 1 // Adding 1 is empirical, but gets us closer -
		// these are still occasionaly aproximate, presumably because of round off error)
	}

	return (strconv.Itoa(intxform))
}


func doShureCommand(socketKey string, commandStr []byte, expectedResponse string, index int) (string, error) {
	function := "doShureCommand"
	resp, err := sendCommand(socketKey, commandStr, expectedResponse)
	// fmt.Printf("Got resp: " + string(resp) + "\n")
	if err != nil {
		errMsg := fmt.Sprintf(function+" - Error sending command lkw32sdf %v", err.Error())
		framework.AddToErrors(socketKey, errMsg)
		return string(resp), errors.New(errMsg)
	}

	// Validate we got the right response type back and for the right parameters
	if !strings.Contains(string(resp), expectedResponse) {
		errMsg := fmt.Sprintf(function + " - Invalid command returned in response received 43drgsfd: " + string(resp))
		framework.AddToErrors(socketKey, errMsg)
		return string(resp), errors.New(errMsg)
	}

	// Now we need to parse the response from the Shure DSP
	respVals := strings.Split(string(resp), " ") // It was nice of Shure to delimit everything with spaces
	if len(respVals) < 5 {
		errMsg := function + " - Not enough tokens returned u98knjl"
		framework.AddToErrors(socketKey, errMsg)
		return string(resp), errors.New(errMsg)
	}

	return `"` + respVals[index] + `"`, nil
}
