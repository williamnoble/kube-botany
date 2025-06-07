package fs

import (
	"fmt"
	"io"
	"os"
)

func CopyImage(srcFileName string, dstFileName string) error {
	src, err := os.Open(srcFileName)
	if err != nil {
		return fmt.Errorf("error opening source file: %w", err)
	}
	defer src.Close()

	dst, err := os.Create(dstFileName)
	if err != nil {
		return fmt.Errorf("error opening destinationImg file: %w", err)
	}

	_, err = io.Copy(dst, src)
	if err != nil {
		return fmt.Errorf("error copying file contents: %w", err)
	}

	return nil
}
