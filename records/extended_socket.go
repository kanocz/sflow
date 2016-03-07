package records

import (
	"net"
)

// ExtendedSocketIPv4Flow - TypeExtendedSocketIPv4FlowRecord
type ExtendedSocketIPv4Flow struct {
	Protocol   uint32
	LocalIP    net.IP `ipVersion:"4"`
	RemoteIP   net.IP `ipVersion:"4"`
	LocalPort  uint32
	RemotePort uint32
}

// RecordName returns the Name of this flow record
func (f ExtendedSocketIPv4Flow) RecordName() string {
	return "ExtendedSocketIPv4Flow"
}

// RecordType returns the ID of the sflow flow record
func (f ExtendedSocketIPv4Flow) RecordType() int {
	return TypeExtendedSocketIPv4FlowRecord
}

// ExtendedSocketIPv6Flow - TypeExtendedSocketIPv6FlowRecord
type ExtendedSocketIPv6Flow struct {
	Protocol   uint32
	LocalIP    net.IP `ipVersion:"6"`
	RemoteIP   net.IP `ipVersion:"6"`
	LocalPort  uint32
	RemotePort uint32
}

// RecordName returns the Name of this flow record
func (f ExtendedSocketIPv6Flow) RecordName() string {
	return "ExtendedSocketIPv6Flow"
}

// RecordType returns the ID of the sflow flow record
func (f ExtendedSocketIPv6Flow) RecordType() int {
	return TypeExtendedSocketIPv6FlowRecord
}

// ExtendedProxySocketIPv4Flow - TypeExtendedProxySocketIPv4FlowRecord
type ExtendedProxySocketIPv4Flow struct {
	Socket ExtendedSocketIPv4Flow
}

// RecordName returns the Name of this flow record
func (f ExtendedProxySocketIPv4Flow) RecordName() string {
	return "ExtendedProxySocketIPv4Flow"
}

// RecordType returns the ID of the sflow flow record
func (f ExtendedProxySocketIPv4Flow) RecordType() int {
	return TypeExtendedProxySocketIPv4FlowRecord
}

// ExtendedProxySocketIPv6Flow - TypeExtendedProxySocketIPv6FlowRecord
type ExtendedProxySocketIPv6Flow struct {
	Socket ExtendedSocketIPv6Flow
}

// RecordName returns the Name of this flow record
func (f ExtendedProxySocketIPv6Flow) RecordName() string {
	return "ExtendedProxySocketIPv6Flow"
}

// RecordType returns the ID of the sflow flow record
func (f ExtendedProxySocketIPv6Flow) RecordType() int {
	return TypeExtendedProxySocketIPv6FlowRecord
}
