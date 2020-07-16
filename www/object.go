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

type ObjectTemplateVars struct {
	Header *html.Header
	Photos []*oembed.Photo
}

func NewObjectHandler(db oembed.OEmbedDatabase, t *template.Template) (http.Handler, error) {

	t = t.Lookup("object")

	if t == nil {
		return nil, errors.New("Missing 'object' template")
	}

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		ctx := req.Context()

		q := req.URL.Query()

		object_uri := q.Get("uri")

		if object_uri == "" {
			http.Error(rsp, "Missing uri", http.StatusBadRequest)
			return
		}

		ph, err := db.GetOEmbedWithObjectURI(ctx, object_uri)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusInternalServerError)
			return
		}

		header := &html.Header{}

		vars := ObjectTemplateVars{
			Header: header,
			Photos: ph,
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
