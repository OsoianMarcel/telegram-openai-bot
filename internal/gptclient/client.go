package gptclient

import (
	"context"
	"errors"
	"time"

	gogpt "github.com/sashabaranov/go-gpt3"
)

type Client struct {
	gptAuthToken string
	GptClient    *gogpt.Client
}

// Creates new OpenAI API client.
//
// It requires a GPT Auth Token.
func NewClient(gptAuthToken string) *Client {
	return &Client{
		gptAuthToken: gptAuthToken,
		GptClient:    gogpt.NewClient(gptAuthToken),
	}
}

// Ask the OpenAI a question and get the response.
//
// The "user" parameter is unique identifier representing your end-user,
// which can help OpenAI to monitor and detect abuse.
func (client *Client) AskAI(ctx context.Context, question string, user string) (string, error) {
	timeoutCtx, cancelCtx := context.WithTimeout(ctx, time.Duration(AI_REQ_TIMEOUT_SEC*time.Second))
	defer cancelCtx()

	req := gogpt.CompletionRequest{
		Model:       AI_MODEL,
		MaxTokens:   AI_MAX_TOKENS,
		Temperature: AI_TEMPERATURE,
		Prompt:      question,
		User:        user,
	}

	resp, err := client.GptClient.CreateCompletion(timeoutCtx, req)
	if err != nil {
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", errors.New("the completion has not choices")
	}

	return resp.Choices[0].Text, nil
}
