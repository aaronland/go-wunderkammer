package oembed

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
)

func DataURL(ctx context.Context, url string) (string, error) {

	rsp, err := http.Get(url)

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
