package serial

import (
	"fmt"

	"github.com/tarm/serial"
)

var c *serial.Config
var s *serial.Port
var buf = make([]byte, 128)
var err error

func Init(device string, baudrate int) error {
	c = &serial.Config{Name: device, Baud: baudrate}
	s, err = serial.OpenPort(c)
	if err != nil {
		return fmt.Errorf("Failed to open port: %s", err)
	}
	return nil
}

func Read() error {
	for {
		n, err := s.Read(buf)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			if closeErr := s.Close();closeErr != nil {
				return fmt.Errorf("Failed to close serial port: %s\n%s", closeErr, err)
			} else {
				return err
			}
		}
		fmt.Printf("%s", buf[:n])
	}
	return nil
}
