// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: author_service.proto

package author_service

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"google.golang.org/protobuf/types/known/anypb"
)

// ensure the imports are used
var (
	_ = bytes.MinRead
	_ = errors.New("")
	_ = fmt.Print
	_ = utf8.UTFMax
	_ = (*regexp.Regexp)(nil)
	_ = (*strings.Reader)(nil)
	_ = net.IPv4len
	_ = time.Duration(0)
	_ = (*url.URL)(nil)
	_ = (*mail.Address)(nil)
	_ = anypb.Any{}
	_ = sort.Sort
)

// Validate checks the field values on Author with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *Author) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Author with the rules defined in the
// proto definition for this message. If any rules are violated, the result is
// a list of violation errors wrapped in AuthorMultiError, or nil if none found.
func (m *Author) ValidateAll() error {
	return m.validate(true)
}

func (m *Author) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Id

	// no validation rules for Name

	// no validation rules for Biography

	// no validation rules for Version

	// no validation rules for CreatedAt

	// no validation rules for UpdatedAt

	if len(errors) > 0 {
		return AuthorMultiError(errors)
	}

	return nil
}

// AuthorMultiError is an error wrapping multiple validation errors returned by
// Author.ValidateAll() if the designated constraints aren't met.
type AuthorMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m AuthorMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m AuthorMultiError) AllErrors() []error { return m }

// AuthorValidationError is the validation error returned by Author.Validate if
// the designated constraints aren't met.
type AuthorValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e AuthorValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e AuthorValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e AuthorValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e AuthorValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e AuthorValidationError) ErrorName() string { return "AuthorValidationError" }

// Error satisfies the builtin error interface
func (e AuthorValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sAuthor.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = AuthorValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = AuthorValidationError{}

// Validate checks the field values on CreateAuthorRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *CreateAuthorRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on CreateAuthorRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// CreateAuthorRequestMultiError, or nil if none found.
func (m *CreateAuthorRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *CreateAuthorRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if utf8.RuneCountInString(m.GetName()) < 3 {
		err := CreateAuthorRequestValidationError{
			field:  "Name",
			reason: "value length must be at least 3 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if utf8.RuneCountInString(m.GetBiography()) < 1 {
		err := CreateAuthorRequestValidationError{
			field:  "Biography",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return CreateAuthorRequestMultiError(errors)
	}

	return nil
}

// CreateAuthorRequestMultiError is an error wrapping multiple validation
// errors returned by CreateAuthorRequest.ValidateAll() if the designated
// constraints aren't met.
type CreateAuthorRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m CreateAuthorRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m CreateAuthorRequestMultiError) AllErrors() []error { return m }

// CreateAuthorRequestValidationError is the validation error returned by
// CreateAuthorRequest.Validate if the designated constraints aren't met.
type CreateAuthorRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e CreateAuthorRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e CreateAuthorRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e CreateAuthorRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e CreateAuthorRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e CreateAuthorRequestValidationError) ErrorName() string {
	return "CreateAuthorRequestValidationError"
}

// Error satisfies the builtin error interface
func (e CreateAuthorRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sCreateAuthorRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = CreateAuthorRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = CreateAuthorRequestValidationError{}

// Validate checks the field values on CreateAuthorResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *CreateAuthorResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on CreateAuthorResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// CreateAuthorResponseMultiError, or nil if none found.
func (m *CreateAuthorResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *CreateAuthorResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if all {
		switch v := interface{}(m.GetAuthor()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, CreateAuthorResponseValidationError{
					field:  "Author",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, CreateAuthorResponseValidationError{
					field:  "Author",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetAuthor()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return CreateAuthorResponseValidationError{
				field:  "Author",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return CreateAuthorResponseMultiError(errors)
	}

	return nil
}

// CreateAuthorResponseMultiError is an error wrapping multiple validation
// errors returned by CreateAuthorResponse.ValidateAll() if the designated
// constraints aren't met.
type CreateAuthorResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m CreateAuthorResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m CreateAuthorResponseMultiError) AllErrors() []error { return m }

// CreateAuthorResponseValidationError is the validation error returned by
// CreateAuthorResponse.Validate if the designated constraints aren't met.
type CreateAuthorResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e CreateAuthorResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e CreateAuthorResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e CreateAuthorResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e CreateAuthorResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e CreateAuthorResponseValidationError) ErrorName() string {
	return "CreateAuthorResponseValidationError"
}

// Error satisfies the builtin error interface
func (e CreateAuthorResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sCreateAuthorResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = CreateAuthorResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = CreateAuthorResponseValidationError{}

// Validate checks the field values on GetAuthorRequest with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *GetAuthorRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GetAuthorRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// GetAuthorRequestMultiError, or nil if none found.
func (m *GetAuthorRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *GetAuthorRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if utf8.RuneCountInString(m.GetId()) < 1 {
		err := GetAuthorRequestValidationError{
			field:  "Id",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return GetAuthorRequestMultiError(errors)
	}

	return nil
}

// GetAuthorRequestMultiError is an error wrapping multiple validation errors
// returned by GetAuthorRequest.ValidateAll() if the designated constraints
// aren't met.
type GetAuthorRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GetAuthorRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GetAuthorRequestMultiError) AllErrors() []error { return m }

// GetAuthorRequestValidationError is the validation error returned by
// GetAuthorRequest.Validate if the designated constraints aren't met.
type GetAuthorRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetAuthorRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetAuthorRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetAuthorRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetAuthorRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetAuthorRequestValidationError) ErrorName() string { return "GetAuthorRequestValidationError" }

// Error satisfies the builtin error interface
func (e GetAuthorRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetAuthorRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetAuthorRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetAuthorRequestValidationError{}

// Validate checks the field values on GetAuthorResponse with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *GetAuthorResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GetAuthorResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// GetAuthorResponseMultiError, or nil if none found.
func (m *GetAuthorResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *GetAuthorResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if all {
		switch v := interface{}(m.GetAuthor()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, GetAuthorResponseValidationError{
					field:  "Author",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, GetAuthorResponseValidationError{
					field:  "Author",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetAuthor()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return GetAuthorResponseValidationError{
				field:  "Author",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return GetAuthorResponseMultiError(errors)
	}

	return nil
}

// GetAuthorResponseMultiError is an error wrapping multiple validation errors
// returned by GetAuthorResponse.ValidateAll() if the designated constraints
// aren't met.
type GetAuthorResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GetAuthorResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GetAuthorResponseMultiError) AllErrors() []error { return m }

// GetAuthorResponseValidationError is the validation error returned by
// GetAuthorResponse.Validate if the designated constraints aren't met.
type GetAuthorResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetAuthorResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetAuthorResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetAuthorResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetAuthorResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetAuthorResponseValidationError) ErrorName() string {
	return "GetAuthorResponseValidationError"
}

// Error satisfies the builtin error interface
func (e GetAuthorResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetAuthorResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetAuthorResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetAuthorResponseValidationError{}

// Validate checks the field values on ListAuthorsRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *ListAuthorsRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on ListAuthorsRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// ListAuthorsRequestMultiError, or nil if none found.
func (m *ListAuthorsRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *ListAuthorsRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if m.GetPage() < 1 {
		err := ListAuthorsRequestValidationError{
			field:  "Page",
			reason: "value must be greater than or equal to 1",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if m.GetPageSize() < 1 {
		err := ListAuthorsRequestValidationError{
			field:  "PageSize",
			reason: "value must be greater than or equal to 1",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return ListAuthorsRequestMultiError(errors)
	}

	return nil
}

// ListAuthorsRequestMultiError is an error wrapping multiple validation errors
// returned by ListAuthorsRequest.ValidateAll() if the designated constraints
// aren't met.
type ListAuthorsRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ListAuthorsRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ListAuthorsRequestMultiError) AllErrors() []error { return m }

// ListAuthorsRequestValidationError is the validation error returned by
// ListAuthorsRequest.Validate if the designated constraints aren't met.
type ListAuthorsRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ListAuthorsRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ListAuthorsRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ListAuthorsRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ListAuthorsRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ListAuthorsRequestValidationError) ErrorName() string {
	return "ListAuthorsRequestValidationError"
}

// Error satisfies the builtin error interface
func (e ListAuthorsRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sListAuthorsRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ListAuthorsRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ListAuthorsRequestValidationError{}

// Validate checks the field values on ListAuthorsResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *ListAuthorsResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on ListAuthorsResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// ListAuthorsResponseMultiError, or nil if none found.
func (m *ListAuthorsResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *ListAuthorsResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	for idx, item := range m.GetAuthors() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, ListAuthorsResponseValidationError{
						field:  fmt.Sprintf("Authors[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, ListAuthorsResponseValidationError{
						field:  fmt.Sprintf("Authors[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return ListAuthorsResponseValidationError{
					field:  fmt.Sprintf("Authors[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	// no validation rules for TotalItems

	// no validation rules for TotalPages

	if len(errors) > 0 {
		return ListAuthorsResponseMultiError(errors)
	}

	return nil
}

// ListAuthorsResponseMultiError is an error wrapping multiple validation
// errors returned by ListAuthorsResponse.ValidateAll() if the designated
// constraints aren't met.
type ListAuthorsResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ListAuthorsResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ListAuthorsResponseMultiError) AllErrors() []error { return m }

// ListAuthorsResponseValidationError is the validation error returned by
// ListAuthorsResponse.Validate if the designated constraints aren't met.
type ListAuthorsResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ListAuthorsResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ListAuthorsResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ListAuthorsResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ListAuthorsResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ListAuthorsResponseValidationError) ErrorName() string {
	return "ListAuthorsResponseValidationError"
}

// Error satisfies the builtin error interface
func (e ListAuthorsResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sListAuthorsResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ListAuthorsResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ListAuthorsResponseValidationError{}

// Validate checks the field values on UpdateAuthorRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *UpdateAuthorRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on UpdateAuthorRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// UpdateAuthorRequestMultiError, or nil if none found.
func (m *UpdateAuthorRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *UpdateAuthorRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if utf8.RuneCountInString(m.GetId()) < 1 {
		err := UpdateAuthorRequestValidationError{
			field:  "Id",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if utf8.RuneCountInString(m.GetName()) < 3 {
		err := UpdateAuthorRequestValidationError{
			field:  "Name",
			reason: "value length must be at least 3 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if utf8.RuneCountInString(m.GetBiography()) < 1 {
		err := UpdateAuthorRequestValidationError{
			field:  "Biography",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if m.GetVersion() < 1 {
		err := UpdateAuthorRequestValidationError{
			field:  "Version",
			reason: "value must be greater than or equal to 1",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return UpdateAuthorRequestMultiError(errors)
	}

	return nil
}

// UpdateAuthorRequestMultiError is an error wrapping multiple validation
// errors returned by UpdateAuthorRequest.ValidateAll() if the designated
// constraints aren't met.
type UpdateAuthorRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m UpdateAuthorRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m UpdateAuthorRequestMultiError) AllErrors() []error { return m }

// UpdateAuthorRequestValidationError is the validation error returned by
// UpdateAuthorRequest.Validate if the designated constraints aren't met.
type UpdateAuthorRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e UpdateAuthorRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e UpdateAuthorRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e UpdateAuthorRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e UpdateAuthorRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e UpdateAuthorRequestValidationError) ErrorName() string {
	return "UpdateAuthorRequestValidationError"
}

// Error satisfies the builtin error interface
func (e UpdateAuthorRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sUpdateAuthorRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = UpdateAuthorRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = UpdateAuthorRequestValidationError{}

// Validate checks the field values on UpdateAuthorResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *UpdateAuthorResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on UpdateAuthorResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// UpdateAuthorResponseMultiError, or nil if none found.
func (m *UpdateAuthorResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *UpdateAuthorResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if all {
		switch v := interface{}(m.GetAuthor()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, UpdateAuthorResponseValidationError{
					field:  "Author",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, UpdateAuthorResponseValidationError{
					field:  "Author",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetAuthor()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return UpdateAuthorResponseValidationError{
				field:  "Author",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return UpdateAuthorResponseMultiError(errors)
	}

	return nil
}

// UpdateAuthorResponseMultiError is an error wrapping multiple validation
// errors returned by UpdateAuthorResponse.ValidateAll() if the designated
// constraints aren't met.
type UpdateAuthorResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m UpdateAuthorResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m UpdateAuthorResponseMultiError) AllErrors() []error { return m }

// UpdateAuthorResponseValidationError is the validation error returned by
// UpdateAuthorResponse.Validate if the designated constraints aren't met.
type UpdateAuthorResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e UpdateAuthorResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e UpdateAuthorResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e UpdateAuthorResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e UpdateAuthorResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e UpdateAuthorResponseValidationError) ErrorName() string {
	return "UpdateAuthorResponseValidationError"
}

// Error satisfies the builtin error interface
func (e UpdateAuthorResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sUpdateAuthorResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = UpdateAuthorResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = UpdateAuthorResponseValidationError{}

// Validate checks the field values on DeleteAuthorRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *DeleteAuthorRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on DeleteAuthorRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// DeleteAuthorRequestMultiError, or nil if none found.
func (m *DeleteAuthorRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *DeleteAuthorRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if utf8.RuneCountInString(m.GetId()) < 1 {
		err := DeleteAuthorRequestValidationError{
			field:  "Id",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if m.GetVersion() < 1 {
		err := DeleteAuthorRequestValidationError{
			field:  "Version",
			reason: "value must be greater than or equal to 1",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return DeleteAuthorRequestMultiError(errors)
	}

	return nil
}

// DeleteAuthorRequestMultiError is an error wrapping multiple validation
// errors returned by DeleteAuthorRequest.ValidateAll() if the designated
// constraints aren't met.
type DeleteAuthorRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m DeleteAuthorRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m DeleteAuthorRequestMultiError) AllErrors() []error { return m }

// DeleteAuthorRequestValidationError is the validation error returned by
// DeleteAuthorRequest.Validate if the designated constraints aren't met.
type DeleteAuthorRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e DeleteAuthorRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e DeleteAuthorRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e DeleteAuthorRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e DeleteAuthorRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e DeleteAuthorRequestValidationError) ErrorName() string {
	return "DeleteAuthorRequestValidationError"
}

// Error satisfies the builtin error interface
func (e DeleteAuthorRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sDeleteAuthorRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = DeleteAuthorRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = DeleteAuthorRequestValidationError{}

// Validate checks the field values on DeleteAuthorResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *DeleteAuthorResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on DeleteAuthorResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// DeleteAuthorResponseMultiError, or nil if none found.
func (m *DeleteAuthorResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *DeleteAuthorResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Message

	if len(errors) > 0 {
		return DeleteAuthorResponseMultiError(errors)
	}

	return nil
}

// DeleteAuthorResponseMultiError is an error wrapping multiple validation
// errors returned by DeleteAuthorResponse.ValidateAll() if the designated
// constraints aren't met.
type DeleteAuthorResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m DeleteAuthorResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m DeleteAuthorResponseMultiError) AllErrors() []error { return m }

// DeleteAuthorResponseValidationError is the validation error returned by
// DeleteAuthorResponse.Validate if the designated constraints aren't met.
type DeleteAuthorResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e DeleteAuthorResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e DeleteAuthorResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e DeleteAuthorResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e DeleteAuthorResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e DeleteAuthorResponseValidationError) ErrorName() string {
	return "DeleteAuthorResponseValidationError"
}

// Error satisfies the builtin error interface
func (e DeleteAuthorResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sDeleteAuthorResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = DeleteAuthorResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = DeleteAuthorResponseValidationError{}
