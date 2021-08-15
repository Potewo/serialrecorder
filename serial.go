package serial

import (
	"fmt"
	"github.com/tarm/serial"
)

func Read() {
	c := &serial.Config{Name: "/dev/cu.ESP33SerialTest-ESP32SPP", Baud: 115200}
	s, err := serial.OpenPort(c)
	if err != nil {
		fmt.Printf("Error: %s", err)
		fmt.Printf("Failed to OpenPort")
	}

	buf := make([]byte, 128)
	for {
		n, err := s.Read(buf)
		if err != nil {
			fmt.Printf("Error: %s", err)
			fmt.Printf("Failed to Read")
		}
		fmt.Printf("%s", buf[:n])
	}
}
