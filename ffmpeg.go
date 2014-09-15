package main

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"os/exec"
	"time"
)

func Frame(filename string, timestamp time.Duration, args ...string) (image.Image, error) {
	cmd := exec.Command("./ffmpeg", append(append([]string{"-ss", Timestamp(timestamp), "-i", filename, "-frames:v", "1"}, args...), "-f", "image2", "-codec", "png", "-")...)
	cmd.Stderr = os.Stderr

	r, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	err = cmd.Start()
	if err != nil {
		return nil, err
	}

	img, err := png.Decode(r)
	if err != nil {
		return nil, err
	}

	err = cmd.Wait()
	if err != nil {
		return nil, err
	}

	return img, nil
}

func Timestamp(d time.Duration) string {
	return fmt.Sprintf("%d:%d:%f", int(d.Hours()), int(d.Minutes()), (d % time.Minute).Seconds())
}
