package helpers

import (
	"encoding/hex"

	"github.com/liobrdev/simplepasswords_api_gateway/utils"
)

var HexHash1 = hex.EncodeToString(utils.HashToken(VALID_EMAIL_1 + VALID_PW_1))
var HexHash2 = hex.EncodeToString(utils.HashToken(VALID_EMAIL_2 + VALID_PW_2))
