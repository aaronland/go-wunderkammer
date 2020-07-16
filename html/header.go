package html

import (
	"time"
)

type Header struct {
	Title string
	Links []Link
	Date  time.Time
}

type Link struct {
	Title string
	Type  string
	Href  string
}
