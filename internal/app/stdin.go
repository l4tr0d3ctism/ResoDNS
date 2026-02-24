package app

import "os"

func HasStdin() bool {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	if stat.Mode()&os.ModeNamedPipe == 0 {
		return false
	}
	return true
}
