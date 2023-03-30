// File: UUIDsetter.go
// Description: Installs the files Go.ideplugin and Go.xclangspec.

package main

import "fmt"
import "os/exec"
import "os"
import _"os/user"
import "strings"
import "io/ioutil"
import "time"


func main() {
    defer checkForError()
	printBanner()
	checkArguments()
	checkFiles()
	copyLangFile()
	copyPluginFile()
	checkPluginUUIDs()
	fmt.Println("Installation completed. Please restart Xcode to load the plug-in.\n")
}

// Handles any arguments sent to this program
func checkArguments() {
	if len(os.Args) == 2 && os.Args[1] == "uninstall" {
		uninstallAllFiles()
		os.Exit(0)
	} else if len(os.Args) > 1 {
		fmt.Println("Usage: installer [uninstall]\n")
		os.Exit(0)
	}
}

// Prints a little banner
func printBanner() {
	outputStr :=
	`
    _______________________________
   /\                              \
   \_| Go syntax support for Xcode |
     | v1.1 - March 30, 2023       |
     |   __________________________|_
     \_/____________________________/`
	fmt.Println(outputStr + "\n\n")
}

// Uninstall files if the user specifies "uninstall" as the argument to this program
func uninstallAllFiles() {
	var err error

	pluginPath := os.Getenv("HOME") + "/Library/Developer/Xcode/Plug-ins/Go.ideplugin"
	specPath := os.Getenv("HOME") + "/Library/Developer/Xcode/Specifications/Go.xclangspec"
	destFolder := fmt.Sprintf("%s - %s", "Xcode Files", time.Now())
	destFolder = destFolder[0 : strings.LastIndex(destFolder, ":") + 5] // add a few digits of the seconds to the name
	
	// Create the Xcode folder in the trash
    err = os.Chdir(os.Getenv("HOME") + "/.Trash/")
	check(err, "Could not change current directory to Trash")
	err = os.Mkdir(destFolder, 0755)
	check(err, "Could not make xcode folder in Trash")
    
	destFolder = os.Getenv("HOME") + "/.Trash/" + destFolder
	
	files := []string{pluginPath, specPath}
	for _,path := range(files) {
        newPath := destFolder + path[strings.LastIndex(path, "/") : ] // get the last path component (the file)
        fmt.Println("Moving file", path, "to", newPath)
        err = os.Rename(path, newPath)
		if err != nil {
			fmt.Println("\aCould not move file", path, "\nError:", err)
		}
	}
	fmt.Println("Files uninstalled. Please restart Xcode for changes to take effect.\n")
}

// Check that all the files are with the installer
func checkFiles() {
	files := []string{"Go.xclangspec", "Go.ideplugin"}
	for _, file := range(files) {
		path, err := os.Executable()
		check(err, "Could not find installer's path")
		index := strings.LastIndex(path, "/")
		path = path[0 : index + 1] // remove this program's name from the path
		path = path + file
		openFile, err := os.Open(path)
		if err != nil {
			message := "Failed to find file " + file + ".\nPlease ensure it is in the same folder as this installer."
			panic(message)
		}
		openFile.Close()
	}
}

// Copy the language file
func copyLangFile() {
	var err error
	var c *exec.Cmd

	// create the folder path
	err = os.Chdir(os.Getenv("HOME"))
	check(err, "Could not change directories to home")
	err = os.MkdirAll("Library/Developer/Xcode/Specifications/", 0755)
	check(err, "Could not create folders for specification file")
	
	// create the source path
	path, err := os.Executable()
	check(err, "Could not find installer's path")
	index := strings.LastIndex(path, "/")
	path = path[0 : index + 1] // remove this program's name from the path
	source := path + "Go.xclangspec"
	
	// copy the file
	destination := os.Getenv("HOME") + "/Library/Developer/Xcode/Specifications"
	fmt.Println("Executing: cp", source, destination)
	c = exec.Command("cp", source, destination)
	output, err := c.CombinedOutput()
	check(err, "Could not copy file Go.xclangspec:" + string(output))
}

// Copy the plug-in file
func copyPluginFile() {
	var err error
	var c *exec.Cmd

	// create the folder path
	err = os.MkdirAll("Library/Developer/Xcode/Plug-ins/", 0755)
	check(err, "Could not create plug-in folder path")
	
	// create the source path
	path, err := os.Executable()
	check(err, "Could not find installer's path")
	index := strings.LastIndex(path, "/")
	path = path[0 : index + 1] // remove this program's name from the path
	source := path + "Go.ideplugin"
	
	// copy the file
	destination := os.Getenv("HOME") + "/Library/Developer/Xcode/Plug-ins/"
	c = exec.Command("cp", "-r", source, destination)
	fmt.Println("Executing: cp", "-r", source, destination)
	output, err := c.CombinedOutput()
	check(err, "Could not copy file Go.ideplugin:" + string(output))
}

// Query the system for Xcode's UUID for plugins and return the string
func getXcodeUUID() string {
	cmd := exec.Command("defaults", "read", "/Applications/Xcode.app/Contents/Info", "DVTPlugInCompatibilityUUID")
	value, err := cmd.Output()
	check(err, "Failed to run defaults command")
	
	xcodeUUID := string(value[:len(value) - 1]) // remove newline
	fmt.Println("Xcode UUID:", xcodeUUID)
	return xcodeUUID
}

// Opens the plugin's Info.plist file and adds Xcode's UUID if needed
func checkPluginUUIDs() {
	// Get file's path
	filePath := "/Library/Developer/Xcode/Plug-ins/Go.ideplugin/Contents/Info.plist"
    homeFolder := os.Getenv("HOME")

	// Read the Info.plist file
	filePath = homeFolder + filePath
	fileData, err := ioutil.ReadFile(filePath)
	check(err, "Failed to open Info.plist file")
	
	// Check that the text DVTPlugInCompatibilityUUIDs is found
	fileStr := string(fileData) // convert byte array into string
	pos := strings.Index(fileStr, "DVTPlugInCompatibilityUUIDs")
	if pos == -1 {
		panic("Failed to find DVTPlugInCompatibilityUUIDs string in Info.plist file")
	}
	
	xcodeUUID := getXcodeUUID()
	
	// If Xcode's UUID not found, add it in
	if strings.Index(fileStr, xcodeUUID) == -1 {
		startPos := pos + len("DVTPlugInCompatibilityUUIDs</key> ") + len("<array> ") + 1
		file, err := os.OpenFile(filePath, os.O_RDWR, 0777)
		check(err, "Failed to open Info.plist file")
		
		file.Seek(int64(startPos), 0)
		_,err = fmt.Fprintf(file, "\n\t\t<string>" + xcodeUUID + "</string>")
		check(err, "Failed to add UUID to file")
		
		// write the rest of the data to the file
		_, err = fmt.Fprintf(file, "%s", fileData[startPos :])
		check(err, "Failed to append data to file")
		
		err = file.Close()
		check(err, "Failed to close Info.plist file")
		
		fmt.Println("Added Xcode UUID to plugin")
	} else {
		fmt.Println("Xcode UUID already in plugin")
	}
}

// handles error handling
func check(err error, message string) {
	if err != nil {
		errStr := fmt.Sprintf(message + ": %s", err)
		panic(errStr)
	}
}

// For debugging the installer
func debugMessage() {
    speak("UUID check done")
}

// This is like the catch block in an exception
func checkForError() {
    err := recover()
    if err != nil {
        message := fmt.Sprintf("Error: %s", err)
        //speak(message)
        displayDialog(message)
    }
}

// Displays a dialog
func displayDialog(message string) {
    c := exec.Command("osascript", "-e",
    `tell application (path to frontmost application as text) to display dialog "` + message + `" buttons {"OK"} with icon stop`)
    err := c.Run()
    if err != nil {
        fmt.Println("\aError: can't display dialog.\nMessage:", message)
    }
}

// Speaks a message
func speak(message string) {
    c := exec.Command("say", message)
    err := c.Run()
    if err != nil {
        println("\aError: can't talk.\nMessage:", message, "\n")
    }
}
