package test_speed

import "net/http"

type requestValidator struct {
}

func (rv requestValidator) Validate(req *http.Request) (errs []string) {
	if providers, ok := req.URL.Query()["provider"]; !ok {
		errs = append(errs, "'provider' param is required")
	} else {
		if len(providers[0]) == 0 {
			errs = append(errs, "'provider' param should not be empty")
		}
	}

	return nil
}
