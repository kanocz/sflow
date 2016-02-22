package records

// sflow flow sample types
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

// Header Protocol Types found in Raw Packet Flow Record
const (
	HeaderProtocolEthernetISO8023   = 1
	HeaderProtocolISO88024Tokenbus  = 2
	HeaderProtocolISO88024Tokenring = 3
	HeaderProtocolFDDI              = 4
	HeaderProtocolFrameRelay        = 5
	HeaderProtocolX24               = 6
	HeaderProtocolPPP               = 7
	HeaderProtocolSMDS              = 8
	HeaderProtocolAAL5              = 9
	HeaderProtocolAAL5IP            = 10
	HeaderProtocolIPv4              = 11
	HeaderProtocolIPv6              = 12
)

// IP Header Protocol Types (see: https://en.wikipedia.org/wiki/List_of_IP_protocol_numbers)
const (
	IPProtocolICMP = 1
	IPProtocolTCP  = 6
	IPProtocolUDP  = 17
)

// Raw Packet Header Types
const (
	HeaderTypeIPv4 = "0800"
	HeaderTypeIPv6 = "86DD"
	//IPX: type_len == 0x0200 || type_len == 0x0201 || type_len == 0x0600
)
