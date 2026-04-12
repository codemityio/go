package huma

// Config is a struct containing server configuration values.
type Config struct {
	Title         string `env:"TITLE,default=API"                       json:"title"`
	Version       string `env:"VERSION,default=v0.0.0"                  json:"version"`
	OpenAPIPath   string `env:"OPEN_API_PATH,default=/openapi"          json:"openApiPath"`
	DocsPath      string `env:"DOCS_PATH,default=/docs"                 json:"docsPath"`
	SchemasPath   string `env:"SCHEMAS_PATH,default=/schemas"           json:"schemasPath"`
	DefaultFormat string `env:"DEFAULT_FORMAT,default=application/json" json:"defaultFormat"`
}
