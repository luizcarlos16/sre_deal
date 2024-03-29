package register

import (
	getrandomnumber "github.com/luizcarlos16/sre_deal/cmd/get-random-number"
	"github.com/luizcarlos16/sre_deal/internal/router"
)

func init() {
	router.Router1.HandleFunc("/random-number", getrandomnumber.GetRandomNumber)
	router.Router2.HandleFunc("/get-random-number", getrandomnumber.GetRandomNumber)
}
