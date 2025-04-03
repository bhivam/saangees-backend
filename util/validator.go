package util

import "regexp"

type Validator struct {
	Errors map[string]string
}

func New() *Validator {
	return &Validator{Errors: make(map[string]string)}
}

func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

func (v *Validator) AddError(key, value string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = value
	}
}

func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.Errors[key] = message
	}
}

// Helper functions for common checks
func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

func Unique[T comparable](values []T) bool {
	uniqueValues := make(map[T]bool)
	for _, value := range values {
		uniqueValues[value] = true
	}
	return len(values) == len(uniqueValues)
}

func Nonempty(value string) bool {
	return len(value) > 0
}

func MaxLen(value string, length int) bool {
	return len(value) <= length
}

func MinLen(value string, length int) bool {
	return len(value) >= length
}

func NonNegativeFl(value float64) bool {
	return value > 0
}
