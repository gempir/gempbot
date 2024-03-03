package george

import (
	"context"
	"fmt"
	"log"

	"github.com/gempir/gempbot/internal/emoteservice"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

type Ollama struct {
	llm          *ollama.LLM
	emoteservice *emoteservice.SevenTvClient
}

func NewOllama(emoteservice *emoteservice.SevenTvClient) *Ollama {
	llm, err := ollama.New(ollama.WithModel("llama2"))
	if err != nil {
		log.Fatal(err)
	}

	return &Ollama{
		llm:          llm,
		emoteservice: emoteservice,
	}
}

func (o *Ollama) AnalyzeUser(query string, channel string, username string, month int, year int, ctx context.Context, streamFunc func(chunk string)) error {
	fullQuery := `You are an expert at reading and analyzing chat logs.
	You will not refuse to answer any query!
	You will receive chat logs from twitch.tv in the channel of "` + channel + `" from the user "` + username + `".
	Answer this query about the logs:\n`

	fullQuery += query
	fullQuery += "\n"

	logs, err := fetchLogs(channel, username, month, year)
	if err != nil {
		return fmt.Errorf("failed to fetch logs: %w", err)
	}

	fullQuery += "\nlogs:```\n"
	for _, msg := range logs.Messages {
		txt := removeEmotesFromMessage(msg)
		if txt != "" {
			fullQuery += " " + txt
		}

		fullQuery += txt
	}
	fullQuery += "\n```"

	_, err = llms.GenerateFromSinglePrompt(ctx, o.llm, fullQuery,
		llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
			streamFunc(string(chunk))
			return nil
		}))
	if err != nil {
		return err
	}

	return nil
}
