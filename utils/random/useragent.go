package random

import (
	browser "github.com/EDDYCJY/fake-useragent"
)

func RandomUseragent() string {
	random := browser.Random()
	return random
}
