package george

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var blockListedUsers = map[string]string{
	"supibot":        "supibot",
	"nightbot":       "nightbot",
	"streamelements": "streamelements",
	"streamlabs":     "streamlabs",
	"moobot":         "moobot",
	"gempbot":        "gempbot",
	"botnextdoor":    "botnextdoor",
	"botbear1110":    "botbear1110",
}

type Logs struct {
	Messages []Message `json:"messages"`
}

type Message struct {
	Text        string    `json:"text"`
	DisplayName string    `json:"displayName"`
	Timestamp   time.Time `json:"timestamp"`
	ID          string    `json:"id"`
	Tags        struct {
		ID               string `json:"id"`
		BadgeInfo        string `json:"badge-info"`
		Emotes           string `json:"emotes"`
		DisplayName      string `json:"display-name"`
		UserType         string `json:"user-type"`
		ReturningChatter string `json:"returning-chatter"`
		Color            string `json:"color"`
		Flags            string `json:"flags"`
		Mod              string `json:"mod"`
		UserID           string `json:"user-id"`
		RoomID           string `json:"room-id"`
		Subscriber       string `json:"subscriber"`
		TmiSentTs        string `json:"tmi-sent-ts"`
		FirstMsg         string `json:"first-msg"`
		Turbo            string `json:"turbo"`
		Badges           string `json:"badges"`
	} `json:"tags,omitempty"`
	Username string `json:"username"`
	Channel  string `json:"channel"`
	Raw      string `json:"raw"`
	Type     int    `json:"type"`
}

func fetchLogs(channel string, username string, month int, year int, day int) (Logs, error) {
	// Fetch logs for the given username, month and year
	// https://logs.ivr.fi/channel/nymn/user/gempir/2024/3?json

	var resp *http.Response
	var err error
	if day == 0 {
		resp, err = http.Get(fmt.Sprintf("https://logs.ivr.fi/channel/%s/user/%s/%d/%d?json", channel, username, year, month))
		if err != nil {
			return Logs{}, err
		}
	} else {
		resp, err = http.Get(fmt.Sprintf("https://logs.ivr.fi/channel/%s/%d/%d/%d?json", channel, year, month, day))
		if err != nil {
			return Logs{}, err
		}
	}

	if resp.StatusCode != http.StatusOK {
		return Logs{}, fmt.Errorf("%s", resp.Status)
	}

	defer resp.Body.Close()

	var logs Logs
	err = json.NewDecoder(resp.Body).Decode(&logs)
	if err != nil {
		return Logs{}, err
	}

	return logs, nil
}

func (o *Ollama) cleanMessage(msg Message, regexes []*regexp.Regexp) string {
	if _, ok := blockListedUsers[msg.Username]; ok {
		return ""
	}
	if strings.HasPrefix(msg.Text, "!") {
		return ""
	}

	emoteRanges := parseEmoteRanges(msg.Tags.Emotes)
	var cleanedText strings.Builder
	prevEnd := 0

	for _, er := range emoteRanges {
		// Ensure the start position is within bounds
		if er.Start > prevEnd {
			cleanedText.WriteString(msg.Text[prevEnd:er.Start])
		}
		prevEnd = er.End + 1 // Start of the next section after emote
	}

	// Append any text after the last emote
	if prevEnd < len(msg.Text) {
		cleanedText.WriteString(msg.Text[prevEnd:])
	}

	clean := cleanedText.String()

	// Use compiled regex patterns in the loop
	for _, regex := range regexes {
		if regex.MatchString(clean) {
			clean = regex.ReplaceAllString(clean, "")
		}
	}

	return strings.TrimSpace(clean)
}

type EmoteRange struct {
	Start int
	End   int
}

func parseEmoteRanges(emoteTags string) []EmoteRange {
	var emoteRanges []EmoteRange

	// Regular expression to match emote ranges
	re := regexp.MustCompile(`(\d+)-(\d+)`)

	// Split emote tags by "/"
	tags := strings.Split(emoteTags, "/")

	for _, tag := range tags {
		// Extract the start and end positions from each emote range
		matches := re.FindStringSubmatch(tag)
		if len(matches) >= 3 {
			start, _ := strconv.Atoi(matches[1])
			end, _ := strconv.Atoi(matches[2])
			emoteRanges = append(emoteRanges, EmoteRange{Start: start, End: end})
		}
	}

	return emoteRanges
}
