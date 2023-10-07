package main

//
//// IMPORTS
//

import (
	// Modules in GOROOT
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	// External modules
	color "github.com/gookit/color"
)

//
//// PROGRAM CONFIGURATION
//

type Program struct {
	name                    string
	nameAscii               string
	version                 string
	exec                    string
	shortDescription        string
	longDescription         string
	defaultShell            string
	userDataDir             string
	userCratesDir           string
	userTemplatesDir        string
	userTargetsTemplatesDir string
	userCratesTemplatesDir  string
	indentLevel             int
}

func getDefaultShellAbsolutePath(shellName string) string {
	// Get shell absolute path using `which`
	cmd := exec.Command("which", shellName)

	output, err := cmd.CombinedOutput()
	outputString := string(output)
	outputString = strings.TrimLeft(outputString, "\n")
	outputString = strings.TrimRight(outputString, "\n")

	if err != nil {
		showError(fmt.Sprintf("Failed to get absolute path to default shell '%v' -> %v", shellName, outputString), 0)
		finishProgram(1)
	}

	return outputString
}

func initializeDefaultProgram(customUserDataDir string) Program {
	// PROGRAM NAME
	programName := "synctropy"

	// PROGRAM NAME ASCII
	programNameAscii := `
 ____                   _
/ ___| _   _ _ __   ___| |_ _ __ ___  _ __  _   _
\___ \| | | | '_ \ / __| __| '__/ _ \| '_ \| | | |
 ___) | |_| | | | | (__| |_| | | (_) | |_) | |_| |
|____/ \__, |_| |_|\___|\__|_|  \___/| .__/ \__, |
       |___/                         |_|    |___/
`

	// PROGRAM VERSION
	programVersion := "0.2.1"

	// PROGRAM EXEC
	programExec := os.Args[0] // Path for program executable

	// DESCRIPTIONS (SHORT AND LONG)
	programShortDescription := "A wrapper for management and syncing of crates via syncing utilities like unison and rsync using hooks, with template support."
	programLongDescription := fmt.Sprintf("%v is a wrapper designed for syncing and managing crate configurations using utilities\nlike unison and rsync via hooks. With a user-friendly structure, users can effortlessly create\nand manage crates that are tailored for specific programs. They can easily set up targets within\nthese crate configurations, allowing for efficient synchronization of data. The program also\noffers the convenience of using templates when creating new crates and targets, ensuring a\nconsistent and streamlined experience. \n\nBearing a name that fuses %v and the scientific concept %v - signifying the shift from\ndisarray to structure, %v aims to manage the mix of your various files and turn them into\na smoothly synchronized collection. It's about evolving from entropy to syntropy, converting the\ndisordered into the organized.", color.HEX("#55ff7f").Sprintf(programName), color.HEX("#ffaa7f").Sprintf("sync"), color.HEX("#ffaa7f").Sprintf("syntropy"), programName)

	// DEFAULT SHELL
	programDefaultShellName := "sh" // Should work on all Unix systems (Linux, Android, ...)
	programDefaultShellPath := getDefaultShellAbsolutePath(programDefaultShellName)

	// USER DIRECTORIES
	var userDataDir string
	if customUserDataDir != "" {
		userDataDir = customUserDataDir
	} else {
		userDataDir = getCurrentUserHomeDir(Program{indentLevel: 0}) + "/" + programName
	}

	userCratesDir := userDataDir + "/crates"
	userTemplatesDir := userDataDir + "/templates"
	userTargetsTemplatesDir := userTemplatesDir + "/targets"
	userCratesTemplatesDir := userTemplatesDir + "/crates"

	// INDENT LEVEL
	indentLevel := 0

	return Program{
		name:                    programName,
		nameAscii:               programNameAscii,
		version:                 programVersion,
		exec:                    programExec,
		shortDescription:        programShortDescription,
		longDescription:         programLongDescription,
		defaultShell:            programDefaultShellPath,
		userDataDir:             userDataDir,
		userCratesDir:           userCratesDir,
		userTemplatesDir:        userTemplatesDir,
		userTargetsTemplatesDir: userTargetsTemplatesDir,
		userCratesTemplatesDir:  userCratesTemplatesDir,
		indentLevel:             indentLevel,
	}
}

func getRootDirectory() string {
	// Check if the "PREFIX" environment variable is set
	prefix := os.Getenv("PREFIX")
	if prefix != "" {
		// Verify if the given prefix string ends with the directory suffix "/usr".
		// If it does, the program will proceed to strip it from the prefix.
		if strings.HasSuffix(prefix, "/usr") {
			return strings.TrimSuffix(prefix, "/usr")
		}
		return prefix
	} else {
		return "/"
	}
}

func displayProgramInfo(program Program) {
	showText(program.nameAscii, program.indentLevel)

	showText("Version: "+green.Sprintf(program.version), program.indentLevel)

	space()

	showText("Running on "+lightCopper.Sprintf(runtime.GOOS+"/"+runtime.GOARCH)+". Built with "+runtime.Version()+" using "+runtime.Compiler+" as compiler.", program.indentLevel)

	space()

	showText("This program is licensed under the GNU General Public License v3.0 (GPL-3.0).\nPlease refer to the LICENSE file for more information.", program.indentLevel)
}
