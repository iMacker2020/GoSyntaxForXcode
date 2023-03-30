# GoSyntaxForXcode
Adds Go syntax highlighting for Xcode

This project adds syntax highlighting to Xcode providing color to various keywords of Go. 

![Xcode Go syntax highlighting](https://user-images.githubusercontent.com/69224955/225790906-737b0400-b811-4a62-a686-61822ea8d835.png)

The included installer is a fat-binary that will do all the work of installing the files for you. It contains binaries for both x86_64 and ARM64 (Sorry i386 and PowerPC users). To use it simply double-click on it in the Finder. Once it is done restart Xcode (if it was open) and open a .go file. 

Notes:
This installer has been tested on:
- Mac OS 10.8 (Xcode 5.1.1)
- Mac OS 10.9 (Xcode 6)
- Mac OS 10.12 (Xcode 9.2)
- Mac OS 10.14 (Xcode 11.3 - does not work due to a bug with Xcode) 
- Mac OS 12.3 (Xcode 13.3.1)

Manual Installation:
- Copy the file Go.ideplugin to ~/Library/Developer/Xcode/Plug-ins/Go.ideplugin.
- Copy the file Go.xclangspec to ~/Library/Developer/Xcode/Specifications/Go.xclangspec.
- Run this command: defaults read /Applications/Xcode.app/Contents/Info DVTPlugInCompatibilityUUID
- Copy the UUID.
- Open the file ~/Library/Developer/Xcode/Plug-ins/Go.ideplugin/Contents/Info.plist.
- Paste the UUID to under DVTPlugInCompatibilityUUIDs.

Uninstall:
- Run the included installer with the option of "uninstall": ./installer uninstall
