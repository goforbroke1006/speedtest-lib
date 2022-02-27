package test_speed

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/goforbroke1006/speedtest-lib/domain"
	"github.com/goforbroke1006/speedtest-lib/pkg/content_len"
)

// HandlerMiddleware use unwarmed upgraders and return bandwidth data in JSON format
func HandlerMiddleware(sources map[string]domain.Upgrader) func(w http.ResponseWriter, req *http.Request) {
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
