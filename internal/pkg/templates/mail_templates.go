package templates

import (
	"fmt"
	"time"

	"github.com/matcornic/hermes/v2"
)

func InitializeHermes() *hermes.Hermes {
	h := &hermes.Hermes{
		Product: hermes.Product{
			Name:        "Mamlaka",
			Link:        "https://mamlaka.com",
			Logo:        "https://github.com/paulodhiambo/zbooks/blob/main/path122%204.png",
			Copyright:   fmt.Sprintf("Copyright © %d Curiocity. All rights reserved.", time.Now().Year()),
			TroubleText: "",
		},
	}
	return h
}
