package content

import (
	"mime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNegotiate(t *testing.T) {
	offers := []string{MediaTypeJSON, MediaTypeHTML, MediaTypePlain}

	tests := []struct {
		name     string
		accept   string
		offers   []string
		expected string
	}{
		{
			name: "browser",
			accept: "text/html,application/xhtml+xml,application/xml;q=0.9,image/jxl,image/avif,image/webp," +
				"image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
			offers:   offers,
			expected: MediaTypeHTML,
		},
		{
			name:     "empty",
			accept:   "",
			offers:   offers,
			expected: MediaTypeJSON,
		},
		{
			name:     "any",
			accept:   "*/*",
			offers:   offers,
			expected: MediaTypeJSON,
		},
		{
			name:     "exact_json",
			accept:   "application/json",
			offers:   offers,
			expected: MediaTypeJSON,
		},
		{
			name:     "exact_plain",
			accept:   "text/plain",
			offers:   offers,
			expected: MediaTypePlain,
		},
		{
			name:     "case_insensitive",
			accept:   "Text/Plain",
			offers:   offers,
			expected: MediaTypePlain,
		},
		{
			name:     "type_wildcard",
			accept:   "text/*",
			offers:   offers,
			expected: MediaTypeHTML,
		},
		{
			name:     "quality",
			accept:   "application/json;q=0.5, text/plain;q=0.9",
			offers:   offers,
			expected: MediaTypePlain,
		},
		{
			name:     "specific_overrides_wildcard",
			accept:   "text/*;q=0.9, text/html;q=0.1",
			offers:   offers,
			expected: MediaTypePlain,
		},
		{
			name:     "not_acceptable_quality",
			accept:   "text/html;q=0, */*;q=0.1",
			offers:   offers,
			expected: MediaTypeJSON,
		},
		{
			name:     "no_match",
			accept:   "image/png",
			offers:   offers,
			expected: MediaTypeJSON,
		},
		{
			name:     "malformed",
			accept:   ";;, text/plain",
			offers:   offers,
			expected: MediaTypePlain,
		},
		{
			name:     "no_offers",
			accept:   "text/html",
			offers:   nil,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, Negotiate(tt.accept, tt.offers...))
		})
	}
}

func TestNegotiateMediaTypes(t *testing.T) {
	const fallback = "x-fallback/none"

	mediaTypes := []string{
		MediaTypeJSON,
		MediaTypeProblemJSON,
		MediaTypeNDJSON,
		MediaTypeXML,
		MediaTypeXHTML,
		MediaTypeYAML,
		MediaTypeHTML,
		MediaTypePlain,
		MediaTypeCSS,
		MediaTypeCSV,
		MediaTypeMarkdown,
		MediaTypeJavaScript,
		MediaTypeEventStream,
		MediaTypeWebAssembly,
		MediaTypeManifestJSON,
		MediaTypeSignedExchange,
		MediaTypeAVIF,
		MediaTypeAPNG,
		MediaTypeBMP,
		MediaTypeGIF,
		MediaTypeICO,
		MediaTypeJPEG,
		MediaTypeJXL,
		MediaTypePNG,
		MediaTypeSVG,
		MediaTypeTIFF,
		MediaTypeWebP,
		MediaTypeAVI,
		MediaTypeMatroska,
		MediaTypeMP4,
		MediaTypeMPEG,
		MediaTypeOGGVideo,
		MediaTypeQuickTime,
		MediaTypeWebMVideo,
		MediaTypeHLS,
		MediaTypeDASH,
		MediaTypeAAC,
		MediaTypeFLAC,
		MediaTypeMP3,
		MediaTypeMP4Audio,
		MediaTypeOGGAudio,
		MediaTypeOpus,
		MediaTypeWAV,
		MediaTypeWebMAudio,
		MediaTypeOctetStream,
		MediaTypePDF,
		MediaTypeZIP,
		MediaTypeGZIP,
		MediaTypeTAR,
		MediaType7Z,
		MediaTypeRAR,
		MediaTypeDOC,
		MediaTypeDOCX,
		MediaTypeXLS,
		MediaTypeXLSX,
		MediaTypePPT,
		MediaTypePPTX,
		MediaTypeODT,
		MediaTypeODS,
		MediaTypeEPUB,
		MediaTypeRTF,
		MediaTypeOTF,
		MediaTypeTTF,
		MediaTypeWOFF,
		MediaTypeWOFF2,
		MediaTypeFormURLEncoded,
		MediaTypeMultipartFormData,
		MediaTypeMultipartMixed,
		MediaTypeMultipartByteRanges,
	}

	for _, mediaType := range mediaTypes {
		t.Run(mediaType, func(t *testing.T) {
			parsed, _, err := mime.ParseMediaType(mediaType)
			require.NoError(t, err)
			assert.Equal(t, mediaType, parsed)

			mainType, _, _ := strings.Cut(mediaType, "/")
			typeRange := mainType + mediaRangeSubtypeAny
			rejected := mediaType + ";" + mediaParamQuality + "=0"

			assert.Equal(t, mediaType, Negotiate(mediaType, fallback, mediaType), "exact")
			assert.Equal(t, mediaType, Negotiate(typeRange, fallback, mediaType), "type")
			assert.Equal(t, fallback, Negotiate(mediaRangeAny, fallback, mediaType), "any")
			assert.Equal(t, fallback, Negotiate(rejected, fallback, mediaType), "q0")
		})
	}
}
