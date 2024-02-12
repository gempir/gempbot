package ysweet

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"os/exec"
)

//go:embed createToken.cjs
var createToken string

type Factory struct {
	ysweetUrl string
}

func NewFactory(ysweetUrl string) *Factory {
	return &Factory{
		ysweetUrl: ysweetUrl,
	}
}

type TokenResponse struct {
	Url   string `json:"url"`
	DocId string `json:"docId"`
	Token string `json:"token"`
}

func (f *Factory) CreateToken(docID string) (TokenResponse, error) {
	cmd := exec.Command("node", "-", "")
	cmd.Env = append(cmd.Env, "YSWEET_URL="+f.ysweetUrl)
	cmd.Env = append(cmd.Env, "YSWEET_DOC_ID=1"+docID)

	cmd.Stdin = bytes.NewBufferString(createToken)

	var out bytes.Buffer
	var errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut

	err := cmd.Run()
	if err != nil {
		return TokenResponse{}, fmt.Errorf("%s %w", errOut.String(), err)
	}

	var tokenResponse TokenResponse
	err = json.Unmarshal(out.Bytes(), &tokenResponse)
	if err != nil {
		return TokenResponse{}, fmt.Errorf("%s %w", out.String(), err)
	}

	return tokenResponse, nil
}
