package test_speed

import (
	"encoding/json"
	"fmt"
	"github.com/goforbroke1006/speedtest-lib/upgrader"
	"net/http"

	"github.com/goforbroke1006/speedtest-lib/pkg/content_len"
)

func TestSpeedHandler(sources map[string]upgrader.Upgrader) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		rv := requestValidator{}
		if errs := rv.Validate(req); len(errs) > 0 {
			w.WriteHeader(http.StatusBadRequest)
			message := ""
			for _, err := range errs {
				message += err + "\n"
			}
			_, _ = w.Write([]byte(message))
			return
		}

		provider := req.URL.Query()["provider"][0]

		cUpgrader, ok := sources[provider]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(fmt.Sprintf("unexpected provider type: %s", provider)))
			return
		}

		resp := Response{
			DLSpeed: content_len.DataLen(cUpgrader.GetDLSpeedMbps()).MegaBites(),
			ULSpeed: content_len.DataLen(cUpgrader.GetULSpeedMbps()).MegaBites(),
		}
		respBody, err := json.Marshal(resp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(fmt.Sprintf("error: %s", err.Error())))
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(respBody)
	}
}
