package chatgpt

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var API_KEY string

type sendQuestionBody struct {
	Messages         []Message `json:"messages"`
	Temperature      float64   `json:"temperature"`
	MaxTokens        int       `json:"max_tokens"`
	TopP             int       `json:"top_p"`
	FrequencyPenalty int       `json:"frequency_penalty"`
	PresencePenalty  int       `json:"presence_penalty"`
	Model            string    `json:"model"`
	Stream           bool      `json:"stream"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type sendQuestionResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Usage   Usage
	Choices []Choice `json:"choices"`
}

type Choice struct {
	Message       Message `json:"message"`
	Index         int     `json:"index"`
	FininshReason string  `json:"finish_reason"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

func init() {
	API_KEY = os.Getenv("API_KEY")
}

func SendQuestion(msg string) (string, error) {
	client := &http.Client{}

	body := sendQuestionBody{
		Model: "gpt-3.5-turbo",
		Messages: []Message{
			{
				Role:    "system",
				Content: "你是隻無尾熊，請以無尾熊的角度回答對話、可愛一點。",
			},
			{
				Role:    "user",
				Content: msg,
			},
		},
		Temperature: 0.7,
		MaxTokens:   256,
		TopP:        1,
	}

	jsonBytes, err := json.Marshal(body)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewReader(jsonBytes))
	if err != nil {
		return "", err
	}

	req.Header.Add("Authorization", API_KEY)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer resp.Body.Close()
	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	var response sendQuestionResponse
	json.Unmarshal([]byte(res), &response)
	if len(response.Choices) == 0 {
		return "openapi 不知道要怎麼回答你這個問題，請你檢討", nil
	}
	return response.Choices[0].Message.Content, nil
}
