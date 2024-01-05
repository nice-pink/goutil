// This is a collection of filesystem functions.
// Files and directories can be read, written, manipulated.
// This collection will be updated eventially.
package filesystem

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// open

func FileExists(path string) bool {
	// Does file exist?
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func OpenFile(path string, printError bool) (file io.Reader, err error) {
	// Open file and return file reader or error.
	openFile, err := os.Open(path)
	if err != nil {
		if printError {
			fmt.Println(path, "does not exist!")
		}
		return nil, err
	}
	defer openFile.Close()

	return openFile, nil
}

// create

func CreateFile(path string, printError bool) (file io.Writer, err error) {
	// Open file and return file writer or error.
	createFile, err := os.Create(path)
	if err != nil {
		if printError {
			fmt.Println(err)
		}
		return nil, err
	}
	defer createFile.Close()

	return createFile, nil
}

// write

func AppendToFile(filepath string, text string, addNewLine bool) (err error) {
	// Append file or create if does not exist.
	// Add a newline at the end if necessary.

	// Open File
	file, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer file.Close()

	// if necessary, add new line.
	if addNewLine {
		text += "\n"
	}

	// Write to file.
	_, err = file.WriteString(text)
	return err
}

func AppendToFileAfter(filepath string, append string, after string) (success bool, err error) {
	// Remove string from file.

	// get started
	file, err := os.OpenFile(filepath, os.O_RDWR, 0644)
	if err != nil {
		fmt.Println(err)
		return false, err
	}

	found := false

	scanner := bufio.NewScanner(file)
	var bs []byte
	buf := bytes.NewBuffer(bs)

	var text string
	for scanner.Scan() {
		text = scanner.Text()
		if text != after {
			_, err := buf.WriteString(text + "\n")
			if err != nil {
				return false, err
			}
		} else {
			found = true
			_, err := buf.WriteString(text + "\n" + append + "\n")
			if err != nil {
				return false, err
			}
		}
	}
	file.Truncate(0)
	file.Seek(0, 0)
	buf.WriteTo(file)
	return found, nil
}

// remove string

func RemoveLineFromFile(filepath string, remove string) (success bool, err error) {
	// Remove string from file.

	// get started
	file, err := os.OpenFile(filepath, os.O_RDWR, 0644)
	if err != nil {
		fmt.Println(err)
		return false, err
	}

	found := false

	scanner := bufio.NewScanner(file)
	var bs []byte
	buf := bytes.NewBuffer(bs)

	var text string
	for scanner.Scan() {
		text = scanner.Text()
		if text != remove {
			_, err := buf.WriteString(text + "\n")
			if err != nil {
				return false, err
			}
		} else {
			found = true
		}
	}
	file.Truncate(0)
	file.Seek(0, 0)
	buf.WriteTo(file)
	return found, nil
}

// remove string

func RemoveLineWithSubstringFromFile(filepath string, substring string) (success bool, err error) {
	// Remove string from file.

	// get started
	file, err := os.OpenFile(filepath, os.O_RDWR, 0644)
	if err != nil {
		fmt.Println(err)
		return false, err
	}

	found := false

	scanner := bufio.NewScanner(file)
	var bs []byte
	buf := bytes.NewBuffer(bs)

	var text string
	for scanner.Scan() {
		text = scanner.Text()
		match, _ := regexp.MatchString(substring, text)
		if !match {
			_, err := buf.WriteString(text + "\n")
			if err != nil {
				return false, err
			}
		} else {
			found = true
		}
	}
	file.Truncate(0)
	file.Seek(0, 0)
	buf.WriteTo(file)
	return found, nil
}

// delete

func DeleteFile(filepath string) error {
	// Delete file
	return os.Remove(filepath)
}

func DeleteFolder(folderpath string) error {
	// Delete folder
	return os.RemoveAll(folderpath)
}

// copy

func CopyFile(source string, dest string, printError bool) (err error) {
	// Copy file to path.

	// Open file
	sourceFile, err := os.Open(source)
	if err != nil {
		if printError {
			fmt.Println(source, "does not exist!")
		}
		return err
	}
	defer sourceFile.Close()

	// Create new file at destination.
	destFile, err := os.Create(dest)
	if err != nil {
		if printError {
			fmt.Println(err)
		}
		return err
	}
	defer destFile.Close()

	// Copy src file
	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		// if file can't be copied -> panic!
		return err
	}

	// Update dest file permissions.
	sourceInfo, err := os.Stat(source)
	if err == nil {
		err = os.Chmod(dest, sourceInfo.Mode())
	}

	return
}

// replace

func ReplaceInFile(path string, needle string, replacement string, printError bool) (err error) {
	// Replace string in file.
	read, err := os.ReadFile(path)
	if err != nil {
		if printError {
			fmt.Println(err)
		}
		return err
	}

	newContents := strings.Replace(string(read), needle, replacement, -1)

	// fmt.Println(newContents)

	err = os.WriteFile(path, []byte(newContents), 0)
	if err != nil && printError == true {
		fmt.Println(err)
	}
	return err
}

func ReplaceInAllFiles(folder string, recursive bool, needle string, replacement string) (err error) {
	// Replace string in all files in folder.

	// check if folder exists
	if !DirExists(folder) {
		err := errors.New("Folder does not exist!")
		fmt.Println(err)
		return err
	}

	// iterate over items in folder
	dir, _ := os.Open(folder)
	objects, err := dir.Readdir(-1)

	for _, object := range objects {
		filepath := folder + "/" + object.Name()
		if object.IsDir() {
			// Replace string in sub sub-dirs
			err = ReplaceInAllFiles(filepath, recursive, needle, replacement)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			// Replace string in file.
			_ = ReplaceInFile(filepath, needle, replacement, true)
		}
	}
	return
}

func ReplaceRegex(input string, pattern string, replacement string) (output string) {
	// Replace regex in string.
	regex := regexp.MustCompile(pattern)
	return regex.ReplaceAllString(input, replacement)
}

func GetRegexInFile(path string, pattern string, replacement string, printError bool) (string, error) {
	// Get regex from file.
	read, err := os.ReadFile(path)
	if err != nil {
		if printError {
			fmt.Println(err)
		}
		return "", err
	}

	// Find regex and only output based on the pattern specified.
	regex := regexp.MustCompile(pattern)
	return regex.ReplaceAllString(regex.FindString(string(read)), replacement), nil
}

func ReplaceRegexInFile(path string, pattern string, replacement string, printError bool) (err error) {
	// Replace string in file based on regex.

	read, err := os.ReadFile(path)
	if err != nil {
		if printError {
			fmt.Println(err)
		}
		return err
	}

	newContent := ReplaceRegex(string(read), pattern, replacement)

	err = os.WriteFile(path, []byte(newContent), 0)
	if err != nil && printError == true {
		fmt.Println(err)
	}
	return err
}

func ReplaceRegexInAllFiles(folder string, recursive bool, pattern string, replacement string) (err error) {
	// Replace string in all files in folder based on regex.

	// check if folder exists
	if !DirExists(folder) {
		err := errors.New("Folder does not exist!")
		fmt.Println(err)
		return err
	}

	// iterate over items in folder
	dir, _ := os.Open(folder)
	objects, err := dir.Readdir(-1)

	for _, object := range objects {
		filepath := folder + "/" + object.Name()
		if object.IsDir() {
			// Replace in sub sub-dirs
			err = ReplaceRegexInAllFiles(filepath, recursive, pattern, replacement)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			_ = ReplaceRegexInFile(filepath, pattern, replacement, true)
		}
	}
	return
}

// find

func ContainsString(filepath string, needle string) bool {
	// File contains string?

	file, err := os.Open(filepath)
	if err != nil {
		return false
	}

	scanner := bufio.NewScanner(file)

	simpleNeedle := strings.TrimSuffix(needle, "\n")
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), simpleNeedle) {
			return true
		}
	}
	return false
}

// tail

func GetTail(filepath string, lines int) string {
	// Return last X lines of file.

	file, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	found := 0
	line := ""
	var cursor int64 = 0
	stat, _ := file.Stat()
	size := stat.Size()
	for {
		cursor -= 1
		file.Seek(cursor, io.SeekEnd)

		char := make([]byte, 1)
		file.Read(char)

		if cursor != -1 && (char[0] == 10 || char[0] == 13) {
			found += 1
			// found lines end - x
			if found >= lines {
				fmt.Println("Found line " + strconv.Itoa(found))
				break
			}
			// get next last line
			line = ""
		}
		line = fmt.Sprintf("%s%s", string(char), line)

		if cursor == -size {
			fmt.Println("At the beginning of file. Stop!")
			break
		}
	}

	fmt.Println("Full line: " + line)

	return line
}

/////////////////////////// DIR //////////////////////////

func DirExists(path string) bool {
	// Dir exists?
	if _, err := os.Open(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func CreateDir(path string) error {
	// create dir
	return os.MkdirAll(path, os.ModePerm)
}

func CopyDir(source string, dest string, printError bool, failIfExists bool) (err error) {
	// Copy entire director. Define if should fail, if the dest dir already exists.

	// does source exist
	sourceInfo, err := os.Stat(source)
	if err != nil {
		panic(err)
	}

	// is source directory
	if !sourceInfo.IsDir() {
		err := errors.New("Source is not a dir! " + source)
		panic(err)
	}

	// check if dest already exists
	if DirExists(dest) {
		err := errors.New("Dest already exists!")
		fmt.Println(err)
		if failIfExists {
			return err
		} else {
			return nil
		}
	}

	// create dir
	err = os.MkdirAll(dest, sourceInfo.Mode())
	if err != nil {
		return err
	}

	// iterate over items in folder
	dir, _ := os.Open(source)
	objects, err := dir.Readdir(-1)

	for _, object := range objects {
		sourceFile := source + "/" + object.Name()
		destFile := dest + "/" + object.Name()

		if object.IsDir() {
			// create sub-dirs
			err = CopyDir(sourceFile, destFile, printError, failIfExists)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			err = CopyFile(sourceFile, destFile, printError)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
	return
}
