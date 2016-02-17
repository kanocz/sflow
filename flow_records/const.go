package flow_records


const (
	TypeRawPacketFlowRecord     = 1
	TypeEthernetFrameFlowRecord = 2
	TypeIpv4FlowRecord          = 3
	TypeIpv6FlowRecord          = 4

	TypeExtendedSwitchFlowRecord     = 1001
	TypeExtendedRouterFlowRecord     = 1002
	TypeExtendedGatewayFlowRecord    = 1003
	TypeExtendedUserFlowRecord       = 1004
	TypeExtendedUrlFlowRecord        = 1005
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
)