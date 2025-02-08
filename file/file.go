package file

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
)

func checkFileExists(filePath string) (bool, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return false, err
	}

	if len(fileInfo.Name()) <= 0 {
		return false, err
	}

	return true, nil
}

func WriteToFile(filePath, data string) error {
	var file *os.File

	fileExists, err := checkFileExists(filePath)
	if !fileExists {
		if file, err = os.Create(filePath); err != nil {
			return err
		}
	}

	if file, err = os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600); err != nil {
		return err
	}

	defer file.Close()

	if _, err := file.WriteString(data + "\n"); err != nil {
		return err
	}

	return nil
}

func ReadFromFile(filePath string) (string, error) {
	rawData, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	return string(rawData), nil
}

func RecordAudioFile(filePath string) error {
	return nil
}

func SplitAudioFile(filePath string) error {
	cmd := exec.Command("sox", "--i", filePath)

	r, _ := cmd.StdoutPipe()

	cmd.Stderr = cmd.Stdout

	if err := cmd.Start(); err != nil {
		return err
	}

	scanner := bufio.NewScanner(r)

	durationRegex := regexp.MustCompile(`Duration\s+:\s+(\d{2}:\d{2}:\d{2}\.\d{2})`)
	bitRateRegex := regexp.MustCompile(`Bit Rate\s+:\s+([\d.]+[kKmMbB]*)`)

	var (
		duration string
		bitRate  string
	)

	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)

		if result := durationRegex.FindStringSubmatch(line); len(result) > 1 {
			duration = result[1]
			fmt.Println("[+] Extracted Time:", duration)
		}

		if result := bitRateRegex.FindStringSubmatch(line); len(result) > 1 {
			bitRate = result[1]
			fmt.Println("[+] Extracted Time:", bitRate)
		}
	}

	if err := cmd.Wait(); err != nil {
		return err
	}

	// sox example-clips/english-corporate-meeting.wav example-clips/copy-english-corporate-meeting.wav trim 0 5 : newfile : restart
	// if _, err = exec.Command("sox", filePath, filePath, "trim", "0", "5", ":", "newfile", ":", "restart").Output(); err != nil {
	// 	return err
	// }

	return nil
}
