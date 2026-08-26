package validate

import "fmt"

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type Errors struct {
	Fields []FieldError `json:"fields"`
}

func (e *Errors) Error() string {
	if e == nil || len(e.Fields) == 0 {
		return "validation failed"
	}
	return fmt.Sprintf("validation failed for %d field(s)", len(e.Fields))
}

func (e *Errors) Add(field, message string) {
	e.Fields = append(e.Fields, FieldError{Field: field, Message: message})
}

func (e *Errors) Empty() bool {
	return e == nil || len(e.Fields) == 0
}

func (e *Errors) Err() error {
	if e.Empty() {
		return nil
	}
	return e
}
