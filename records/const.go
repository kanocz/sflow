package records

// sflow flow record types
const (
	TypeRawPacketFlowRecord     = 1
	TypeEthernetFrameFlowRecord = 2
	TypeIpv4FlowRecord          = 3
	TypeIpv6FlowRecord          = 4

	TypeExtendedSwitchFlowRecord     = 1001
	TypeExtendedRouterFlowRecord     = 1002
	TypeExtendedGatewayFlowRecord    = 1003
	TypeExtendedUserFlowRecord       = 1004
	TypeExtendedURLFlowRecord        = 1005
	TypeExtendedMlpsFlowRecord       = 1006
	TypeExtendedNatFlowRecord        = 1007
	TypeExtendedMlpsTunnelFlowRecord = 1008
	TypeExtendedMlpsVcFlowRecord     = 1009
	TypeExtendedMlpsFecFlowRecord    = 1010
	TypeExtendedMlpsLvpFecFlowRecord = 1011
	TypeExtendedVlanFlowRecord       = 1012

	TypeExtendedSocketIPv4FlowRecord      = 2102
	TypeExtendedSocketIPv6FlowRecord      = 2103
	TypeExtendedProxySocketIPv4FlowRecord = 2104
	TypeExtendedProxySocketIPv6FlowRecord = 2105
	TypeHTTPRequestFlowRecord             = 2206
	TypeHTTPExtendedProxyFlowRecord       = 2207
)

// flow sample record data structure mapping
var flowRecordTypes = map[uint32]interface{}{
	TypeRawPacketFlowRecord:       RawPacketFlow{},
	TypeEthernetFrameFlowRecord:   EthernetFrameFlow{},
	TypeExtendedSwitchFlowRecord:  ExtendedSwitchFlow{},
	TypeExtendedRouterFlowRecord:  ExtendedRouterFlow{},
	TypeExtendedGatewayFlowRecord: ExtendedGatewayFlow{},
}

// sflow counter record types
const (
	TypeHTTPCounterRecord = 2201
)

// IP Header Protocol Types (see: https://en.wikipedia.org/wiki/List_of_IP_protocol_numbers)
const (
	IPProtocolICMP = 1
	IPProtocolTCP  = 6
	IPProtocolUDP  = 17
	IPProtocolESP  = 50 // IPSEC
	IPProtocolAH   = 51 // IPSEC
)

const (
	// MaximumRecordLength defines the maximum length acceptable for decoded records.
	// This maximum prevents from excessive memory allocation.
	// The value is derived from MAX_PKT_SIZ 65536 in the reference sFlow implementation
	// https://github.com/sflow/sflowtool/blob/bd3df6e11bdf/src/sflowtool.c#L4313.
	MaximumRecordLength = 65536

	// MaximumHeaderLength defines the maximum length acceptable for decoded flow samples.
	// This maximum prevents from excessive memory allocation.
	// The value is set to maximum transmission unit (MTU), as the header of a network packet
	// may not exceed the MTU.
	MaximumHeaderLength = 1500

	// MinimumEthernetHeaderSize defines the minimum header size to be parsed
	MinimumEthernetHeaderSize = 14
	//#define NFT_8022_SIZ 3
	//#define NFT_MAX_8023_LEN 1500
	//#define NFT_MIN_SIZ (NFT_ETHHDR_SIZ + sizeof(struct myiphdr))
)
