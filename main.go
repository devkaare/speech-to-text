package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

var (
	inputPath = "example-clips/english-corporate-meeting.wav"
	// inputPath  = "example-clips/norwegian-topic-explanation.wav"
	outputPath  = "transcriptions/output.txt" + time.Now().GoString()
	summaryPath = "transcriptions/summary.txt" + time.Now().GoString()
)

func writeToFile(filePath, data string) error {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
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

func main() {
	inputFile, err := os.Open(inputPath)
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

	if err := writeToFile(outputPath, audioResp.Text); err != nil {
		panic(err)
	}

	fmt.Printf("[+] Successfully transcribed clip!\nFile: %s\nData: %s\n", inputPath, audioResp.Text)

	data, err := readFromFile(outputPath)
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

	if err := writeToFile(summaryPath, chatResp.Choices[0].Message.Content); err != nil {
		panic(err)
	}

	fmt.Printf("[+] Successfully summarized data!\nData: %s\n", chatResp.Choices[0].Message.Content)
}
