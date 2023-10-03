package main

//
//// IMPORTS
//

import (
	// Modules in GOROOT
	"fmt"
	"os"
	"os/user"
	// External modules
)

//
//// USER DATA/CONFIGURATION
//

func getCurrentUser(program Program) (*user.User, functionResponse) {
	usr, err := user.Current()
	if err != nil {
		return nil, functionResponse{
			exitCode:    1,
			message:     fmt.Sprintf("Failed to get current user -> " + err.Error()),
			logLevel:    "error",
			indentLevel: program.indentLevel,
		}
	}

	return usr, functionResponse{
		exitCode: 0,
	}
}

func getCurrentUserHomeDir(program Program) string {
	currentUser, response := getCurrentUser(program)
	handleFunctionResponse(response, true)

	return currentUser.HomeDir
}

func verifyUserDataDirectory(printSectionTitle bool, program Program) functionResponse {
	if printSectionTitle == true {
		showInfoSectionTitle("Verifying user data directory", program.indentLevel)
	}

	createdDirectories := false
	if _, err := os.Stat(program.userDataDir); os.IsNotExist(err) {
		createdDirectories = true
		showAttention("> User data directory not found. Creating...", program.indentLevel+2)
		err := os.Mkdir(program.userDataDir, 0755)
		if err != nil {
			return functionResponse{
				exitCode:    1,
				message:     fmt.Sprintf("Failed to create user data directory -> " + err.Error()),
				logLevel:    "error",
				indentLevel: program.indentLevel + 3,
			}
		}
	}
	if _, err := os.Stat(program.userCratesDir); os.IsNotExist(err) {
		createdDirectories = true
		showAttention("> User crates directory not found. Creating...", program.indentLevel+2)
		err := os.Mkdir(program.userCratesDir, 0755)
		if err != nil {
			return functionResponse{
				exitCode:    1,
				message:     fmt.Sprintf("Failed to create user crates directory -> " + err.Error()),
				logLevel:    "error",
				indentLevel: program.indentLevel + 3,
			}
		}
	}
	if _, err := os.Stat(program.userTemplatesDir); os.IsNotExist(err) {
		createdDirectories = true
		showAttention("> User templates directory not found. Creating...", program.indentLevel+2)
		err := os.Mkdir(program.userTemplatesDir, 0755)
		if err != nil {
			return functionResponse{
				exitCode:    1,
				message:     fmt.Sprintf("Failed to create user templates directory -> " + err.Error()),
				logLevel:    "error",
				indentLevel: program.indentLevel + 3,
			}
		}
	}
	if _, err := os.Stat(program.userCratesTemplatesDir); os.IsNotExist(err) {
		createdDirectories = true
		showAttention("> User crates templates directory not found. Creating...", program.indentLevel+2)
		err := os.Mkdir(program.userCratesTemplatesDir, 0755)
		if err != nil {
			return functionResponse{
				exitCode:    1,
				message:     fmt.Sprintf("Failed to create user crates templates directory -> " + err.Error()),
				logLevel:    "error",
				indentLevel: program.indentLevel + 3,
			}
		}
	}
	if _, err := os.Stat(program.userTargetsTemplatesDir); os.IsNotExist(err) {
		createdDirectories = true
		showAttention("> User targets templates directory not found. Creating...", program.indentLevel+2)
		err := os.Mkdir(program.userTargetsTemplatesDir, 0755)
		if err != nil {
			return functionResponse{
				exitCode:    1,
				message:     fmt.Sprintf("Failed to create user targets templates directory -> " + err.Error()),
				logLevel:    "error",
				indentLevel: program.indentLevel + 3,
			}
		}
	}

	response := functionResponse{
		exitCode:    0,
		logLevel:    "success",
		indentLevel: program.indentLevel + 1,
	}

	if createdDirectories == true {
		response.message = "Finished"
	} else {
		response.message = "Passed"
	}

	return response
}
