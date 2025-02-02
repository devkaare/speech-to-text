package main

import (
	"context"
	// "encoding/json"
	"fmt"
	"io"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

func writeResp(data string) error {
	respFile, err := os.OpenFile("transcriptions/responses.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer respFile.Close()

	if _, err := respFile.WriteString(data + "\n"); err != nil {
		return err
	}
	return nil
}

func main() {
	inputFile, err := os.Open("example-clips/english-corporate-meeting.wav")
	// inputFile, err := os.Open("example-clips/norwegian-topic-explanation.wav")
	if err != nil {
		panic(err)
	}
	defer inputFile.Close()

	apiKey := os.Getenv("OPENAI_API_KEY")
	fmt.Println("[+] Loaded API key: ", apiKey)

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

	if err := writeResp(audioResp.Text); err != nil {
		panic(err)
	}

	fmt.Println("[+] Received API response:", audioResp.Text)

	// TODO: Create a summary with AI using all the responses stored

}
