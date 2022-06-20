package validator

import (
	"encoding/json"
	"errors"
	"io/ioutil"

	"github.com/go-chassis/openlog"
	"github.com/xeipuuv/gojsonschema"
)

//ValidatePaylaod function
func ValidatePaylaod(schemaPath string, payload map[string]interface{}) (interface{}, error) {
	payload_bytes, _ := json.Marshal(payload)
	schema_bytes, _ := ioutil.ReadFile(schemaPath)
	schemaLoader := gojsonschema.NewBytesLoader(schema_bytes)
	documentLoader := gojsonschema.NewBytesLoader(payload_bytes)
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		openlog.Error("error occured here" + err.Error())
		return err, errors.New("Invalid Payload")
	}
	if result.Valid() {
		return nil, nil
	} else {

		validationErrors := make([]string, 0)
		for _, desc := range result.Errors() {
			validationErrors = append(validationErrors, desc.String())
		}
		openlog.Error("schema validation errors")
		return validationErrors, errors.New("Invalid Payload")
	}
}
