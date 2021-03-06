package tezos

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/trustwallet/blockatlas/coin"
	"github.com/trustwallet/blockatlas/models"
	"github.com/trustwallet/blockatlas/platform/tezos/source"
	"github.com/trustwallet/blockatlas/util"
	"net/http"
	"time"
)

var client = source.Client{
	HttpClient: http.DefaultClient,
}

func Setup(router gin.IRouter) {
	router.Use(util.RequireConfig("tezos.api"))
	router.Use(func(c *gin.Context) {
		client.RpcUrl = viper.GetString("tezos.api")
		c.Next()
	})
	router.GET("/:address", getTransactions)
}

func getTransactions(c *gin.Context) {
	s, err := client.GetTxsOfAddress(c.Param("address"))
	if apiError(c, err) {
		return
	}

	txs := make([]models.Tx, 0)
	for _, srcTx := range s {
		if srcTx.Type.Kind != "manager" {
			continue
		}

		if len(srcTx.Type.Operations) < 1 {
			continue
		}

		op := srcTx.Type.Operations[0]

		date, err := time.Parse("2006-01-02T15:04:05Z", op.Timestamp)
		unix := date.Unix()
		if err != nil {
			unix = 0
		}

		if op.Kind != "transaction" {
			continue
		}
		var status, errMsg string
		if !op.Failed {
			status = models.StatusCompleted
		} else {
			status = models.StatusFailed
			errMsg = "transaction failed"
		}
		txs = append(txs, models.Tx{
			Id:     srcTx.Hash,
			Coin:   coin.IndexXTZ,
			Date:   unix,
			From:   op.Src.Tz,
			To:     op.Dest.Tz,
			Fee:    op.Fee,
			Block:  op.OpLevel,
			Meta:   models.Transfer{
				Value: op.Amount,
			},
			Status: status,
			Error:  errMsg,
		})
	}

	page := models.Response(txs)
	page.Sort()
	c.JSON(http.StatusOK, &page)
}

func apiError(c *gin.Context, err error) bool {
	if err != nil {
		logrus.WithError(err).Errorf("Unhandled error: %s", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return true
	}
	return false
}

