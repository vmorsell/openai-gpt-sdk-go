package gpt_test

import (
	"fmt"

	"github.com/vmorsell/openai-gpt-sdk-go/gpt"
)

func Example_stream() {
	apiKey := ""
	assistantName := "HaikuGPT"

	config := gpt.NewConfig().WithAPIKey(apiKey)
	client := gpt.NewClient(config)

	msg := `Can you write a haiku about the phrase "Hello, World!"?`
	fmt.Printf("User: %s\n", msg)

	ch := make(chan *gpt.ChatCompletionChunkEvent)
	go func() {
		in := gpt.ChatCompletionInput{
			Messages: []gpt.Message{
				{
					Role:    gpt.System,
					Content: fmt.Sprintf("You are %s, a helpful assistant specialized in writing Japanese haikus.", assistantName),
				},
				{
					Role:    gpt.User,
					Content: msg,
				},
			},
			Stream: true,
		}
		err := client.ChatCompletionStream(in, ch)
		if err != nil {
			panic(err)
		}
	}()

	fmt.Printf("%s: ", assistantName)
	for {
		ev, ok := <-ch
		if !ok {
			break
		}

		if len(ev.Choices) == 0 {
			continue
		}

		if ev.Choices[0].Delta.Content != nil {
			fmt.Print(*ev.Choices[0].Delta.Content)
		}
	}
	fmt.Println()
}
