package validator

import (
	"net/url"
	"regexp"
)

var youtubeRegex = regexp.MustCompile(`^(https?:\/\/)?(www\.)?(youtube\.com|youtu\.?be)\/.+$`)

func IsValidYouTubeURL(urlString string) bool {
	parsedURL, err := url.Parse(urlString)
	if err != nil {
		return false
	}

	if parsedURL.Scheme == "" {
		urlString = "https://" + urlString
	}

	return youtubeRegex.MatchString(urlString)
}
