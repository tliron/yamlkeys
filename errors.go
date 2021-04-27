package yamlkeys

import (
	"fmt"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

//
// DecodeError
//

type DecodeError struct {
	Message string
	Line    int
	Column  int
}

func NewDecodeError(message string, line int, column int) *DecodeError {
	return &DecodeError{
		Message: message,
		Line:    line,
		Column:  column,
	}
}

// error interface
func (self *DecodeError) Error() string {
	if self.Line != -1 {
		if self.Column != -1 {
			return fmt.Sprintf("malformed YAML @%d,%d: %s", self.Line, self.Column, self.Message)
		} else {
			return fmt.Sprintf("malformed YAML @%d: %s", self.Line, self.Message)
		}
	} else {
		return fmt.Sprintf("malformed YAML: %s", self.Message)
	}
}

func NewDecodeErrorFor(message string, node *yaml.Node) *DecodeError {
	return NewDecodeError(message, node.Line, node.Column)
}

func NewDuplicateKeyErrorFor(key interface{}, node *yaml.Node) *DecodeError {
	return NewDecodeErrorFor(fmt.Sprintf("duplicate map key: %s", key), node)
}

func WrapWithDecodeError(err error) error {
	// Unfortunately, "gopkg.in/yaml.v3" just uses fmt.Errorf to create its errors,
	// so the only way we can extract line number information is by parsing the error string

	message := err.Error()
	if strings.HasPrefix(message, "yaml: ") {
		if strings.HasPrefix(message, "yaml: line ") {
			suffix := message[11:]
			if colon := strings.Index(suffix, ": "); colon != -1 {
				line := suffix[:colon]
				if row, err := strconv.Atoi(line); err == nil {
					return NewDecodeError(suffix[colon+2:], row, -1)
				}
			}
		} else {
			return NewDecodeError(message[6:], -1, -1)
		}
	}

	return err
}
