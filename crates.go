package main

//
//// IMPORTS
//

import (
	// Modules in GOROOT
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	// External modules
	survey "github.com/AlecAivazis/survey/v2"
	copy "github.com/otiai10/copy"
)

//
//// CRATES
//

type Crate struct {
	name        string
	path        string
	hooksDir    string
	targetsDir  string
	tempDir     string
	environment map[string]string
}

func getSelectedCratesFromCLI(crateNames []string, allCrates bool, interactiveSelection bool, multiple bool, program Program) ([]Crate, functionResponse) {
	var selectedCrates []Crate

	if interactiveSelection == true {
		if allCrates != false || crateNames != nil {
			return []Crate{}, functionResponse{
				exitCode:    1,
				logLevel:    "error",
				message:     fmt.Sprintf("Flag '--interactive/-i' cannot be used together with flags '--crate/-c' or '--all/-a'"),
				indentLevel: program.indentLevel,
			}
		}

		availableCrates, response := getUserCrates(program)
		handleFunctionResponse(response, true)

		availableCratesStrings := make([]string, len(availableCrates))
		for i, element := range availableCrates {
			availableCratesStrings[i] = element.name
		}

		var selectedIndices []int
		prompt := &survey.MultiSelect{
			Message: "Select the crate(s)",
			Options: availableCratesStrings,
			Default: selectedIndices,
		}
		err := survey.AskOne(prompt, &selectedIndices, survey.WithPageSize(10))
		if err != nil {
			if err.Error() == "interrupt" {
				return []Crate{}, functionResponse{
					exitCode:    1,
					message:     "Operation cancelled by user",
					logLevel:    "error",
					indentLevel: program.indentLevel,
				}
			}
		}

		selectedCrates := make([]Crate, len(selectedIndices))
		for i, index := range selectedIndices {
			selectedCrates[i] = availableCrates[index]
		}

		if len(selectedCrates) == 0 {
			return selectedCrates, functionResponse{
				exitCode:    1,
				logLevel:    "attention",
				message:     fmt.Sprintf("No crates were selected"),
				indentLevel: program.indentLevel,
			}
		}

		return selectedCrates, functionResponse{exitCode: 0}
	} else {
		if allCrates != false && crateNames != nil {
			return []Crate{}, functionResponse{
				exitCode:    1,
				logLevel:    "error",
				message:     fmt.Sprintf("Conflicting flags: both '--crate/-c' and '--all/-a' flags cannot be specified at the same time"),
				indentLevel: program.indentLevel,
			}
		}

		if allCrates == false && crateNames == nil {
			return []Crate{}, functionResponse{
				exitCode:    1,
				logLevel:    "error",
				message:     fmt.Sprintf("Missing required flag: '--interactive/--i' or '--crate/-c' or '--all/-a' flag must be specified"),
				indentLevel: program.indentLevel,
			}
		}

		if allCrates == true {
			selectedCrates, response := getUserCrates(program)
			handleFunctionResponse(response, true)

			return selectedCrates, functionResponse{exitCode: 0}
		} else {
			selectedCrates = make([]Crate, len(crateNames))
			for i, element := range crateNames {
				selectedCrates[i] = generateCrateObj(element, program)
			}

			// Verify crates
			response := verifyCratesDirectories(selectedCrates, program)
			handleFunctionResponse(response, true)

			return selectedCrates, functionResponse{exitCode: 0}
		}
	}
}

func displayCrateTag(msg string, crate Crate) string {
	return fmt.Sprintf(msg) + fmt.Sprintf(" (") + salmonPink.Sprintf(crate.name) + fmt.Sprintf(")")
}

func getUserCrates(program Program) ([]Crate, functionResponse) {
	crateNames, err := ioutil.ReadDir(program.userCratesDir)
	if err != nil {
		return []Crate{}, functionResponse{
			exitCode:    1,
			message:     fmt.Sprintf("Error reading the crates directory -> " + err.Error()),
			logLevel:    "error",
			indentLevel: program.indentLevel,
		}
	}
	crateNames = filterHiddenFilesAndDirectories(crateNames)

	// Generate a crates array
	crates := make([]Crate, len(crateNames))
	for i, element := range crateNames {
		crates[i] = generateCrateObj(element.Name(), program)
	}

	if len(crates) == 0 {
		return crates, functionResponse{
			exitCode:    1,
			logLevel:    "attention",
			message:     "No crates found",
			indentLevel: program.indentLevel,
		}
	}

	return crates, functionResponse{
		exitCode: 0,
	}
}

func generateCrateObj(crate string, program Program) Crate {
	defaultCrateEnv := map[string]string{
		"PROGRAM_NAME":               program.name,
		"DEFAULT_SHELL":              program.defaultShell,
		"SYNCTROPY_EXEC":             program.exec,
		"SYNCTROPY_UTILS":            fmt.Sprintf("%v utils", program.exec),
		"USER_DATA_DIR":              program.userDataDir,
		"USER_CRATES_DIR":            program.userCratesDir,
		"USER_TEMPLATES_DIR":         program.userTemplatesDir,
		"USER_CRATES_TEMPLATES_DIR":  program.userCratesTemplatesDir,
		"USER_TARGETS_TEMPLATES_DIR": program.userTargetsTemplatesDir,
		"CRATE_NAME":                 crate,
		"CRATE_DIR":                  program.userCratesDir + "/" + crate,
		"CRATE_HOOKS_DIR":            program.userCratesDir + "/" + crate + "/hooks",
		"CRATE_TARGETS_DIR":          program.userCratesDir + "/" + crate + "/targets",
		"CRATE_TEMP_DIR":             program.userCratesDir + "/" + crate + "/.tmp",
	}

	return Crate{
		name:        crate,
		path:        program.userCratesDir + "/" + crate,
		hooksDir:    program.userCratesDir + "/" + crate + "/hooks",
		targetsDir:  program.userCratesDir + "/" + crate + "/targets",
		tempDir:     program.userCratesDir + "/" + crate + "/.tmp",
		environment: defaultCrateEnv,
	}
}

func verifyCratesDirectories(crates []Crate, program Program) functionResponse {
	var failedCrates []Crate

	for _, crate := range crates {
		response := verifyCrateDirectory(crate, program)
		if response.exitCode != 0 {
			failedCrates = append(failedCrates, crate)
		}
	}

	var response functionResponse
	if len(failedCrates) > 0 {
		message := "The following crate(s) was/were not found:\n"
		for _, crate := range failedCrates {
			message = message + fmt.Sprintf("\n    - %s", crate.name)
		}
		response = functionResponse{
			exitCode:    1,
			message:     message,
			logLevel:    "error",
			indentLevel: program.indentLevel,
		}
	} else {
		response = functionResponse{
			exitCode: 0,
		}
	}

	return response
}

func verifyCrateDirectory(crate Crate, program Program) functionResponse {
	if _, err := os.Stat(crate.path); os.IsNotExist(err) {
		return functionResponse{
			exitCode: 1,
		}
	}
	return functionResponse{
		exitCode: 0,
	}
}

func removeCrateTempDirectory(crate Crate, showCrateName bool, notRemoveTempDir bool, program Program) {
	var response functionResponse

	if showCrateName == true {
		showInfoSectionTitle(displayCrateTag("Removing temporary directory", crate), program.indentLevel)
	} else {
		showInfoSectionTitle("Removing temporary directory", program.indentLevel)
	}

	if notRemoveTempDir == true {
		response = functionResponse{
			exitCode:    0,
			message:     "Skipping",
			logLevel:    "attention",
			indentLevel: program.indentLevel + 1,
		}
		handleFunctionResponse(response, false)
		return
	}

	if _, err := os.Stat(crate.tempDir); os.IsNotExist(err) {
		response = functionResponse{
			exitCode:    0,
			logLevel:    "attention",
			message:     fmt.Sprintf("Temporary directory not found"),
			indentLevel: program.indentLevel + 1,
		}
	} else {
		err = os.RemoveAll(crate.tempDir)
		if err != nil {
			response = functionResponse{
				exitCode:    1,
				logLevel:    "error",
				message:     fmt.Sprintf("Failed to remove temporary directory -> %v", err.Error()),
				indentLevel: program.indentLevel + 1,
			}
		} else {
			response = functionResponse{
				exitCode:    0,
				message:     "Finished",
				logLevel:    "success",
				indentLevel: program.indentLevel + 1,
			}
		}
	}

	handleFunctionResponse(response, true)
}

func setupCrateTempDirectory(crate Crate, showCrateName bool, notCreateTempDir bool, program Program) {
	var response functionResponse

	if showCrateName == true {
		showInfoSectionTitle(displayCrateTag("Setting up temporary directory", crate), program.indentLevel)
	} else {
		showInfoSectionTitle("Setting up temporary directory", program.indentLevel)
	}

	if notCreateTempDir == true {
		response = functionResponse{
			exitCode:    0,
			message:     "Skipping",
			logLevel:    "attention",
			indentLevel: program.indentLevel + 1,
		}
		handleFunctionResponse(response, false)
		return
	}

	if _, err := os.Stat(crate.tempDir); err == nil {
		response = functionResponse{
			exitCode:    0,
			message:     "Temporary directory already exists. Recreating it...",
			logLevel:    "attention",
			indentLevel: program.indentLevel + 2,
		}
		handleFunctionResponse(response, false)

		err = os.RemoveAll(crate.tempDir)
		if err != nil {
			response = functionResponse{
				exitCode:    1,
				logLevel:    "error",
				message:     fmt.Sprintf("Failed to recreate temporary directory -> '%v'", err.Error()),
				indentLevel: program.indentLevel + 3,
			}

			handleFunctionResponse(response, true)
		}
	}

	perm := os.FileMode(0755)
	err := os.Mkdir(crate.tempDir, perm)
	if err != nil {
		response = functionResponse{
			exitCode:    1,
			logLevel:    "error",
			message:     fmt.Sprintf("Failed to create temporary directory -> '%v'", err.Error()),
			indentLevel: program.indentLevel + 1,
		}
	} else {
		response = functionResponse{
			exitCode:    0,
			logLevel:    "success",
			message:     "Finished",
			indentLevel: program.indentLevel + 1,
		}
	}

	handleFunctionResponse(response, true)
}

func cratesCreate(program Program) functionResponse {
	// Ask for a crate name
	var crateName string
	var crateTemplate string

	promptName := &survey.Input{
		Message: "Crate name:",
	}
	err := survey.AskOne(promptName, &crateName, survey.WithValidator(survey.MinLength(2)))
	if err != nil {
		if err.Error() == "interrupt" {
			return functionResponse{
				exitCode:    1,
				message:     "Operation cancelled by user",
				logLevel:    "error",
				indentLevel: program.indentLevel,
			}
		}
	}

	// Generate a crate object
	crate := generateCrateObj(crateName, program)

	// Ask for a crate template
	availableTemplates, err := ioutil.ReadDir(program.userCratesTemplatesDir)
	if err != nil {
		return functionResponse{
			exitCode:    1,
			message:     fmt.Sprintf("Failed to read the user's crate templates directory -> " + err.Error()),
			logLevel:    "error",
			indentLevel: program.indentLevel,
		}
	}
	availableTemplates = filterHiddenFilesAndDirectories(availableTemplates)

	availableTemplatesStrings := make([]string, len(availableTemplates)+1)
	// Add a 'scratch' (empty) pseudo-template
	availableTemplatesStrings[0] = "scratch"
	// Add the available templates
	for i, element := range availableTemplates {
		availableTemplatesStrings[i+1] = element.Name()
	}

	promptTemplate := &survey.Select{
		Message: "Crate template:",
		Options: availableTemplatesStrings,
	}
	err = survey.AskOne(promptTemplate, &crateTemplate)
	if err != nil {
		if err.Error() == "interrupt" {
			return functionResponse{
				exitCode:    1,
				message:     "Operation cancelled by user",
				logLevel:    "error",
				indentLevel: program.indentLevel,
			}
		}
	}

	// Verify if the scratch template was selected
	var crateTemplateDir string
	scratchTemplate := false
	if crateTemplate == "scratch" {
		scratchTemplate = true
	} else {
		crateTemplateDir = program.userCratesTemplatesDir + "/" + crateTemplate
	}

	showInfoSectionTitle(displayCrateTag("Creating crate", crate), program.indentLevel)

	// Verify if crate already exists
	response := verifyCrateDirectory(crate, program)
	if response.exitCode == 0 {
		return functionResponse{
			exitCode:    1,
			message:     fmt.Sprintf("Crate '%s' already exists", crate.name),
			logLevel:    "attention",
			indentLevel: program.indentLevel + 1,
		}
	}

	// Create crate directory
	space()
	showInfoSectionTitle("Creating crate directory", program.indentLevel+1)
	err = os.Mkdir(crate.path, 0755)
	if err != nil {
		return functionResponse{
			exitCode:    1,
			message:     fmt.Sprintf("Failed to create crate directory -> " + err.Error()),
			logLevel:    "error",
			indentLevel: program.indentLevel + 2,
		}
	}

	err = os.Mkdir(crate.targetsDir, 0755)
	if err != nil {
		return functionResponse{
			exitCode:    1,
			message:     fmt.Sprintf("Failed to create crate's targets directory -> " + err.Error()),
			logLevel:    "error",
			indentLevel: program.indentLevel + 2,
		}
	}

	response = functionResponse{
		exitCode:    0,
		message:     "Finished",
		logLevel:    "success",
		indentLevel: program.indentLevel + 2,
	}
	handleFunctionResponse(response, false)

	// Copy template to crate directory
	if scratchTemplate == false {
		space()
		showInfoSectionTitle("Copying template to crate directory", program.indentLevel+1)
		copyOptions := copy.Options{
			PreserveTimes: true,
			PreserveOwner: true,
		}

		err = copy.Copy(crateTemplateDir, crate.path, copyOptions)
		if err != nil {
			// Remove crate directory
			_ = os.RemoveAll(crate.path)

			return functionResponse{
				exitCode:    1,
				message:     fmt.Sprintf("Failed to copy template -> " + err.Error()),
				logLevel:    "error",
				indentLevel: program.indentLevel + 2,
			}
		} else {
			response := functionResponse{
				exitCode:    0,
				message:     "Finished",
				logLevel:    "success",
				indentLevel: program.indentLevel + 2,
			}
			handleFunctionResponse(response, false)
		}
	}

	// Run post_create hook (if any)
	space()
	showInfoSectionTitle(lightGray.Sprintf("Running ")+orange.Sprintf("post_create")+lightGray.Sprintf(" hook"), program.indentLevel+1)
	if _, err := os.Stat(crate.hooksDir + "/post_create"); err == nil {
		_, response := runHook(crate.hooksDir+"/post_create", crate.environment, true, true, true, true, true, incrementProgramIndentLevel(program, 1))

		if response.exitCode != 0 {
			handleFunctionResponse(response, false)

			space()

			// Remove crate directory
			showInfoSectionTitle(lightGray.Sprintf("Removing crate directory"), program.indentLevel+1)
			err = os.RemoveAll(crate.path)

			if err != nil {
				response = functionResponse{
					exitCode:    response.exitCode,
					message:     "Failed to remove crate directory",
					logLevel:    "attention",
					indentLevel: program.indentLevel + 2,
				}
			} else {
				response = functionResponse{
					exitCode:    response.exitCode,
					message:     "Removed",
					logLevel:    "attention",
					indentLevel: program.indentLevel + 2,
				}
			}

			return response
		} else {
			response := functionResponse{
				exitCode:    0,
				message:     "Finished",
				logLevel:    "success",
				indentLevel: program.indentLevel + 2,
			}
			handleFunctionResponse(response, false)
		}
	} else {
		response := functionResponse{
			exitCode:    0,
			message:     "Hook not found",
			logLevel:    "attention",
			indentLevel: program.indentLevel + 2,
		}
		handleFunctionResponse(response, false)
	}

	space()

	return functionResponse{
		exitCode:    0,
		message:     "Finished",
		logLevel:    "success",
		indentLevel: program.indentLevel + 1,
	}
}

func cratesRm(crates []Crate, program Program) functionResponse {
	for index, crate := range crates {
		space()

		orange.Println(fmt.Sprintf("(%v/%v)", index+1, len(crates)))
		showInfoSectionTitle(displayCrateTag("Removing", crate), program.indentLevel)

		// Show number of targets and ask for confirmation
		targets, response := getCrateTargets(crate, program)
		if response.exitCode != 0 && response.logLevel != "attention" {
			handleFunctionResponse(response, true)
		}

		if len(targets) > 0 {
			showAttention(fmt.Sprintf("This action will delete %v targets from this crate", len(targets)), program.indentLevel+1)

			userConfirmation := askConfirmation("Enter 'yes/y' to confirm or 'no/n' to cancel the operation", incrementProgramIndentLevel(program, 2))
			if userConfirmation == false {
				return functionResponse{
					exitCode:    1,
					logLevel:    "error",
					message:     "Operation cancelled by user",
					indentLevel: program.indentLevel + 1,
				}
			}

			space()

			// Run the pre_rm hook for each target
			for _, target := range targets {
				showInfoSectionTitle(displayTargetTag(lightGray.Sprintf("Running ")+gray.Sprintf("pre_rm")+lightGray.Sprintf(" hook"), target), program.indentLevel+1)

				if _, err := os.Stat(target.hooksDir + "/pre_rm"); os.IsNotExist(err) {
					response = functionResponse{
						exitCode:    0,
						message:     "Hook not found",
						logLevel:    "attention",
						indentLevel: program.indentLevel + 2,
					}
					handleFunctionResponse(response, false)
				} else {
					program := incrementProgramIndentLevel(program, 1)

					_, response := runHook(target.hooksDir+"/pre_rm", target.environment, true, true, true, true, true, program)
					response.indentLevel = program.indentLevel + 2
					handleFunctionResponse(response, true)
				}

				space()
			}
		} else {
			response = functionResponse{
				exitCode:    0,
				message:     "Crate has no targets. Keeping on...",
				logLevel:    "attention",
				indentLevel: program.indentLevel + 1,
			}
			handleFunctionResponse(response, false)
		}

		// Run pre_rm hook for the crate
		space()

		showInfoSectionTitle(displayCrateTag(lightGray.Sprintf("Running ")+gray.Sprintf("pre_rm")+lightGray.Sprintf(" hook"), crate), program.indentLevel+1)

		if _, err := os.Stat(crate.hooksDir + "/pre_rm"); os.IsNotExist(err) {
			response = functionResponse{
				exitCode:    0,
				message:     "Hook not found",
				logLevel:    "attention",
				indentLevel: program.indentLevel + 2,
			}
			handleFunctionResponse(response, false)
		} else {
			program = incrementProgramIndentLevel(program, 1)

			_, response := runHook(crate.hooksDir+"/pre_rm", crate.environment, true, true, true, true, true, program)
			response.indentLevel = program.indentLevel + 2
			handleFunctionResponse(response, true)
		}

		space()

		err := os.RemoveAll(crate.path)
		if err != nil {
			response = functionResponse{
				exitCode:    1,
				logLevel:    "error",
				message:     "Failed to remove crate -> " + err.Error(),
				indentLevel: program.indentLevel + 1,
			}
		} else {
			response = functionResponse{
				exitCode:    0,
				logLevel:    "success",
				message:     "Finished",
				indentLevel: program.indentLevel + 1,
			}
		}

		handleFunctionResponse(response, true)
	}

	return functionResponse{
		exitCode: 0,
	}
}

func cratesLs(program Program) functionResponse {
	space()
	showInfoSectionTitle("Listing crates", program.indentLevel)

	crates, response := getUserCrates(program)
	if response.exitCode != 0 {
		response.indentLevel = program.indentLevel + 1
		handleFunctionResponse(response, true)
	}

	space()

	for _, crate := range crates {
		// Get optional description (if hook exists)
		crateDescription, response := runHook(crate.hooksDir+"/ls", crate.environment, false, false, false, false, false, program)
		crateDescriptionString := crateDescription.Output

		if response.exitCode == 0 {
			if len(crateDescriptionString) > 0 {
				showText(fmt.Sprintf("- %s (%s)", crate.name, blue.Sprintf(crateDescriptionString)), program.indentLevel+1)
			} else {
				showText(fmt.Sprintf("- %s", crate.name), program.indentLevel+1)
			}
		} else {
			showText(fmt.Sprintf("- %s", crate.name), program.indentLevel+1)
		}
	}

	return functionResponse{
		exitCode: 0,
	}
}

func cratesEdit(crates []Crate, program Program) functionResponse {
	hook := "edit"
	response := cratesRunHooks(crates, []string{hook}, false, false, false, false, false, program)

	return response
}

func cratesView(crates []Crate, program Program) functionResponse {
	hook := "view"
	response := cratesRunHooks(crates, []string{hook}, false, false, false, false, false, program)

	return response
}

func cratesRunHooks(crates []Crate, hooks []string, notCreateTempDir bool, notRemoveTempDir bool, notPrintOutput bool, notPrintEntryCmd bool, notPrintAlerts bool, program Program) functionResponse {
	var response functionResponse

	for index, crate := range crates {
		space()

		orange.Println(fmt.Sprintf("(%v/%v)", index+1, len(crates)))

		showInfoSectionTitle(displayCrateTag("Running hook(s)", crate), program.indentLevel)

		program = incrementProgramIndentLevel(program, 1)

		setupCrateTempDirectory(crate, false, notCreateTempDir, program)

		response = func(crate Crate, hooks []string, notRemoveTempDir bool, notPrintOutput bool, notPrintEntryCmd bool, program Program) functionResponse {
			for _, hook := range hooks {
				space()
				space()

				showInfoSectionTitle(lightGray.Sprintf("Running ")+orange.Sprintf(hook)+lightGray.Sprintf(" hook"), program.indentLevel)

				// Run hook
				if _, err := os.Stat(crate.hooksDir + "/" + hook); os.IsNotExist(err) {
					response = functionResponse{
						exitCode:    1,
						message:     fmt.Sprintf("No '%v' hook found", hook),
						logLevel:    "error",
						indentLevel: program.indentLevel + 1,
					}
					handleFunctionResponse(response, false)
				} else {
					_, hookResponse := runHook(crate.hooksDir+"/"+hook, crate.environment, !notPrintOutput, true, !notPrintEntryCmd, true, !notPrintAlerts, program)

					if hookResponse.exitCode != 0 {
						hookResponse.indentLevel = program.indentLevel + 1
						return hookResponse
					} else {
						response = functionResponse{
							exitCode:    0,
							message:     "Finished",
							logLevel:    "success",
							indentLevel: program.indentLevel + 1,
						}
						handleFunctionResponse(response, false)
					}
				}
			}

			return functionResponse{
				exitCode: 0,
			}
		}(crate, hooks, notRemoveTempDir, notPrintOutput, notPrintEntryCmd, program)

		if response.exitCode != 0 {
			handleFunctionResponse(response, false)

			space()
			space()

			removeCrateTempDirectory(crate, false, notRemoveTempDir, program)

			space()
			finishProgram(response.exitCode)
		} else {
			handleFunctionResponse(response, false)
		}

		space()
		space()

		removeCrateTempDirectory(crate, false, notRemoveTempDir, program)

		program = decrementProgramIndentLevel(program, 1)
	}

	return functionResponse{
		exitCode: 0,
	}
}

func cratesHooksLs(crates []Crate, program Program) functionResponse {
	for index, crate := range crates {
		space()

		orange.Println(fmt.Sprintf("(%v/%v)", index+1, len(crates)))
		showInfoSectionTitle(displayCrateTag("Listing hooks", crate), program.indentLevel)

		hooks, err := ioutil.ReadDir(crate.hooksDir)
		if err != nil {
			return functionResponse{
				exitCode:    1,
				message:     fmt.Sprintf("Error reading the hooks directory -> " + err.Error()),
				logLevel:    "error",
				indentLevel: program.indentLevel + 1,
			}
		}

		// Filter out .entry files
		filteredHooks := make([]os.FileInfo, 0)
		for _, element := range hooks {
			if !strings.HasSuffix(element.Name(), ".entry") {
				filteredHooks = append(filteredHooks, element)
			}
		}
		filteredHooks = filterHiddenFilesAndDirectories(filteredHooks)

		if len(filteredHooks) == 0 {
			showAttention("No hooks found", program.indentLevel)
		}

		for _, element := range filteredHooks {
			// Verify if hook has custom entry command
			var entryCommand string
			customEntryFilePath := crate.hooksDir + "/" + element.Name() + ".entry"
			if _, err = os.Stat(customEntryFilePath); err == nil {
				contents, err := ioutil.ReadFile(customEntryFilePath)
				if err != nil {
					return functionResponse{
						exitCode:    1,
						message:     fmt.Sprintf("Failed to read custom entry configuration file for hook " + element.Name() + " -> " + err.Error()),
						logLevel:    "error",
						indentLevel: program.indentLevel,
					}
				}
				entryCommand = string(contents)
				entryCommand = strings.TrimLeft(entryCommand, "\n")
				entryCommand = strings.TrimRight(entryCommand, "\n")
			} else {
				entryCommand = program.defaultShell
			}
			showText(fmt.Sprintf("- %s (%s)", element.Name(), coral.Sprintf(entryCommand)), program.indentLevel+1)
		}
	}

	return functionResponse{
		exitCode: 0,
	}
}
