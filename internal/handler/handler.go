package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/daniilsolovey/gaming-task/internal/config"
	"github.com/daniilsolovey/gaming-task/internal/database"
	"github.com/daniilsolovey/gaming-task/internal/requester"
	"github.com/gin-gonic/gin"
	"github.com/reconquest/pkg/log"
)

type Handler struct {
	database  *database.Database
	config    *config.Config
	requester requester.RequesterInterface
}

func NewHandler(
	database *database.Database,
	config *config.Config,
	requester requester.RequesterInterface,
) *Handler {
	return &Handler{
		database:  database,
		config:    config,
		requester: requester,
	}
}

type NewPlayerBalance struct {
	Balance float64 `json:"balance"`
}

type BalanceResponse struct {
	JSONRPC string `json:"jsonrpc"`
	ID      string `json:"id"`
	Result  struct {
		Balance float64 `json:"balance"`
	} `json:"result"`
}

func (handler *Handler) StartServer(config *config.Config) {
	router := gin.Default()
	router.GET("/", handler.ActionIndex)
	router.GET("/get_player_balance/:id", handler.GetPlayerBalanceForPlatform)
	router.POST("/create_player", handler.CreatePlayer)
	router.POST("/create_bank_group", handler.CreateBankGroup)
	router.POST("/create_session", handler.CreateSession)
	router.POST("/update_player_balance/:id", handler.UpdatePlayerBalanceFromPlatform)

	router.Run(handler.config.Handler.Port)
}

func (handler *Handler) ActionIndex(context *gin.Context) {
	context.Data(
		200,
		"text/plain; charset=UTF-8",
		[]byte("This api version: "+handler.config.Handler.ApiVersion),
	)
}

func (handler *Handler) GetPlayerBalanceForPlatform(context *gin.Context) {
	playerID := context.Param("id")
	playerFromDatabase, err := handler.database.GetPlayerByID(playerID)
	if err != nil {
		log.Error(err)
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "player not found"})
		return
	}
	var response BalanceResponse
	response.ID = playerID
	response.JSONRPC = "2.0"
	response.Result.Balance = playerFromDatabase.Balance
	context.IndentedJSON(http.StatusOK, response)
}

func (handler *Handler) UpdatePlayerBalanceFromPlatform(context *gin.Context) {
	playerID := context.Param("id")
	body, err := ioutil.ReadAll(context.Request.Body)
	if err != nil {
		log.Error(err)
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "unable to read request body"})
		return
	}

	var balance NewPlayerBalance
	err = json.Unmarshal(body, &balance)
	if err != nil {
		log.Errorf(
			err,
			"unable to unmarshal data from body: %s", string(body))
		return
	}

	err = handler.database.UpdatePlayerBalance(playerID, balance.Balance)
	if err != nil {
		log.Error(err)
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "player not found"})
		return
	}

	context.IndentedJSON(http.StatusOK, gin.H{"message": "balance updated"})

}

func (handler *Handler) CreatePlayer(context *gin.Context) {
	player, err := handler.requester.CreatePlayer()
	if err != nil {
		log.Error(err)
	}

	player = &requester.ResponseCreatePlayer{}
	player.ID = 1928822491
	player.JSONRPC = "2.0"
	log.Info("player created: ", player)

	err = handler.database.InsertPlayer(*player, 113.12)
	if err != nil {
		log.Error(err)
	}

	context.IndentedJSON(http.StatusOK, player)

}

func (handler *Handler) CreateBankGroup(context *gin.Context) {
	bankGroup, err := handler.requester.CreateBankGroup()
	if err != nil {
		log.Error(err)
	}

	bankGroup = &requester.ResponseBankGroup{}
	bankGroup.ID = 1225625456
	bankGroup.JSONRPC = "2.0"
	log.Info("bankGroup created: ", bankGroup)
	context.IndentedJSON(http.StatusOK, bankGroup)

}

func (handler *Handler) CreateSession(context *gin.Context) {
	session, err := handler.requester.CreateSession()
	if err != nil {
		log.Error(err)
	}

	session = &requester.ResponseSession{}
	session.ID = 321864203
	session.JSONRPC = "2.0"
	sessionResult := &requester.SessionResult{SessionID: "exampleSession", SessionURL: "exampleURL"}
	session.Result = append(session.Result, *sessionResult)
	log.Info("session created: ", session)
	context.IndentedJSON(http.StatusOK, session)

}
