package requester

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/daniilsolovey/gaming-task/internal/config"
	"github.com/reconquest/karma-go"
	"github.com/reconquest/pkg/log"
)

const (
	CREATE_PLAYER_REQUEST_BODY = `{
		"jsonrpc": "2.0",
		"method": "Player.Set",
		"id": 1928822491,
		"params": {
			"Id": "noname",
			"Nick": "Noname",
			"BankGroupId": "new_bank_group"
			}
		}`

	CREATE_BANK_GROUP_REQUEST_BODY = `{
		"jsonrpc": "2.0", "method": "BankGroup.Set", "id": 1225625456, "params": {
		"Id": "new_bank_group", "Currency": "EUR"
		}}
	`

	CREATE_SESSION_REQUEST_BODY = `{
		"jsonrpc": "2.0", "method": "Session.Create", "id": 321864203,
		"params": {
		"PlayerId": "noname",
		"GameId": "bennys_the_biggest_game" }
		}`
)

type ResponseCreatePlayer struct {
	JSONRPC string   `json:"jsonrpc"`
	ID      int      `json:"id"`
	Result  []string `json:"result"`
}

type ResponseBankGroup struct {
	JSONRPC string   `json:"jsonrpc"`
	ID      int      `json:"id"`
	Result  []string `json:"result"`
}

type ResponseSession struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      int             `json:"id"`
	Result  []SessionResult `json:"result"`
}

type SessionResult struct {
	SessionID  string `json:"sessionId"`
	SessionURL string `json:"sessionUrl"`
}

type Requester struct {
	config *config.Config
}

func NewRequester(
	config *config.Config,
) *Requester {
	return &Requester{
		config: config,
	}
}

type RequesterInterface interface {
	CreatePlayer() (*ResponseCreatePlayer, error)
	CreateBankGroup() (*ResponseBankGroup, error)
	CreateSession() (*ResponseSession, error)
}

func prepareBody(textBody string) ([]byte, error) {
	var jsonData []byte
	jsonData, err := json.Marshal(textBody)
	if err != nil {
		return nil, karma.Format(
			err,
			"unable to marshal text_body with: %s",
			textBody,
		)
	}

	return jsonData, nil
}

func (requester *Requester) CreatePlayer() (*ResponseCreatePlayer, error) {
	log.Info("creating player on platform")
	url := requester.config.Platform.URL
	requestBody, err := prepareBody(CREATE_PLAYER_REQUEST_BODY)
	if err != nil {
		return nil, karma.Format(
			err,
			"unable to prepare text_body for create_player request",
		)
	}

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, karma.Format(
			err,
			"unable to create request by path: %s",
			url,
		)
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	client, err := requester.setupClient()
	if err != nil {
		return nil, karma.Format(
			err,
			"unable to setup client for request by url: %s", url,
		)
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, karma.Format(
			err,
			"unable to send http request by url: %s", url,
		)
	}

	defer response.Body.Close()

	var result ResponseCreatePlayer
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return nil, karma.Format(
			err,
			"unable to decode response for player, response status code: %d ",
			response.StatusCode,
		)
	}

	return &result, nil
}

func (requester *Requester) CreateBankGroup() (*ResponseBankGroup, error) {
	log.Info("creating bank_group on platform")
	url := requester.config.Platform.URL
	requestBody, err := prepareBody(CREATE_BANK_GROUP_REQUEST_BODY)
	if err != nil {
		return nil, karma.Format(
			err,
			"unable to prepare text_body for bank_group request",
		)
	}

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, karma.Format(
			err,
			"unable to create request by path: %s",
			url,
		)
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	client, err := requester.setupClient()
	if err != nil {
		return nil, karma.Format(
			err,
			"unable to setup client for request by url: %s", url,
		)
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, karma.Format(
			err,
			"unable to send http request by url: %s", url,
		)
	}

	defer response.Body.Close()

	var result ResponseBankGroup
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return nil, karma.Format(
			err,
			"unable to decode response for bank_group, response status code: %d ",
			response.StatusCode,
		)
	}

	return &result, nil
}

func (requester *Requester) CreateSession() (*ResponseSession, error) {
	log.Info("creating session on platform")
	url := requester.config.Platform.URL
	requestBody, err := prepareBody(CREATE_SESSION_REQUEST_BODY)
	if err != nil {
		return nil, karma.Format(
			err,
			"unable to prepare text_body for session request",
		)
	}

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, karma.Format(
			err,
			"unable to create request by path: %s",
			url,
		)
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	client, err := requester.setupClient()
	if err != nil {
		return nil, karma.Format(
			err,
			"unable to setup client for request by url: %s", url,
		)
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, karma.Format(
			err,
			"unable to send http request by url: %s", url,
		)
	}

	defer response.Body.Close()

	var result ResponseSession
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return nil, karma.Format(
			err,
			"unable to decode response for session, response status code: %d ",
			response.StatusCode,
		)
	}

	return &result, nil
}

func (requester *Requester) setupClient() (*http.Client, error) {
	caCert, err := ioutil.ReadFile(requester.config.SSLPath)
	if err != nil {
		return nil, karma.Format(
			err,
			"unable to read ssl by path:% s",
			requester.config.SSLPath,
		)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: caCertPool,
			},
		},

		Timeout: 1 * time.Minute,
	}

	return client, nil
}