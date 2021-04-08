package mtgfail

import (
	"errors"
	"fmt"
	"net/url"
)

var (
	// supported is a map of supported sites mapped to the specific postfix for text only decklists
	// example:
	// "tappedout.net":"?fmt=txt""
	// Given the specific URL, https://deckbox.org/sets/2805419 the text only deck link is https://deckbox.org/sets/2805419?fmt=txt
	// Usage:
	//
	supported = map[string]string{
		"tappedout.net": "?fmt=txt",
		"deckbox.org":   "/export",
	}
)

// IsSupportedDomain is a predicate that returns true if the provided URL can be used. This means it is a page of one of our supported site and it is a valid URL. Any error is swallowed.
func IsSupportedDomain(uri string) bool {
	valid, _, _ := ParseDomain(uri)
	return valid
}

// SupportedDomain checks to see if the provided URL is a page on one of our supported site and returns any error
func SupportedDomain(uri string) (bool, error) {
	valid, _, err := ParseDomain(uri)
	return valid, err
}

// ParseDomain parses a URL for validity and then checks any associated postfix. It returns
func ParseDomain(uri string) (valid bool, fullURI string, err error) {
	var (
		postfix string
		u       *url.URL
	)
	u, err = url.Parse(uri)
	if err != nil {
		return false, "", fmt.Errorf("Unable to parse domain %s", uri)
	}

	host := u.Hostname()
	if host == "" {
		return false, "", errors.New("empty hostname")
	}
	if postfix, valid = supported[host]; !valid {
		return false, "", nil
	}

	return valid, fmt.Sprintf("%s%s", uri, postfix), nil
}
func TranslateDomain(uri string) (string, error) {
	_, fullURL, err := ParseDomain(uri)
	return fullURL, err
}
