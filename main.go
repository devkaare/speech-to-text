package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

func main() {
	fmt.Println("[+] Opening audio file...")
	file, err := os.Open("test.wav")
	if err != nil {
		panic(err)
	}

	apiKey := os.Getenv("OPENAI_API_KEY")
	fmt.Println("[+] Making new client...")
	client := openai.NewClient(
		option.WithAPIKey(apiKey),
	)

	fmt.Println("[+] Sending request...")
	audioResp, err := client.Audio.Translations.New(
		context.TODO(),
		openai.AudioTranslationNewParams{
			File: openai.F[io.Reader](file),
		},
	)
	if err != nil {
		panic(err)
	}

	fmt.Println(audioResp)
}
