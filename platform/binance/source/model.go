package source

import (
	"errors"
	"fmt"
	"github.com/trustwallet/blockatlas/models"
)

type Account struct {
	AccountNumber int       `json:"account_number"`
	Address       string    `json:"address"`
	Balances      []Balance `json:"balances"`
	PublicKey     []byte    `json:"public_key"`
	Sequence      uint64    `json:"sequence"`
}

type Balance struct {
	Symbol        string `json:"symbol"`
	Free          uint64 `json:"free"`
	Locked        uint64 `json:"locked"`
	Frozen        uint64 `json:"frozen"`
}

type Error struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
}

type Tx struct {
	BlockHeight   uint64 `json:"blockHeight"`
	Code          int    `json:"code"`
	ConfirmBlocks int    `json:"confirmBlocks"`
	Data          string `json:"data"`
	FromAddr      string `json:"fromAddr"`
	OrderId       string `json:"orderId"`
	Timestamp     string `json:"timeStamp"`
	ToAddr        string `json:"toAddr"`
	Age           int64  `json:"txAge"`
	Asset         string `json:"txAsset"`
	Fee           models.Amount `json:"txFee"`
	Hash          string `json:"txHash"`
	Value         models.Amount `json:"value"`
}

type TxPage struct {
	Total uint64 `json:"total"`
	Txs   []Tx   `json:"tx"`
}

var ErrSourceConn  = errors.New("connection to servers failed")
var ErrInvalidAddr = errors.New("invalid address")
var ErrNotFound    = errors.New("not found")

func (e *Error) Error() string {
	return fmt.Sprintf("%d: %s", e.Code, e.Message)
}
