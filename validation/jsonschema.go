package validation

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/xeipuuv/gojsonschema"
)

// ValidateJSONSchema performs the JSON Schema validation: https://json-schema.org/
// This fetches the schema definition via http(s) or local filesystem.
// If jsonInput is not a valid JSON or if jsonInput doesn't conform to the desired JSON schema, an error is returned.
//
// TODO: accept io.Reader instead of []byte
func ValidateJSONSchema(jsonInput []byte, desiredSchemaURI string) error {
	schemaBody, err := getDesiredSchema(desiredSchemaURI)
	if err != nil {
		return fmt.Errorf("failed to get JSON schema: %w", err)
	}

	schemaLoader := gojsonschema.NewBytesLoader(schemaBody)
	docLoader := gojsonschema.NewBytesLoader(jsonInput)

	result, err := gojsonschema.Validate(schemaLoader, docLoader)
	if err != nil {
		return fmt.Errorf("failed to validate JSON schema: %w", err)
	}

	if !result.Valid() {
		var sb strings.Builder
		for _, err := range result.Errors() {
			sb.WriteString("\n\t")
			sb.WriteString(err.String())
		}
		return fmt.Errorf("JSON doc doesn't conform to the desired JSON schema: %s", sb.String())
	}

	return nil
}

func getDesiredSchema(desiredSchemaURI string) ([]byte, error) {
	tlsConfig := &tls.Config{InsecureSkipVerify: true}
	client := http.Client{Transport: &http.Transport{TLSClientConfig: tlsConfig}}
	res, err := client.Get(desiredSchemaURI)
	if err != nil {
		return nil, fmt.Errorf("failed to get JSON schema: %w", err)
	}
	return ioutil.ReadAll(res.Body)
}
