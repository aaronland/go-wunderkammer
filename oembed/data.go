package oembed

import (
	"bufio"
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"github.com/aaronland/go-image-decode"
	"github.com/aaronland/go-image-encode"
	"github.com/aaronland/go-image-halftone"
	"github.com/aaronland/go-image-resize"
	"github.com/aaronland/go-image-rotate"
	"github.com/esimov/caire"
	"image"
	"image/draw"
	"log"
	"io/ioutil"
	"net/http"
	"strings"
)

type DataURLOptions struct {
	ContentAwareResize bool
	ContentAwareHeight int
	ContentAwareWidth  int
	Halftone           bool
	MaxDimension	int
}

func DataURL(ctx context.Context, url string, opts *DataURLOptions) (string, error) {

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

	if strings.HasPrefix(content_type, "image/") {

		dec, err := decode.NewDecoder(ctx, "image://")
		
		if err != nil {
			return "", err
		}

		// FIX ME...
		enc, err := encode.NewEncoder(ctx, "gif://")
		
		if err != nil {
			return "", err
		}
		
		br := bytes.NewReader(body)

		im, format, err := dec.Decode(ctx, br)

		if err != nil {
			return "", err
		}

		orientation := "0"

		if format == "jpeg" {

			_, err := br.Seek(0, 0)

			if err != nil {
				return "", err
			}

			o, err := rotate.GetImageOrientation(ctx, br)

			if err != nil {
				log.Println(err)	
			} else {
				orientation = o
			}
		}

		new_im, err := rotate.RotateImageWithOrientation(ctx, im, orientation)

		if err != nil {
			return "", err
		}

		var content_aware_err error

		if opts.ContentAwareResize {

			caire_w := opts.ContentAwareWidth
			caire_h := opts.ContentAwareHeight

			// failing here because... ?
			// https://github.com/esimov/caire/blob/eb499d00d8b9e45511b0a5fc3418b26b24123081/process.go#L165

			pr := &caire.Processor{
				NewWidth:  caire_w,
				NewHeight: caire_h,
				Scale:     true,
			}

			b := new_im.Bounds()
			caire_im := image.NewNRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
			draw.Draw(caire_im, caire_im.Bounds(), new_im, b.Min, draw.Src)

			resized_im, err := pr.Resize(caire_im)

			if err != nil {
				log.Printf("Failed to resize %s, %v\n", url, err)
				content_aware_err = err
			} else {
				new_im = resized_im
			}

		}

		if !opts.ContentAwareResize || content_aware_err != nil {

			new_im, err = resize.ResizeImageMax(ctx, new_im, opts.MaxDimension)

			if err != nil {
				return "", err
			}
		}

		// end of caire stuff

		if opts.Halftone {

			opts := halftone.NewDefaultHalftoneOptions()
			new_im, err = halftone.HalftoneImage(ctx, new_im, opts)

			if err != nil {
				return "", err
			}
		}

		var buf bytes.Buffer
		wr := bufio.NewWriter(&buf)

		err = enc.Encode(ctx, new_im, wr)

		if err != nil {
			return "", err
		}

		wr.Flush()

		body = buf.Bytes()
	}

	b64_data := base64.StdEncoding.EncodeToString(body)
	data_url := fmt.Sprintf("data:%s;base64,%s", content_type, b64_data)

	return data_url, nil
}
