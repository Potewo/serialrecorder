package file

import (
	"os"
	"fmt"
)

var file *os.File

func Open(filename string) error {
	_file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		return fmt.Errorf("Failed to open file\n%s", err)
	}
	file = _file
	return nil
}

func Close() error {
	if err := file.Close(); err != nil {
		return fmt.Errorf("Failed to close file.\n%s", err)
	}
	return nil
}

func Append(s string) error {
	_, err := fmt.Fprint(file, s)
	if err != nil {
		return fmt.Errorf("Failed to append string to file\n%s", err)
	}
	return nil
}
