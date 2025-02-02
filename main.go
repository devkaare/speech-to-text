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
	inputPath = "example-clips/english-corporate-meeting.wav"
	// inputPath  = "example-clips/norwegian-topic-explanation.wav"
	outputPath = "transcriptions/output.txt"
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

	fmt.Printf("Successfully transcribed clip!\nFile: %s\nData: %s\n", inputPath, audioResp.Text)

	chatResp, err := client.Chat.Completions.New(
		context.TODO(),
		openai.ChatCompletionNewParams{
			Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
				openai.UserMessage("Say this is a test"),
			}),
			Model: openai.F(openai.ChatModelGPT4o),
		},
	)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Successfully summarized data!\nData: %s\n", chatResp.Choices[0].Message.Content)
}
