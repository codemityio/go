package content

// text and structured data media types.
const (
	MediaTypeJSON           = "application/json"
	MediaTypeProblemJSON    = "application/problem+json"
	MediaTypeNDJSON         = "application/x-ndjson"
	MediaTypeXML            = "application/xml"
	MediaTypeXHTML          = "application/xhtml+xml"
	MediaTypeYAML           = "application/yaml"
	MediaTypeHTML           = "text/html"
	MediaTypePlain          = "text/plain"
	MediaTypeCSS            = "text/css"
	MediaTypeCSV            = "text/csv"
	MediaTypeMarkdown       = "text/markdown"
	MediaTypeJavaScript     = "text/javascript"
	MediaTypeEventStream    = "text/event-stream"
	MediaTypeWebAssembly    = "application/wasm"
	MediaTypeManifestJSON   = "application/manifest+json"
	MediaTypeSignedExchange = "application/signed-exchange"
)

// image media types.
const (
	MediaTypeAVIF = "image/avif"
	MediaTypeAPNG = "image/apng"
	MediaTypeBMP  = "image/bmp"
	MediaTypeGIF  = "image/gif"
	MediaTypeICO  = "image/vnd.microsoft.icon"
	MediaTypeJPEG = "image/jpeg"
	MediaTypeJXL  = "image/jxl"
	MediaTypePNG  = "image/png"
	MediaTypeSVG  = "image/svg+xml"
	MediaTypeTIFF = "image/tiff"
	MediaTypeWebP = "image/webp"
)

// video media types.
const (
	MediaTypeAVI       = "video/x-msvideo"
	MediaTypeMatroska  = "video/x-matroska"
	MediaTypeMP4       = "video/mp4"
	MediaTypeMPEG      = "video/mpeg"
	MediaTypeOGGVideo  = "video/ogg"
	MediaTypeQuickTime = "video/quicktime"
	MediaTypeWebMVideo = "video/webm"
	MediaTypeHLS       = "application/vnd.apple.mpegurl"
	MediaTypeDASH      = "application/dash+xml"
)

// audio media types.
const (
	MediaTypeAAC       = "audio/aac"
	MediaTypeFLAC      = "audio/flac"
	MediaTypeMP3       = "audio/mpeg"
	MediaTypeMP4Audio  = "audio/mp4"
	MediaTypeOGGAudio  = "audio/ogg"
	MediaTypeOpus      = "audio/opus"
	MediaTypeWAV       = "audio/wav"
	MediaTypeWebMAudio = "audio/webm"
)

// downloadable file and archive media types.
const (
	MediaTypeOctetStream = "application/octet-stream"
	MediaTypePDF         = "application/pdf"
	MediaTypeZIP         = "application/zip"
	MediaTypeGZIP        = "application/gzip"
	MediaTypeTAR         = "application/x-tar"
	MediaType7Z          = "application/x-7z-compressed"
	MediaTypeRAR         = "application/vnd.rar"
	MediaTypeDOC         = "application/msword"
	MediaTypeDOCX        = "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	MediaTypeXLS         = "application/vnd.ms-excel"
	MediaTypeXLSX        = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	MediaTypePPT         = "application/vnd.ms-powerpoint"
	MediaTypePPTX        = "application/vnd.openxmlformats-officedocument.presentationml.presentation"
	MediaTypeODT         = "application/vnd.oasis.opendocument.text"
	MediaTypeODS         = "application/vnd.oasis.opendocument.spreadsheet"
	MediaTypeEPUB        = "application/epub+zip"
	MediaTypeRTF         = "application/rtf"
)

// font media types.
const (
	MediaTypeOTF   = "font/otf"
	MediaTypeTTF   = "font/ttf"
	MediaTypeWOFF  = "font/woff"
	MediaTypeWOFF2 = "font/woff2"
)

// form and multipart media types.
const (
	MediaTypeFormURLEncoded      = "application/x-www-form-urlencoded"
	MediaTypeMultipartFormData   = "multipart/form-data"
	MediaTypeMultipartMixed      = "multipart/mixed"
	MediaTypeMultipartByteRanges = "multipart/byteranges"
)

const (
	mediaRangeAny        = "*/*"
	mediaRangeSubtypeAny = "/*"
	mediaRangeWildcard   = "*"
	mediaParamQuality    = "q"
)

// media range match specificity, higher value means more specific match.
const (
	specificityNone = iota - 1
	specificityAny
	specificityType
	specificityExact
)
