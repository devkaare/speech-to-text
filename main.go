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

type APIResp struct {
	text string
}

func main() {
	file, err := os.Open("example-clips/english-corporate-meeting.wav")
	// file, err := os.Open("example-clips/norwegian-topic-explanation.wav")
	if err != nil {
		panic(err)
	}

	apiKey := os.Getenv("OPENAI_API_KEY")
	fmt.Println("[+] Loaded API key: ", apiKey)

	client := openai.NewClient(
		option.WithAPIKey(apiKey),
	)

	audioResp, err := client.Audio.Translations.New(
		context.TODO(),
		openai.AudioTranslationNewParams{
			File:  openai.F[io.Reader](file),
			Model: openai.F(openai.AudioModelWhisper1),
		},
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("[+] Received API response:", audioResp.Text)
}
