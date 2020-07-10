package oembed

import (
	"context"
)

type OEmbedRecord struct {
	Version      string `json:"version,xml:"version""`
	Type         string `json:"type"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	Title        string `json:"title"`
	URL          string `json:"url"`
	AuthorName   string `json:"author_name"`
	AuthorURL    string `json:"author_url"`
	ProviderName string `json:"provider_name"`
	ProviderURL  string `json:"provider_url"`
	ObjectURI    string `json:"object_uri"`
	DataURL      string `json:"data_url"`
}

type OEmbedDatabase interface {
	AddOEmbed(context.Context, *OEmbedRecord) error
	GetRandomOEmbed(context.Context) (*OEmbedRecord, error)
	GetOEmbedWithObjectURI(context.Context, string) ([]*OEmbedRecord, error)
	Close() error
}
