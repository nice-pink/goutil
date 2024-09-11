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
	"path"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/nice-pink/goutil/pkg/log"
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

// list

func IsOlderThan(timeStamp time.Time, sec int64) bool {
	now := time.Now()
	passed := time.Duration(sec) * time.Second
	return timeStamp.Add(passed).Before(now)
}

func ListFiles(folder string, olderThanSeconds int64, ignoreHiddenFiles bool) []string {
	files, err := os.ReadDir(folder)
	if err != nil {
		log.Err(err)
	}

	filenames := []string{}
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		// ignore hidden files
		if ignoreHiddenFiles && strings.HasPrefix(file.Name(), ".") {
			continue
		}

		filepath := filepath.Join(folder, file.Name())
		fileInfo, err := os.Stat(filepath)
		if err != nil {
			log.Err(err)
			continue
		}

		// check time
		if olderThanSeconds <= 0 {
			filenames = append(filenames, filepath)
		}

		if IsOlderThan(fileInfo.ModTime(), olderThanSeconds) {
			// append file
			filenames = append(filenames, filepath)
		}
	}

	return filenames
}

// delete

func DeleteFile(filepath string) error {
	// Delete file
	return os.Remove(filepath)
}

func DeleteFiles(filepaths []string) {
	// Delete file
	for _, path := range filepaths {
		err := os.Remove(path)
		if err != nil {
			log.Info("---")
			log.Err(err)
			log.Info("---")
		}
	}
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

func ReplaceInFile(path string, needle string, replacement string, printError bool) (replaced bool, err error) {
	// Replace string in file.
	read, err := os.ReadFile(path)
	if err != nil {
		if printError {
			fmt.Println(err)
		}
		return false, err
	}

	newContent := strings.Replace(string(read), needle, replacement, -1)

	eq := bytes.Equal(read, []byte(newContent))
	if eq {
		return false, nil
	}
	// fmt.Println(newContents)

	err = os.WriteFile(path, []byte(newContent), 0)
	if err != nil && printError {
		fmt.Println(err)
	}
	return false, err
}

func ReplaceInAllFiles(folder string, recursive bool, needle string, replacement string) (replaced bool, err error) {
	// Replace string in all files in folder.

	// check if folder exists
	if !DirExists(folder) {
		err := errors.New("folder does not exist")
		fmt.Println(err)
		return false, err
	}

	// iterate over items in folder
	objects, err := os.ReadDir(folder)
	if err != nil {
		fmt.Println(err)
		return false, err
	}

	for _, object := range objects {
		filepath := folder + "/" + object.Name()
		if object.IsDir() {
			// Replace string in sub sub-dirs
			_, err = ReplaceInAllFiles(filepath, recursive, needle, replacement)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			// Replace string in file.
			r, _ := ReplaceInFile(filepath, needle, replacement, true)
			if r {
				replaced = true
			}
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

func GetAllRegexInFile(path string, pattern string, replacement string, printError bool) ([]string, error) {
	// Get regex from file.
	read, err := os.ReadFile(path)
	if err != nil {
		if printError {
			fmt.Println(err)
		}
		return nil, err
	}

	// Find regex and only output based on the pattern specified.
	values := []string{}
	regex := regexp.MustCompile(pattern)
	items := regex.FindAllString(string(read), -1)
	for _, item := range items {
		values = append(values, regex.ReplaceAllString(item, replacement))
	}
	return values, nil
}

func GetRegexInAllFiles(folder string, recursive bool, pattern string, replacement string, fileExtensions []string) ([]string, error) {
	// check if folder exists
	if !DirExists(folder) {
		err := errors.New("folder does not exist")
		fmt.Println(err)
		return nil, err
	}

	// iterate over items in folder
	objects, err := os.ReadDir(folder)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	values := []string{}
	for _, object := range objects {
		filepath := folder + "/" + object.Name()
		// fmt.Println(filepath)
		if object.IsDir() {
			// Replace in sub sub-dirs
			val, err := GetRegexInAllFiles(filepath, recursive, pattern, replacement, fileExtensions)
			if err != nil {
				fmt.Println(err)
			}
			if len(val) > 0 {
				values = append(values, val...)
			}
		} else {
			if len(fileExtensions) > 0 && !slices.Contains(fileExtensions, path.Ext(filepath)) {
				continue
			}
			val, _ := GetAllRegexInFile(filepath, pattern, replacement, false)
			if len(val) > 0 {
				values = append(values, val...)
			}
		}
	}
	return values, nil
}

func ReplaceRegexInFile(path string, pattern string, replacement string, printError bool) (replaced bool, err error) {
	// Replace string in file based on regex.

	read, err := os.ReadFile(path)
	if err != nil {
		if printError {
			fmt.Println(err)
		}
		return false, err
	}

	newContent := ReplaceRegex(string(read), pattern, replacement)

	eq := bytes.Equal(read, []byte(newContent))
	if eq {
		return false, nil
	}

	err = os.WriteFile(path, []byte(newContent), 0)
	if err != nil && printError {
		fmt.Println(err)
	}
	return true, err
}

func ReplaceRegexInAllFiles(folder string, recursive bool, pattern string, replacement string) (replaced bool, err error) {
	// Replace string in all files in folder based on regex.

	// check if folder exists
	if !DirExists(folder) {
		err := errors.New("folder does not exist")
		fmt.Println(err)
		return false, err
	}

	// iterate over items in folder
	objects, err := os.ReadDir(folder)
	if err != nil {
		fmt.Println(err)
		return false, err
	}

	for _, object := range objects {
		filepath := folder + "/" + object.Name()
		if object.IsDir() {
			// Replace in sub sub-dirs
			_, err = ReplaceRegexInAllFiles(filepath, recursive, pattern, replacement)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			r, _ := ReplaceRegexInFile(filepath, pattern, replacement, true)
			if r {
				replaced = true
			}
		}
	}
	return replaced, nil
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

func FindAllStringsInFile(path string, pattern string) []string {
	read, err := os.ReadFile(path)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	regex := regexp.MustCompile(pattern)
	return regex.FindAllString(string(read), -1)
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
		err := errors.New("source is not a dir: " + source)
		panic(err)
	}

	// check if dest already exists
	if DirExists(dest) {
		err := errors.New("dest already exists")
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
	objects, err := os.ReadDir(source)

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
