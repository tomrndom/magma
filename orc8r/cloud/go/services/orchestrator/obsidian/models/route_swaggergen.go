// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/swag"
)

// Route route
// swagger:model route
type Route struct {

	// destination ip
	DestinationIP string `json:"destination_ip,omitempty" magma_alt_name:"DestinationIp"`

	// gateway ip
	GatewayIP string `json:"gateway_ip,omitempty" magma_alt_name:"GatewayIp"`

	// genmask
	Genmask string `json:"genmask,omitempty"`

	// network interface id
	NetworkInterfaceID string `json:"network_interface_id,omitempty" magma_alt_name:"NetworkInterfaceId"`
}

// Validate validates this route
func (m *Route) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *Route) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Route) UnmarshalBinary(b []byte) error {
	var res Route
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
