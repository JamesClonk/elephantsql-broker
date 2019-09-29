package broker

import (
	"io/ioutil"

	"github.com/JamesClonk/elephantsql-broker/log"
)

func init() {
	log.SetOutput(ioutil.Discard)
}
