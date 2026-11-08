package content

import (
	"mime"
	"strconv"
	"strings"
)

// Negotiate returns the best matching media type offered for the `Accept` header value (e.g. `text/html` for
// browser requests), offers are listed in the order of preference and the first one is used when the header is
// empty or none of the offers is acceptable.
func Negotiate(accept string, offers ...string) string {
	if len(offers) == 0 {
		return ""
	}

	if strings.TrimSpace(accept) == "" {
		return offers[0]
	}

	best, bestQuality := offers[0], 0.0

	for _, offer := range offers {
		quality := offerQuality(accept, offer)

		// Strict comparison keeps the order of preference for offers of the same quality.
		if quality > bestQuality {
			best, bestQuality = offer, quality
		}
	}

	return best
}

// offerQuality returns the quality of the offer taken from the most specific matching media range.
func offerQuality(accept, offer string) float64 {
	quality, specificity := 0.0, specificityNone

	for mediaRange := range strings.SplitSeq(accept, ",") {
		mediaType, params, err := mime.ParseMediaType(strings.TrimSpace(mediaRange))
		if err != nil {
			continue
		}

		s := matchSpecificity(mediaType, offer)
		if s <= specificity {
			continue
		}

		quality, specificity = 1.0, s

		if value, ok := params[mediaParamQuality]; ok {
			if q, e := strconv.ParseFloat(value, 64); e == nil {
				quality = q
			}
		}
	}

	return quality
}

// matchSpecificity returns how specifically the media range matches the offer.
func matchSpecificity(mediaRange, offer string) int {
	switch {
	case mediaRange == offer:
		return specificityExact

	case mediaRange == mediaRangeAny:
		return specificityAny

	case strings.HasSuffix(mediaRange, mediaRangeSubtypeAny) &&
		strings.HasPrefix(offer, strings.TrimSuffix(mediaRange, mediaRangeWildcard)):
		return specificityType

	default:
		return specificityNone
	}
}
