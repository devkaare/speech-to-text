package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

var (
	inputFileName = "english-corporate-meeting.wav"
	// inputFileName  = "norwegian-topic-explanation.wav"
	outputFileName  = "output.txt"
	summaryFileName = "summary.txt"

	outputID = time.Now().Format(time.RFC3339)

	outputDir = "transcriptions/" + outputID + "/"
	inputDir  = "example-clips/"

	outputFilePath  = outputDir + outputFileName
	summaryFilePath = outputDir + summaryFileName
	inputFilePath   = inputDir + inputFileName
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

func writeToFile(filePath, data string) error {
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

func readFromFile(filePath string) (string, error) {
	rawData, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	return string(rawData), nil
}

func recordAudioFile(filePath string) error {
	return nil
}

func splitAudioFile(filePath string) error {
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

func main() {
	inputFile, err := os.Open(inputFilePath)
	if err != nil {
		panic(err)
	}
	defer inputFile.Close()

	apiKey := os.Getenv("OPENAI_API_KEY")
	client := openai.NewClient(
		option.WithAPIKey(apiKey),
	)

	audioResp, err := client.Audio.Translations.New(
		context.TODO(),
		openai.AudioTranslationNewParams{
			File:  openai.F[io.Reader](inputFile),
			Model: openai.F(openai.AudioModelWhisper1),
		},
	)
	if err != nil {
		panic(err)
	}

	if err := os.Mkdir(outputDir, 0777); err != nil {
		panic(err)
	}

	if err := writeToFile(outputFilePath, audioResp.Text); err != nil {
		panic(err)
	}

	fmt.Printf("[+] Successfully transcribed clip!\nFile: %s\nData: %s\n", inputFilePath, audioResp.Text)

	data, err := readFromFile(outputFilePath)
	if err != nil {
		panic(err)
	}

	chatResp, err := client.Chat.Completions.New(
		context.TODO(),
		openai.ChatCompletionNewParams{
			Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
				openai.UserMessage("Create a short summary of this audio clip transcription: " + data),
			}),
			Model: openai.F(openai.ChatModelGPT4o),
		},
	)
	if err != nil {
		panic(err)
	}

	if err := writeToFile(summaryFilePath, chatResp.Choices[0].Message.Content); err != nil {
		panic(err)
	}

	fmt.Printf("[+] Successfully summarized data!\nData: %s\n", chatResp.Choices[0].Message.Content)
}
