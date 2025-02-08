package main

import (
	"context"
	"fmt"
	"github.com/devkaare/speech-to-text/file"
	"io"
	"os"
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
	inputDir  = "audio-recordings/"

	outputFilePath  = outputDir + outputFileName
	summaryFilePath = outputDir + summaryFileName
	inputFilePath   = inputDir + inputFileName
)

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

	if err := file.WriteToFile(outputFilePath, audioResp.Text); err != nil {
		panic(err)
	}

	fmt.Printf("[+] Successfully transcribed clip!\nFile: %s\nData: %s\n", inputFilePath, audioResp.Text)

	data, err := file.ReadFromFile(outputFilePath)
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

	if err := file.WriteToFile(summaryFilePath, chatResp.Choices[0].Message.Content); err != nil {
		panic(err)
	}

	fmt.Printf("[+] Successfully summarized data!\nData: %s\n", chatResp.Choices[0].Message.Content)
}
