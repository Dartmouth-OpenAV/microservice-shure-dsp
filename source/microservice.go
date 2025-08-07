package main

import (
	"errors"
	"github.com/Dartmouth-OpenAV/microservice-framework/framework"
)

func setFrameworkGlobals() {
	// Define device specific globals
	framework.DefaultSocketPort = 2202    // Shure's default socket port
	framework.GlobalDelimiter = 0x3E // ">" is Shure's line delimiter for socket commands

	// globals that change modes in the microservice framework:
	framework.CheckFunctionAppendBehavior = "Remove older instance"  // So the most recent set takes precedence
	framework.MicroserviceName = "OpenAV Shure DSP Microservice"

	framework.RegisterMainGetFunc(doDeviceSpecificGet)
	framework.RegisterMainSetFunc(doDeviceSpecificSet)
}

// Every microservice using this golang microservice framework needs to provide this function to invoke functions to do sets.
// socketKey is the network connection for the framework to use to communicate with the device.
// setting is the first parameter in the URI.
// arg1 are the second and third parameters in the URI.
//   Example PUT URIs that will result in this function being invoked:
// 	 ":address/:setting/"
//   ":address/:setting/:arg1"
//   ":address/:setting/:arg1/:arg2"
func doDeviceSpecificSet(socketKey string, setting string, arg1 string, arg2 string, arg3 string) (string, error) {
	function := "doDeviceSpecificSet"

	// Add a case statement for each set function your microservice implements.  These calls can use 0, 1, or 2 arguments.
	switch setting {
	case "volume":
		return setVolume(socketKey, arg1, arg2)     // arg1 is output, arg2 is value to set
	case "audiomute":
		return setAudioMute(socketKey, arg1, arg2)  // arg1 is output, arg2 is value to set
	case "matrixmute":
		return setMatrixMute(socketKey, arg1, arg2, arg3)  // arg1 is input, arg2 is output, arg3 is value to set
	case "matrixvolume":
		return setMatrixVolume(socketKey, arg1, arg2, arg3)  // arg1 is input, arg2 is output, arg3 is value to set
	}

	// If we get here, we didn't recognize the setting.  Send an error back to the config writer who had a bad URL.
	errMsg := function + " - unrecognized setting in URI: " + setting
	framework.AddToErrors(socketKey, errMsg)
	err := errors.New(errMsg)
	return setting, err
}

// Every microservice using this golang microservice framework needs to provide this function to invoke functions to do gets.
// socketKey is the network connection for the framework to use to communicate with the device.
// setting is the first parameter in the URI.
// arg1 are the second and third parameters in the URI.
//   Example GET URIs that will result in this function being invoked:
// 	 ":address/:setting/"
//   ":address/:setting/:arg1"
//   ":address/:setting/:arg1/:arg2"
// Every microservice using this golang microservice framework needs to provide this function to invoke functions to do gets.
func doDeviceSpecificGet(socketKey string, setting string, arg1 string, arg2 string) (string, error) {
	function := "doDeviceSpecificGet"

	switch setting {
	case "volume":
		return getVolume(socketKey, arg1)    // arg1 is output 
	case "audiomute":
		return getAudioMute(socketKey, arg1) // arg1 is output
	case "matrixmute":
		return getMatrixMute(socketKey, arg1, arg2)	// arg1 is input, arg2 is output
	case "matrixvolume":
		return getMatrixVolume(socketKey, arg1, arg2)	// arg1 is input, arg2 is output
	case "healthcheck":
		return healthCheck(socketKey)
	}

	// If we get here, we didn't recognize the setting.  Send an error back to the config writer who had a bad URL.
	errMsg := function + " - unrecognized setting in URI: " + setting
	framework.AddToErrors(socketKey, errMsg)
	err := errors.New(errMsg)
	return setting, err
}

func main() {
	setFrameworkGlobals()
	framework.Startup()
 }
 
