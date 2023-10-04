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
//// TARGETS
//

type Target struct {
	crate        Crate
	name         string
	path         string
	hooksDir     string
	tempDir      string
	disabledPath string
	environment  map[string]string
}

func getSelectedTargetsFromCLI(crateName string, targetNames []string, allTargets bool, interactiveSelection bool, multiple bool, program Program) (Crate, []Target, functionResponse) {
	var selectedCrate Crate
	var selectedTargets []Target

	if crateName == "" && interactiveSelection == false {
		return Crate{}, []Target{}, functionResponse{
			exitCode:    1,
			logLevel:    "error",
			message:     fmt.Sprintf("Flag '--crate/-c' or '--interactive/-i' should be specified"),
			indentLevel: program.indentLevel,
		}
	}

	if interactiveSelection == true {
		if allTargets != false || targetNames != nil || crateName != "" {
			return Crate{}, []Target{}, functionResponse{
				exitCode:    1,
				logLevel:    "error",
				message:     fmt.Sprintf("Flag '--interactive/-i' cannot be used with flags '--crate/-c' and '--target/-t' or '--all/-a'"),
				indentLevel: program.indentLevel,
			}
		}

		availableCrates, response := getUserCrates(program)
		handleFunctionResponse(response, true)

		availableCratesStrings := make([]string, len(availableCrates))
		for i, element := range availableCrates {
			availableCratesStrings[i] = element.name
		}

		var selectedCrateIndex int
		promptCrate := &survey.Select{
			Message: "Select the crate",
			Options: availableCratesStrings,
			Default: selectedCrateIndex,
		}
		err := survey.AskOne(promptCrate, &selectedCrateIndex, survey.WithPageSize(10))
		if err != nil {
			if err.Error() == "interrupt" {
				return Crate{}, []Target{}, functionResponse{
					exitCode:    1,
					message:     "Operation cancelled by user",
					logLevel:    "error",
					indentLevel: program.indentLevel,
				}
			}
		}

		selectedCrate = availableCrates[selectedCrateIndex]

		availableTargets, response := getCrateTargets(selectedCrate, program)
		handleFunctionResponse(response, true)

		availableTargetsStrings := make([]string, len(availableTargets))
		for i, element := range availableTargets {
			availableTargetsStrings[i] = element.name
		}

		var selectedTargetsIndices []int
		promptTarget := &survey.MultiSelect{
			Message: "Select the target(s)",
			Options: availableTargetsStrings,
			Default: selectedTargetsIndices,
		}
		err = survey.AskOne(promptTarget, &selectedTargetsIndices, survey.WithPageSize(10))
		if err != nil {
			if err.Error() == "interrupt" {
				return Crate{}, []Target{}, functionResponse{
					exitCode:    1,
					message:     "Operation cancelled by user",
					logLevel:    "error",
					indentLevel: program.indentLevel,
				}
			}
		}

		selectedTargets := make([]Target, len(selectedTargetsIndices))
		for i, index := range selectedTargetsIndices {
			selectedTargets[i] = availableTargets[index]
		}

		if len(selectedTargets) == 0 {
			return selectedCrate, selectedTargets, functionResponse{
				exitCode:    1,
				logLevel:    "attention",
				message:     fmt.Sprintf("No targets were selected"),
				indentLevel: program.indentLevel,
			}
		}

		return selectedCrate, selectedTargets, functionResponse{exitCode: 0}
	} else {
		if allTargets != false && targetNames != nil {
			return Crate{}, []Target{}, functionResponse{
				exitCode:    1,
				logLevel:    "error",
				message:     fmt.Sprintf("Conflicting flags: both '--target/-t' and '--all/-a' flags cannot be specified at the same time"),
				indentLevel: program.indentLevel,
			}
		}

		if allTargets == false && targetNames == nil {
			return Crate{}, []Target{}, functionResponse{
				exitCode:    1,
				logLevel:    "error",
				message:     fmt.Sprintf("Missing required flag: '--interactive/-i' or '--target/-t' or '--all/-a' flag must be specified"),
				indentLevel: program.indentLevel,
			}
		}

		crate := generateCrateObj(crateName, program)

		response := verifyCrateDirectory(crate, program)
		if response.exitCode != 0 {
			return Crate{}, []Target{}, functionResponse{
				exitCode:    1,
				message:     fmt.Sprintf("Crate '%s' not found", crate.name),
				logLevel:    "error",
				indentLevel: program.indentLevel,
			}
		}

		if allTargets == true {
			selectedTargets, response := getCrateTargets(crate, program)
			handleFunctionResponse(response, true)

			return crate, selectedTargets, functionResponse{exitCode: 0}
		} else {
			selectedTargets = make([]Target, len(targetNames))
			for i, element := range targetNames {
				selectedTargets[i] = generateTargetObj(crate.name, element, program)
			}

			// Verify targets
			response := verifyTargetsDirectories(selectedTargets, program)
			handleFunctionResponse(response, true)

			return crate, selectedTargets, functionResponse{exitCode: 0}
		}
	}
}

func displayTargetTag(msg string, target Target) string {
	return fmt.Sprintf(msg) + fmt.Sprintf(" (") + salmonPink.Sprintf(target.crate.name) + fmt.Sprintf("/") + green.Sprintf(target.name) + fmt.Sprintf(")")
}

func getCrateTargets(crate Crate, program Program) ([]Target, functionResponse) {
	targetNames, err := ioutil.ReadDir(crate.targetsDir)
	if err != nil {
		return []Target{}, functionResponse{
			exitCode:    1,
			message:     fmt.Sprintf("Failed to read the crate's targets directory -> " + err.Error()),
			logLevel:    "error",
			indentLevel: program.indentLevel,
		}
	}
	targetNames = filterHiddenFilesAndDirectories(targetNames)

	// Generate an array of targets
	targets := make([]Target, len(targetNames))
	for i, element := range targetNames {
		targets[i] = generateTargetObj(crate.name, element.Name(), program)
	}

	if len(targets) == 0 {
		return targets, functionResponse{
			exitCode:    1,
			logLevel:    "attention",
			message:     "No targets found",
			indentLevel: program.indentLevel,
		}
	}

	return targets, functionResponse{
		exitCode: 0,
	}
}

func generateTargetObj(crate string, target string, program Program) Target {
	defaultTargetEnv := map[string]string{
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
		"TARGET_NAME":                target,
		"TARGET_DIR":                 program.userCratesDir + "/" + crate + "/targets" + "/" + target,
		"TARGET_HOOKS_DIR":           program.userCratesDir + "/" + crate + "/targets" + "/" + target + "/hooks",
		"TARGET_TEMP_DIR":            program.userCratesDir + "/" + crate + "/targets" + "/" + target + "/.tmp",
	}

	return Target{
		crate:        generateCrateObj(crate, program),
		name:         target,
		path:         program.userCratesDir + "/" + crate + "/targets" + "/" + target,
		hooksDir:     program.userCratesDir + "/" + crate + "/targets" + "/" + target + "/hooks",
		tempDir:      program.userCratesDir + "/" + crate + "/targets" + "/" + target + "/.tmp",
		disabledPath: program.userCratesDir + "/" + crate + "/targets" + "/" + target + "/disabled",
		environment:  defaultTargetEnv,
	}
}

func isTargetDisabled(target Target, program Program) (bool, functionResponse) {
	if _, err := os.Stat(target.disabledPath); os.IsNotExist(err) {
		return false, functionResponse{
			exitCode:    0,
			indentLevel: program.indentLevel + 1,
		}
	} else if err != nil {
		return false, functionResponse{
			exitCode:    1,
			message:     "Failed to verify if target is disabled -> " + err.Error(),
			logLevel:    "error",
			indentLevel: program.indentLevel + 1,
		}
	} else {
		return true, functionResponse{
			exitCode: 0,
		}
	}
}

func enableTarget(target Target, program Program) functionResponse {
	if _, err := os.Stat(target.disabledPath); os.IsNotExist(err) {
		return functionResponse{
			exitCode:    0,
			message:     "Target already enabled",
			logLevel:    "attention",
			indentLevel: program.indentLevel + 1,
		}
	} else if err != nil {
		return functionResponse{
			exitCode:    1,
			message:     "Failed to verify if target is disabled -> " + err.Error(),
			logLevel:    "error",
			indentLevel: program.indentLevel + 1,
		}
	} else {
		err = os.Remove(target.disabledPath)
		if err != nil {
			return functionResponse{
				exitCode:    1,
				message:     "Failed to remove disabled lock file -> " + err.Error(),
				logLevel:    "error",
				indentLevel: program.indentLevel + 1,
			}
		} else {
			return functionResponse{
				exitCode:    0,
				message:     "Finished",
				logLevel:    "success",
				indentLevel: program.indentLevel + 1,
			}
		}
	}
}

func targetsEnable(crate Crate, targets []Target, program Program) functionResponse {
	for _, target := range targets {
		space()
		showInfoSectionTitle(displayTargetTag("Enabling", target), program.indentLevel)

		response := enableTarget(target, program)
		if response.exitCode != 0 {
			return response
		}
		handleFunctionResponse(response, false)
	}

	return functionResponse{
		exitCode: 0,
	}
}

func disableTarget(target Target, program Program) functionResponse {
	if _, err := os.Stat(target.disabledPath); os.IsNotExist(err) {
		file, err := os.Create(target.disabledPath)
		defer file.Close()
		if err != nil {
			return functionResponse{
				exitCode:    1,
				message:     "Failed to create disabled lock file -> " + err.Error(),
				logLevel:    "error",
				indentLevel: program.indentLevel + 1,
			}
		} else {
			return functionResponse{
				exitCode:    0,
				message:     "Finished",
				logLevel:    "success",
				indentLevel: program.indentLevel + 1,
			}
		}
	} else if err != nil {
		return functionResponse{
			exitCode:    1,
			message:     "Failed to verify if target is enabled -> " + err.Error(),
			logLevel:    "error",
			indentLevel: program.indentLevel + 1,
		}
	} else {
		return functionResponse{
			exitCode:    0,
			message:     "Target already disabled",
			logLevel:    "attention",
			indentLevel: program.indentLevel + 1,
		}
	}
}

func targetsDisable(crate Crate, targets []Target, program Program) functionResponse {
	for _, target := range targets {
		space()
		showInfoSectionTitle(displayTargetTag("Disabling", target), program.indentLevel)

		response := disableTarget(target, program)
		if response.exitCode != 0 {
			return response
		}
		handleFunctionResponse(response, false)
	}

	return functionResponse{
		exitCode: 0,
	}
}

func verifyTargetsDirectories(targets []Target, program Program) functionResponse {
	var failedTargets []Target

	for _, target := range targets {
		response := verifyTargetDirectory(target, program)
		if response.exitCode != 0 {
			failedTargets = append(failedTargets, target)
		}
	}

	var response functionResponse
	if len(failedTargets) > 0 {
		message := "The following target(s) was/were not found:\n"
		for _, target := range failedTargets {
			message = message + fmt.Sprintf("\n    - %s", target.name)
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

func verifyTargetDirectory(target Target, program Program) functionResponse {
	if _, err := os.Stat(target.path); os.IsNotExist(err) {
		return functionResponse{
			exitCode: 1,
		}
	}
	return functionResponse{
		exitCode: 0,
	}
}

func removeTargetTempDirectory(target Target, notRemoveTempDir bool, program Program) {
	var response functionResponse

	showInfoSectionTitle(fmt.Sprintf("Removing temporary directory"), program.indentLevel)

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

	if _, err := os.Stat(target.tempDir); os.IsNotExist(err) {
		response = functionResponse{
			exitCode:    0,
			logLevel:    "attention",
			message:     fmt.Sprintf("Temporary directory not found"),
			indentLevel: program.indentLevel + 1,
		}
	} else {
		err = os.RemoveAll(target.tempDir)
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

func setupTargetTempDirectory(target Target, notCreateTempDir bool, program Program) {
	var response functionResponse

	showInfoSectionTitle(lightGray.Sprintf("Setting up temporary directory"), program.indentLevel)

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

	if _, err := os.Stat(target.tempDir); err == nil {
		response = functionResponse{
			exitCode:    0,
			message:     "Temporary directory already exists. Recreating it...",
			logLevel:    "attention",
			indentLevel: program.indentLevel + 2,
		}
		handleFunctionResponse(response, false)

		err = os.RemoveAll(target.tempDir)
		if err != nil {
			response = functionResponse{
				exitCode:    1,
				message:     fmt.Sprintf("Failed to recreate temporary directory -> '%v'", err.Error()),
				indentLevel: program.indentLevel + 3,
			}

			handleFunctionResponse(response, true)
		}
	}

	perm := os.FileMode(0755)
	err := os.Mkdir(target.tempDir, perm)
	if err != nil {
		response = functionResponse{
			exitCode:    1,
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

func targetsCreate(program Program) functionResponse {
	var selectedCrate Crate
	var targetName string
	var targetTemplate string

	// Ask for a parent crate

	availableCrates, response := getUserCrates(program)
	handleFunctionResponse(response, true)

	availableCratesStrings := make([]string, len(availableCrates))
	for i, element := range availableCrates {
		availableCratesStrings[i] = element.name
	}

	var selectedIndex int
	promptCrate := &survey.Select{
		Message: "Select a crate",
		Options: availableCratesStrings,
		Default: selectedIndex,
	}
	err := survey.AskOne(promptCrate, &selectedIndex, survey.WithPageSize(10))
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

	selectedCrate = availableCrates[selectedIndex]

	//
	////
	//

	// Ask for a target name

	promptTargetName := &survey.Input{
		Message: "Target name:",
	}
	err = survey.AskOne(promptTargetName, &targetName, survey.WithValidator(survey.MinLength(2)))
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

	// Generate a target object
	target := generateTargetObj(selectedCrate.name, targetName, program)

	// Ask for a target template
	availableTemplates, err := ioutil.ReadDir(program.userTargetsTemplatesDir)
	if err != nil {
		return functionResponse{
			exitCode:    1,
			message:     fmt.Sprintf("Failed to read the user's target templates directory -> " + err.Error()),
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

	promptTargetTemplate := &survey.Select{
		Message: "Target template:",
		Options: availableTemplatesStrings,
	}
	err = survey.AskOne(promptTargetTemplate, &targetTemplate)
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
	var targetTemplateDir string
	scratchTemplate := false
	if targetTemplate == "scratch" {
		scratchTemplate = true
	} else {
		targetTemplateDir = program.userTargetsTemplatesDir + "/" + targetTemplate
	}

	showInfoSectionTitle(displayTargetTag("Creating target", target), program.indentLevel)

	// Verify if target already exists
	response = verifyTargetDirectory(target, program)
	if response.exitCode == 0 {
		return functionResponse{
			exitCode:    1,
			message:     fmt.Sprintf("Target '%s' already exists", target.name),
			logLevel:    "attention",
			indentLevel: program.indentLevel + 1,
		}
	}

	// Create target directory
	space()
	showInfoSectionTitle("Creating target directory", program.indentLevel+1)
	err = os.Mkdir(target.path, 0755)
	if err != nil {
		return functionResponse{
			exitCode:    1,
			message:     fmt.Sprintf("Failed to create target directory -> " + err.Error()),
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

	// Copy template to target directory
	if scratchTemplate == false {
		space()
		showInfoSectionTitle("Copying template to target directory", program.indentLevel+1)
		copyOptions := copy.Options{
			PreserveTimes: true,
			PreserveOwner: true,
		}

		err = copy.Copy(targetTemplateDir, target.path, copyOptions)
		if err != nil {
			// Remove target directory
			_ = os.RemoveAll(target.path)

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
	if _, err := os.Stat(target.hooksDir + "/post_create"); err == nil {
		_, response := runHook(target.hooksDir+"/post_create", target.environment, true, true, true, true, true, incrementProgramIndentLevel(program, 1))

		if response.exitCode != 0 {
			handleFunctionResponse(response, false)

			space()

			// Remove target directory
			showInfoSectionTitle(lightGray.Sprintf("Removing target directory"), program.indentLevel+1)
			err = os.RemoveAll(target.path)

			if err != nil {
				response = functionResponse{
					exitCode:    response.exitCode,
					message:     "Failed to remove target directory",
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

func targetsRm(crate Crate, targets []Target, program Program) functionResponse {
	var response functionResponse

	for index, target := range targets {
		space()

		orange.Println(fmt.Sprintf("(%v/%v)", index+1, len(targets)))
		showInfoSectionTitle(displayTargetTag("Removing", target), program.indentLevel)

		showInfoSectionTitle(lightGray.Sprintf("Running ")+orange.Sprintf("pre_rm")+lightGray.Sprintf(" hook"), program.indentLevel+1)

		if _, err := os.Stat(target.hooksDir + "/pre_rm"); os.IsNotExist(err) {
			response := functionResponse{
				exitCode:    0,
				message:     "Hook not found",
				logLevel:    "attention",
				indentLevel: program.indentLevel + 2,
			}
			handleFunctionResponse(response, false)
		} else {
			program = incrementProgramIndentLevel(program, 1)

			_, response := runHook(target.hooksDir+"/pre_rm", target.environment, true, true, true, true, true, program)
			response.indentLevel = program.indentLevel + 2
			handleFunctionResponse(response, true)
		}

		space()

		err := os.RemoveAll(target.path)
		if err != nil {
			response = functionResponse{
				exitCode:    1,
				logLevel:    "error",
				message:     "Failed to remove target -> " + err.Error(),
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

func targetsLs(crates []Crate, program Program) functionResponse {
	for index, crate := range crates {
		space()

		orange.Println(fmt.Sprintf("(%v/%v)", index+1, len(crates)))
		showInfoSectionTitle(displayCrateTag("Listing targets", crate), program.indentLevel)

		targets, response := getCrateTargets(crate, program)
		if response.exitCode != 0 {
			response.indentLevel = program.indentLevel + 1
			return response
		}

		space()

		for _, target := range targets {
			// Get optional description (if hook exists)
			targetDescription, response := runHook(target.hooksDir+"/ls", target.environment, false, false, false, false, false, program)
			targetDescriptionString := targetDescription.Output

			if response.exitCode == 0 {
				if len(targetDescriptionString) > 0 {
					showText(fmt.Sprintf("- %s (%s)", target.name, blue.Sprintf(targetDescriptionString)), program.indentLevel+1)
				} else {
					showText(fmt.Sprintf("- %s", target.name), program.indentLevel+1)
				}
			} else {
				showText(fmt.Sprintf("- %s", target.name), program.indentLevel+1)
			}
		}
	}

	return functionResponse{
		exitCode: 0,
	}
}

func targetsEdit(crate Crate, targets []Target, program Program) functionResponse {
	hook := "edit"
	response := targetsRunHooks(crate, targets, []string{hook}, false, false, false, false, false, program)

	return response
}

func targetsView(crate Crate, targets []Target, program Program) functionResponse {
	hook := "view"
	response := targetsRunHooks(crate, targets, []string{hook}, false, false, false, false, false, program)

	return response
}

func targetsSync(crate Crate, targets []Target, program Program) functionResponse {
	var response functionResponse

	setupCrateTempDirectory(crate, true, false, program)

	// Run pre_transaction hook for crate (if any)
	space()
	showInfoSectionTitle(lightGray.Sprintf("Running pre_transaction hook")+lightGray.Sprintf(" (")+salmonPink.Sprintf(crate.name)+lightGray.Sprintf(")"), program.indentLevel)
	if _, err := os.Stat(crate.hooksDir + "/pre_transaction"); os.IsNotExist(err) {
		response = functionResponse{
			exitCode:    0,
			message:     "Hook not found",
			logLevel:    "attention",
			indentLevel: program.indentLevel + 1,
		}
		handleFunctionResponse(response, false)
	} else {
		_, response := runHook(crate.hooksDir+"/pre_transaction", crate.environment, true, true, true, true, true, program)
		response.indentLevel = program.indentLevel + 1

		handleFunctionResponse(response, false)

		if response.exitCode != 0 {
			space()
			removeCrateTempDirectory(crate, true, false, program)

			space()

			finishProgram(response.exitCode)
		}
	}

	space()

	response = func(targets []Target, program Program) functionResponse {
		for _, target := range targets {
			space()
			space()

			showInfoSectionTitle(lightGray.Sprintf("Syncing")+lightGray.Sprintf(" (")+salmonPink.Sprintf(target.crate.name)+"/"+green.Sprintf(target.name)+lightGray.Sprintf(")"), program.indentLevel)

			isTargetDisabled, response := isTargetDisabled(target, program)
			if response.exitCode != 0 {
				response.indentLevel = program.indentLevel + 1

				return response
			}

			if isTargetDisabled == true {
				response = functionResponse{
					exitCode:    0,
					message:     "Target is disabled",
					logLevel:    "attention",
					indentLevel: program.indentLevel + 1,
				}
				handleFunctionResponse(response, false)
			} else {
				program = incrementProgramIndentLevel(program, 1)

				space()
				setupTargetTempDirectory(target, false, program)

				space()
				space()
				showInfoSectionTitle(lightGray.Sprintf("Running pre_transaction hook"), program.indentLevel)
				_, response = runHook(target.hooksDir+"/pre_transaction", target.environment, true, true, true, true, true, program)
				response.indentLevel = program.indentLevel + 1
				handleFunctionResponse(response, false)

				space()
				space()
				showInfoSectionTitle(lightGray.Sprintf("Running sync hook"), program.indentLevel)
				_, response = runHook(target.hooksDir+"/sync", target.environment, true, true, true, true, true, program)
				response.indentLevel = program.indentLevel + 1

				if response.exitCode != 0 {
					response.indentLevel = program.indentLevel + 1
					handleFunctionResponse(response, false)

					space()
					space()
					removeTargetTempDirectory(target, false, program)

					return response
				}

				space()
				space()
				showInfoSectionTitle(lightGray.Sprintf("Running post_transaction hook"), program.indentLevel)
				_, response = runHook(target.hooksDir+"/post_transaction", target.environment, true, true, true, true, true, program)
				response.indentLevel = program.indentLevel + 1
				handleFunctionResponse(response, false)

				space()
				space()
				removeTargetTempDirectory(target, false, program)

				program = decrementProgramIndentLevel(program, 1)
			}
		}

		return functionResponse{
			exitCode: 0,
		}
	}(targets, program)

	if response.exitCode != 0 {
		space()
		removeCrateTempDirectory(crate, true, false, program)

		space()

		finishProgram(response.exitCode)
	}

	// Run post_transaction hook for crate (if any)
	space()
	space()
	space()

	showText(lightGray.Sprintf("Running post_transaction hook")+lightGray.Sprintf(" (")+salmonPink.Sprintf(crate.name+lightGray.Sprintf(")")), program.indentLevel)
	if _, err := os.Stat(crate.hooksDir + "/post_transaction"); os.IsNotExist(err) {
		response = functionResponse{
			exitCode:    0,
			message:     "Hook not found",
			logLevel:    "attention",
			indentLevel: program.indentLevel + 1,
		}
		handleFunctionResponse(response, false)
	} else {
		_, response := runHook(crate.hooksDir+"/post_transaction", crate.environment, true, true, true, true, true, program)
		response.indentLevel = program.indentLevel + 1
		handleFunctionResponse(response, true)
	}

	space()
	removeCrateTempDirectory(crate, true, false, program)

	return functionResponse{
		exitCode: 0,
	}
}

func targetsRunHooks(crate Crate, targets []Target, hooks []string, notCreateTempDir bool, notRemoveTempDir bool, notPrintOutput bool, notPrintEntryCmd bool, notPrintAlerts bool, program Program) functionResponse {
	var response functionResponse

	setupCrateTempDirectory(crate, true, notCreateTempDir, program)

	space()

	for index, target := range targets {
		orange.Println(fmt.Sprintf("(%v/%v)", index+1, len(targets)))

		showInfoSectionTitle(displayTargetTag("Running hook(s)", target), program.indentLevel)

		program = incrementProgramIndentLevel(program, 1)

		space()

		setupTargetTempDirectory(target, notCreateTempDir, program)

		response = func(crate Crate, target Target, hooks []string, notRemoveTempDir bool, notPrintOutput bool, notPrintEntryCmd bool, program Program) functionResponse {
			for _, hook := range hooks {
				space()
				space()

				showInfoSectionTitle(lightGray.Sprintf("Running ")+orange.Sprintf(hook)+lightGray.Sprintf(" hook"), program.indentLevel)

				// Run hook
				if _, err := os.Stat(target.hooksDir + "/" + hook); os.IsNotExist(err) {
					response = functionResponse{
						exitCode:    1,
						message:     fmt.Sprintf("No '%v' hook found", hook),
						logLevel:    "error",
						indentLevel: program.indentLevel + 1,
					}
					handleFunctionResponse(response, false)
				} else {
					_, hookResponse := runHook(target.hooksDir+"/"+hook, target.environment, !notPrintOutput, true, !notPrintEntryCmd, true, !notPrintAlerts, program)

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
		}(crate, target, hooks, notRemoveTempDir, notPrintOutput, notPrintEntryCmd, program)

		if response.exitCode != 0 {
			handleFunctionResponse(response, false)

			space()
			space()

			removeTargetTempDirectory(target, notRemoveTempDir, program)

			space()

			removeCrateTempDirectory(crate, true, notRemoveTempDir, decrementProgramIndentLevel(program, 1))

			space()
			finishProgram(response.exitCode)
		} else {
			handleFunctionResponse(response, false)
		}

		space()
		space()

		removeTargetTempDirectory(target, notRemoveTempDir, program)

		space()

		removeCrateTempDirectory(crate, true, notRemoveTempDir, decrementProgramIndentLevel(program, 1))

		program = decrementProgramIndentLevel(program, 1)
	}

	return functionResponse{
		exitCode: 0,
	}
}

func targetsHooksLs(crate Crate, targets []Target, program Program) functionResponse {
	for index, target := range targets {
		space()

		orange.Println(fmt.Sprintf("(%v/%v)", index+1, len(targets)))
		showInfoSectionTitle(displayTargetTag("Listing hooks", target), program.indentLevel)

		hooks, err := ioutil.ReadDir(target.hooksDir)
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
			customEntryFilePath := target.hooksDir + "/" + element.Name() + ".entry"
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
