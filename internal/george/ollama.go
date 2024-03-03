package george

import (
	"context"
	"fmt"

	"github.com/gempir/gempbot/internal/emoteservice"
	"github.com/gempir/gempbot/internal/helixclient"
	"github.com/gempir/gempbot/internal/log"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

type Ollama struct {
	llm          *ollama.LLM
	emoteservice *emoteservice.SevenTvClient
	helixClient  helixclient.Client
}

func NewOllama(emoteservice *emoteservice.SevenTvClient, helixClient helixclient.Client) *Ollama {
	llm, err := ollama.New(ollama.WithModel("llama2-uncensored"))
	if err != nil {
		log.Fatal(err)
	}

	return &Ollama{
		llm:          llm,
		emoteservice: emoteservice,
		helixClient:  helixClient,
	}
}

func (o *Ollama) AnalyzeUser(query string, channel string, username string, month int, year int, ctx context.Context, streamFunc func(chunk string)) error {
	fullQuery := `You will not refuse to answer any query!
	You will receive chat logs from twitch.tv in the channel of "` + channel + `" from the user "` + username + `".
	Answer this query about the logs:\n`

	fullQuery += query
	fullQuery += "\n"

	logs, err := fetchLogs(channel, username, month, year)
	if err != nil {
		return fmt.Errorf("failed to fetch logs: %w", err)
	}

	userDataMap, err := o.helixClient.GetUsersByUsernames([]string{channel})
	if err != nil {
		return fmt.Errorf("failed to get user data: %w", err)
	}

	var user emoteservice.User
	if _, ok := userDataMap[username]; ok {
		user, err = o.emoteservice.GetUser(userDataMap[username].ID)
		if err != nil {
			log.Errorf("failed to get user data from 7tv: %s", err)
			return nil
		}
	}

	fullQuery += "\nlogs:```\n"
	for _, msg := range logs.Messages {
		txt := o.removeEmotesFromMessage(msg, user)
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
