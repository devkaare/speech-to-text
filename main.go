package main

import (
	"context"
	"fmt"
	"io"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

const (
	inputPath  = "example-clips/english-corporate-meeting.wav"
	outputPath = "transcriptions/responses.txt"
)

func writeOutput(data string) error {
	outputFile, err := os.OpenFile(outputPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	if _, err := outputFile.WriteString(data + "\n"); err != nil {
		return err
	}
	return nil
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

	if err := writeOutput(audioResp.Text); err != nil {
		panic(err)
	}

	fmt.Printf("Successfully transcribed clip!\nFile: %s\nData: %s", inputPath, audioResp.Text)

}
