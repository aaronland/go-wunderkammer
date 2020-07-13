package oembed

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
)

func DataURL(ctx context.Context, url string) (string, error) {

	select {
	case <-ctx.Done():
		return "", nil
	default:
		// pass
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)

	if err != nil {
		return "", err
	}

	// make this a singleton?
	cl := &http.Client{}

	rsp, err := cl.Do(req)

	if err != nil {
		return "", err
	}

	defer rsp.Body.Close()

	content_type := rsp.Header.Get("Content-type")

	body, err := ioutil.ReadAll(rsp.Body)

	if err != nil {
		return "", err
	}

	b64_data := base64.StdEncoding.EncodeToString(body)
	data_url := fmt.Sprintf("data:%s;base64,%s", content_type, b64_data)

	return data_url, nil
}
