package mtgfail

import (
	"errors"
	"fmt"
	"net/url"
	"path"

	"github.com/inconshreveable/log15"
)

var (
	// supported is a map of supported sites mapped to the specific postfix for text only decklists
	// example:
	// "tappedout.net":"?fmt=txt""
	// Given the specific URL, https://deckbox.org/sets/2805419 the text only deck link is https://deckbox.org/sets/2805419?fmt=txt
	// Usage:
	//
	supportedWithPath = map[string]string{
		"deckbox.org": "/export",
	}

	supportedWithParam = map[string]*pair{
		"tappedout.net": {
			key:   "fmt",
			value: "txt",
		},
	}
	supported map[string]bool
)

type pair struct {
	key   string
	value string
}

func init() {
	supported = make(map[string]bool)
	for k := range supportedWithPath {
		supported[k] = true
	}
	for k := range supportedWithParam {
		supported[k] = true
	}
}

// IsSupportedDomain is a predicate that returns true if the provided URL can be used. This means it is a page of one of our supported site and it is a valid URL. Any error is swallowed.
func IsSupportedDomain(uri string) bool {
	valid, _, _ := SanitizeURL(uri)
	return valid
}

// SupportedDomain checks to see if the provided URL is a page on one of our supported site and returns any error
func SupportedDomain(uri string) (bool, error) {
	valid, _, err := SanitizeURL(uri)
	return valid, err
}

// SanitizeURL parses a URL string for validity and creates a new *url.URL if it is approved.
func SanitizeURL(uri string) (valid bool, cleaned *url.URL, err error) {
	var (
		u *url.URL
	)
	u, err = url.ParseRequestURI(uri)
	if err != nil {
		return false, nil, fmt.Errorf("Unable to parse domain %s", uri)
	}

	host := u.Hostname()
	if host == "" {
		return false, nil, errors.New("empty hostname")
	}
	if !supported[host] {
		fmt.Println("unsupported host", host, supported)
		return false, nil, nil
	}

	return true, u, nil
}
func TranslateDomain(uri string) (*url.URL, error) {
	log := log15.New()
	valid, cleaned, err := SanitizeURL(uri)
	if err != nil {
		return nil, err
	}

	if !valid || cleaned == nil {
		return nil, fmt.Errorf("invalid upstream domain: %s", uri)
	}

	host := cleaned.Hostname()
	if pair, ok := supportedWithParam[host]; ok && pair != nil {
		log.Debug("supported with params",
			"pair", pair,
			"host", host,
			"url", cleaned,
		)
		// query params
		q := cleaned.Query()
		q.Set(pair.key, pair.value)
		cleaned.RawQuery = q.Encode()
		return cleaned, nil
	} else if seg, ok := supportedWithPath[host]; ok && seg != "" {
		log.Debug("supported with path postfix",
			"pathseg", seg,
			"host", host,
			"url", cleaned,
		)
		// path segment
		cleaned.Path = path.Join(cleaned.Path, seg)
		return cleaned, nil
	}
	log.Debug("unsupported",
		"host", host,
		"url", cleaned,
	)
	return nil, fmt.Errorf("unsupported upstream")
}
