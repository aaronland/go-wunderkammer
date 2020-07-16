package www

// THIS IS EARLY STAGES AND EVERYTHING IS IN FLUX

import (
	"errors"
	"github.com/aaronland/go-wunderkammer/html"
	"github.com/aaronland/go-wunderkammer/oembed"
	"html/template"
	_ "log"
	"net/http"
)

type ImageTemplateVars struct {
	Header *html.Header
	Photo  *oembed.Photo
}

func NewImageHandler(db oembed.OEmbedDatabase, t *template.Template) (http.Handler, error) {

	t = t.Lookup("image")

	if t == nil {
		return nil, errors.New("Missing 'object' template")
	}

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		ctx := req.Context()

		q := req.URL.Query()

		url := q.Get("url")

		if url == "" {
			http.Error(rsp, "Missing url", http.StatusBadRequest)
			return
		}

		ph, err := db.GetOEmbedWithURL(ctx, url)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusInternalServerError)
			return
		}

		header := &html.Header{}

		vars := ImageTemplateVars{
			Header: header,
			Photo:  ph,
		}

		err = t.Execute(rsp, vars)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusInternalServerError)
			return
		}

		return
	}

	h := http.HandlerFunc(fn)
	return h, nil
}
