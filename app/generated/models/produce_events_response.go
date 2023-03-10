// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// ProduceEventsResponse produce events response
//
// swagger:model ProduceEventsResponse
type ProduceEventsResponse struct {

	// status
	// Enum: [OK]
	Status string `json:"status,omitempty"`
}

// Validate validates this produce events response
func (m *ProduceEventsResponse) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateStatus(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

var produceEventsResponseTypeStatusPropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["OK"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		produceEventsResponseTypeStatusPropEnum = append(produceEventsResponseTypeStatusPropEnum, v)
	}
}

const (

	// ProduceEventsResponseStatusOK captures enum value "OK"
	ProduceEventsResponseStatusOK string = "OK"
)

// prop value enum
func (m *ProduceEventsResponse) validateStatusEnum(path, location string, value string) error {
	if err := validate.EnumCase(path, location, value, produceEventsResponseTypeStatusPropEnum, true); err != nil {
		return err
	}
	return nil
}

func (m *ProduceEventsResponse) validateStatus(formats strfmt.Registry) error {

	if swag.IsZero(m.Status) { // not required
		return nil
	}

	// value enum
	if err := m.validateStatusEnum("status", "body", m.Status); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *ProduceEventsResponse) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ProduceEventsResponse) UnmarshalBinary(b []byte) error {
	var res ProduceEventsResponse
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
