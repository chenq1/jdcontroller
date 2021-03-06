package ofp11

import (
	"jd.com/jdcontroller/lib/buffer"
	"jd.com/jdcontroller/lib/packet/eth"
	_ "net"
)

/* OpenFlow Switch Specification      */
/* Version 1.1.0 (Wire Protocol 0x02) */
/* February 28, 2011                  */

const (
	Version = 2
)

/* A.1 OpenFlow Header */
type Header struct {
	Version uint8
	Type    uint8
	Length  uint16
	XID     uint32
}

// Type
const (
	// Immutable messages.
	OFPTHello = iota
	OFPTError
	OFPTEchoRequest
	OFPTEchoReply
	OFPTExperimenter

	// Switch configuration messages.
	OFPTFeaturesRequest
	OFPTFeaturesReply
	OFPTGetConfigRequest
	OFPTGetConfigReply
	OFPTSetConfig

	// Asynchronous messages.
	OFPTPacketIn
	OFPTFlowRemoved
	OFPTPortStatus

	// Controller command messages.
	OFPTPacketOut
	OFPTFlowMod
	OFPTGroupMod
	OFPTPortMod
	OFPTTableMod

	// Statistics messages.
	OFPTStatsRequest
	OFPTStatsReply

	// Barrier messages.
	OFPTBarrierRequest
	OFPTBarrierReply

	// Queue Configuration messages.
	OFPTQueueGetConfigRequest
	OFPTQueueGetConfigReply
)

/* A.2 Common Structures */
// constant
const (
	OFPEthALen        = 6
	OFPMaxPortNameLen = 16
)

// A.2.1 Port Structures
type Port struct {
	PortNO uint32
	pad    [4]byte // Size 4
	HWAddr [OFPEthALen]byte
	pad2   [2]byte                 // Size 2
	Name   [OFPMaxPortNameLen]byte // Size 16

	Config uint32
	State  uint32

	Curr       uint32
	Advertised uint32
	Supported  uint32
	Peer       uint32

	CurrSpeed uint32
	MaxSpeed  uint32
}

// Port config
const (
	OFPPCPortDown   = 1 << 0
	OFPPCNORecv     = 1 << 2
	OFPPCNOFWD      = 1 << 5
	OFPPCNOPacketIn = 1 << 6
)

// Port state
const (
	OFPPSLinkDown = 1 << 0
	OFPPSBlocked  = 1 << 1
	OFPPSLive     = 1 << 2
)

// Port numbering
const (
	OFPPMax = 0xffffff00

	OFPPInPort = 0xfffffff8
	OFPPTable  = 0xfffffff9

	OFPPNormal = 0xfffffffa
	OFPPFlood  = 0xfffffffb

	OFPPAll        = 0xfffffffc
	OFPPController = 0xfffffffd
	OFPPLocal      = 0xfffffffe
	OFPPAny        = 0xffffffff
)

// Port features
const (
	OFPPF10MBHD  = 1 << 0
	OFPPF10MBFD  = 1 << 1
	OFPPF100MBHD = 1 << 2
	OFPPF100MBFD = 1 << 3
	OFPPF1GBHD   = 1 << 4
	OFPPF1GBFD   = 1 << 5
	OFPPF10GBFD  = 1 << 6
	OFPPF40GBFD  = 1 << 7
	OFPPF100GBFD = 1 << 8
	OFPPF1TBFD   = 1 << 9
	OFPPFOther   = 1 << 10

	OFPPFCopper    = 1 << 11
	OFPPFFiber     = 1 << 12
	OFPPFAutoneg   = 1 << 13
	OFPPFPause     = 1 << 14
	OFPPFPauseAsym = 1 << 15
)

// A.2.2 Queue Structures
type PacketQueue struct {
	QueueID    uint32
	Length     uint16
	pad        [2]byte // Size 2
	Properties []QueueProp
}

// Queue properties
const (
	OFPQTNone = iota
	OFPQTMinRate
)

type QueueProp interface {
	PackBinary() (data []byte, err error)
	UnpackBinary(data []byte) (err error)
	Len() int
}

type QueuePropHeader struct {
	Property uint16
	Length   uint16
	pad      [4]byte // Size 4
}

type QueuePropMinRate struct {
	Header QueuePropHeader
	Rate   uint16
	pad    [6]byte // Size 6
}

// A.2.3 Flow Match Structures
// constat
const (
	OFPMTStandardLen = 88
)

// Match type
const (
	OFPMTStandard = iota
)

type Match struct {
	Type         uint16
	Length       uint16
	InPort       uint32            /* Input switch port. */
	Wildcards    uint32            /* Wildcard fields. */
	DLSrc        [OFPEthALen]uint8 /* Ethernet source address. */
	DLSrcMask    [OFPEthALen]uint8 /* Ethernet source address. */
	DLDst        [OFPEthALen]uint8 /* Ethernet destination address. */
	DLDstMask    [OFPEthALen]uint8 /* Ethernet destination address. */
	DLVLAN       uint16            /* Input VLAN id. */
	DLVLANPCP    uint8             /* Input VLAN priority. */
	pad          [1]byte           /* Align to 64-bits Size 1 */
	DLType       uint16            /* Ethernet frame type. */
	NWTos        uint8             /* IP ToS (actually DSCP field, 6 bits). */
	NWProto      uint8             /* IP protocol or lower 8 bits of ARP opcode. */
	NWSrc        [4]byte           /* IP source address. */
	NWSrcMask    [4]byte           /* IP source address. */
	NWDst        [4]byte           /* IP destination address. */
	NWDstMask    [4]byte           /* IP destination address. */
	TPSrc        uint16            /* TCP/UDP source port. */
	TPDst        uint16            /* TCP/UDP destination port. */
	MPLSLabel    uint32            /* TCP/UDP destination port. */
	MPLSTC       uint8             /* TCP/UDP destination port. */
	pad2         [3]byte           /* Align to 64-bits Size 3 */
	Metadata     uint64
	MetadataMask uint64
}

// Flow wildcards
const (
	OFPFWInPort    = 1 << 0
	OFPFWDLVLAN    = 1 << 1
	OFPFWDLVLANPCP = 1 << 2
	OFPFWDLType    = 1 << 3
	OFPFWNWTOS     = 1 << 4
	OFPFWNWProto   = 1 << 5
	OFPFWTPSrc     = 1 << 6
	OFPFWTPDst     = 1 << 7
	OFPFWMPLSLabel = 1 << 8
	OFPFWMPLSTC    = 1 << 9

	OFPFWAll = ((1 << 10) - 1)
)

// VLAN id
const (
	OFPVIDAny  = 0xfffe
	OFPVIDNone = 0xffff
)

// A.2.4 Flow Instruction Structures
// Instruction type
const (
	OFPITGotoTable     = 1
	OFPITWriteMetadata = 2
	OFPITWriteActions  = 3
	OFPITApplyActions  = 4
	OFPITClearActions  = 5
	OFPITExperimenter  = 0xffff
)

type Instruction interface {
	PackBinary() (data []byte, err error)
	UnpackBinary(data []byte) (err error)
	Len() int
}

type InstructionHeader struct {
	Type   uint16
	Length uint16
}

type InstructionGotoTable struct {
	Header  InstructionHeader
	TableID uint8
	pad     [3]byte // Size 3
}

type InstructionWriteMetadata struct {
	Header       InstructionHeader
	pad          [4]byte // Size 4
	Metadata     uint64
	MetadataMask uint64
}

type InstructionActions struct {
	Header  InstructionHeader
	pad     [4]byte // Size 4
	Actions []Action
}

// A.2.5 Action Structures
// Action type
const (
	OFPATOutput       = iota // Output to switch port.
	OFPATSetVLANVID          // Set the 802.1q VLAN id.
	OFPATSetVLANPCP          // Set the 802.1q priority.
	OFPATSetDLSrc            // Set ethernet source address.
	OFPATSetDLDst            // Set ethernet destination address.
	OFPATSetNWSrc            // Set IP source address.
	OFPATSetNWDst            // Set IP destination address.
	OFPATSetNWTos            // Set IP ToS (DSCP field, 6bits).
	OFPATSetNWECN            // Set IP ToS (DSCP field, 6bits).
	OFPATSetTPSrc            // Set TCP/UDP source port.
	OFPATSetTPDst            // Set TCP/UDP destination port.
	OFPATCopyTTLOut          // Set TCP/UDP destination port.
	OFPATCopyTTLIn           // Set TCP/UDP destination port.
	OFPATSetMPLSLabel        // Set TCP/UDP destination port.
	OFPATSetMPLSTC           // Set TCP/UDP destination port.
	OFPATSetMPLSTTL          // Set TCP/UDP destination port.
	OFPATDecMPLSTTL          // Set TCP/UDP destination port.

	OFPATPushVLAN // Set TCP/UDP destination port.
	OFPATPopVLAN  // Set TCP/UDP destination port.
	OFPATPushMPLS // Set TCP/UDP destination port.
	OFPATPopMPLS  // Set TCP/UDP destination port.

	OFPATSetQueue // Set TCP/UDP destination port.
	OFPATGroup    // Set TCP/UDP destination port.
	OFPATSetNWTTL // Set TCP/UDP destination port.
	OFPATDecNWTTL // Set TCP/UDP destination port.

	OFPATExperimenter = 0xffff
)

type Action interface {
	PackBinary() (data []byte, err error)
	UnpackBinary(data []byte) (err error)
	Len() int
}

type ActionHeader struct {
	Type   uint16
	Length uint16
}

type ActionOutput struct {
	Header ActionHeader
	Port   uint32
	MaxLen uint16
	pad    [6]byte // Size 6
}

type ActionGroup struct {
	Header  ActionHeader
	GroupID uint32
}

type ActionSetQueue struct {
	Header  ActionHeader
	QueueID uint32
}

type ActionVLANVID struct {
	Header  ActionHeader
	VLANVID uint16
	pad     [2]byte // Size 2
}

type ActionVLANPCP struct {
	Header  ActionHeader
	VLANPCP uint8
	pad     [3]byte // Size 3
}

type ActionMPLSLabel struct {
	Header    ActionHeader
	MPLSLabel uint32
}

type ActionMPLSTC struct {
	Header ActionHeader
	MPLSTC uint8
	pad    [3]byte // Size 3
}

type ActionMPLSTTL struct {
	Header  ActionHeader
	MPLSTTL uint8
	pad     [3]byte // Size 3
}

type ActionDLAddr struct {
	Header ActionHeader
	DLAddr [OFPEthALen]uint8
	pad    [6]byte // Size 6
}

type ActionNWAddr struct {
	Header ActionHeader
	NWAddr [4]byte
}

type ActionNWTOS struct {
	Header ActionHeader
	NWTOS  uint8
	pad    [3]byte // Size 3
}

type ActionNWECN struct {
	Header ActionHeader
	NWECN  uint8
	pad    [3]byte // Size 3
}

type ActionNWTTL struct {
	Header ActionHeader
	NWTTL  uint8
	pad    [3]byte // Size 3
}

type ActionTPPort struct {
	Header ActionHeader
	TPPort uint16
	pad    [2]byte // Size 2
}

type ActionPush struct {
	Header    ActionHeader
	Ethertype uint16
	pad       [2]byte // Size 2
}

type ActionPopMPLS struct {
	Header    ActionHeader
	Ethertype uint16
	pad       [2]byte // Size 2
}

type ActionExperimenterHeader struct {
	Header       ActionHeader
	Experimenter uint32
}

/* A.3 Controller to Switch Messages */
// A.3.1 Handshake
type SwitchFeatures struct {
	Header       Header
	DPID         uint64
	Buffers      uint32
	Tables       uint8
	pad2         [3]byte // Size 3
	Capabilities uint32
	Reserved     uint32
	Ports        []Port
}

// Capabilities
const (
	OFPCFlowStats  = 1 << 0 // Flow statistics.
	OFPCTableStats = 1 << 1 // Table statistics.
	OFPCPortStats  = 1 << 2 // Port statistics.
	OFPCGroupStats = 1 << 3 // 802.1d spanning tree.
	OFPCIPReasm    = 1 << 5 // Can reassemble IP fragments.
	OFPCQueueStats = 1 << 6 // Queue statistics.
	OFPCARPMatchIP = 1 << 7 // Match IP address in ARP packets.
)

// A.3.2 Switch Configuration
// Switch config
type SwitchConfig struct {
	Header      Header
	Flags       uint16
	MissSendLen uint16
}

// Config flags
const (
	OFPCFragNormal             = 0
	OFPCFragDrop               = 1 << 0
	OFPCFragReasm              = 1 << 1
	OFPCFragMask               = 3
	OFPCInvalidTTLToController = 1 << 2
)

// A.3.3 Flow Table Configuration
// Table mod
type TableMod struct {
	Header  Header
	TableID uint8
	pad     [3]byte // Size 3
	Config  uint32
}

// Table config
const (
	OFPTCTableMissController = 0
	OFPTCTableMissContinue   = 1 << 0
	OFPTCTableMissDrop       = 1 << 1
	OFTCTableMissMask        = 3
)

// A.3.4 Modify State Messages
// Flow mod
type FlowMod struct {
	Header     Header
	Cookie     uint64
	CookieMask uint64

	TableID      uint8
	Command      uint8
	IdleTimeout  uint16
	HardTimeout  uint16
	Priority     uint16
	BufferID     uint32
	OutPort      uint32
	OutGroup     uint32
	Flags        uint16
	pad          [2]byte // Size 2
	Match        Match
	Instructions []Instruction
}

// Flow mod command
const (
	OFPFCAdd          = iota // New flow.
	OFPFCModify              // Modify all matching flow.
	OFPFCModifyStrict        // Modify entry strictly matching wildcards.
	OFPFCDelete              // Delete all matching flow.
	OFPFCDeleteStrict        // Strictly match wildcards and priority.
)

// Flow mod flags
const (
	OFPFFSendFlowRem  = 1 << 0
	OFPFFCheckOverlap = 1 << 1
)

// Group mod
type GroupMod struct {
	Header  Header
	Command uint16
	Type    uint8
	pad     uint8 // change [1]unint8 to uint8
	GroupID uint32
	Buckets []Bucket
}

// Group mod command
const (
	OFPGCAdd = iota
	OFPGCModify
	OFPGCDelete
)

// Group type
const (
	OFPGTAll = iota
	OFPGTSelect
	OFPGTIndirect
	OFPGTFF
)

type Bucket struct {
	Length     uint16
	Weight     uint16
	WatchPort  uint32
	WatchGroup uint32
	pad        [4]byte // Size 4
	Actions    []Action
}

// Port mod
type PortMod struct {
	Header    Header
	PortNO    uint32
	pad       [4]byte // Size 4
	HWAddr    [OFPEthALen]uint8
	pad2      [2]byte // Size 2
	Config    uint32
	Mask      uint32
	Advertise uint32
	pad3      [4]byte // Size 4
}

// A.3.5 Queue Configuration Messages
type QueueGetConfigRequest struct {
	Header Header
	Port   uint32
	pad    [4]byte // Size 4
}

type QueueGetConfigReply struct {
	Header Header
	Port   uint32
	pad    [4]byte // Size 4
	Queues []PacketQueue
}

// A.3.6 Read State Messages
// constats
const (
	OFPDescStrLen      = 256
	OFPSerialNumLen    = 32
	OFPMaxTableNameLen = 32
)

// Stats request
type StatsRequest struct {
	Header Header
	Type   uint16
	Flags  uint16
	pad    [4]byte // Size 4
	Body   buffer.Message
}

// Stats reply
type StatsReply struct {
	Header Header
	Type   uint16
	Flags  uint16
	pad    [4]byte // Size 4
	Body   buffer.Message
}

// Stats types
const (
	OFPSTDesc = iota
	OFPSTFlow
	OFPSTAggregate
	OFPSTTable
	OFPSTPort
	OFPSTQueue
	OFPSTGroup
	OFPSTGroupDesc
	OFPSTExperimenter = 0xffff
)

// Desc stats
type DescStats struct {
	MfrDesc   [OFPDescStrLen]byte   // Size 256
	HWDesc    [OFPDescStrLen]byte   // Size 256
	SWDesc    [OFPDescStrLen]byte   // Size 256
	SerialNum [OFPSerialNumLen]byte // Size 32
	DPDesc    [OFPDescStrLen]byte   // Size 256
}

// Flow stats request
type FlowStatsRequest struct {
	TableID    uint8
	pad        [3]byte // Size 3
	OutPort    uint32
	OutGroup   uint32
	pad2       [4]byte // Size 4
	Cookie     uint64
	CookieMask uint64
	Match      Match
}

// Flow stats
type FlowStats struct {
	Length       uint16
	TableID      uint8
	pad          uint8 // change [1]byte to uint8
	DurationSec  uint32
	DurationNSec uint32
	Priority     uint16
	IdleTimeout  uint16
	HardTimeout  uint16
	pad2         [6]byte // Size 6
	Cookie       uint64
	PacketCount  uint64
	ByteCount    uint64
	Match        Match
	Instructions []Instruction
}

// Aggregate stats request
type AggregateStatsRequest struct {
	TableID    uint8
	pad        [3]byte // Size 3
	OutPort    uint32
	OutGroup   uint32
	pad2       [4]byte // Size 4
	Cookie     uint64
	CookieMask uint64
	Match      Match
}

// Aggregate stats reply
type AggregateStatsReply struct {
	PacketCount uint64
	ByteCount   uint64
	FlowCount   uint32
	pad         [4]byte // Size 4
}

// Table stats
type TableStats struct {
	TableID      uint8
	pad          [7]uint8                 // Size 7
	Name         [OFPMaxTableNameLen]byte // Size 32
	Wildcards    uint32
	Match        uint32
	Instructions uint32
	WriteActions uint32
	ApplyActions uint32
	Config       uint32
	MaxEntries   uint32
	ActiveCount  uint32
	LookupCount  uint64
	MatchedCount uint64
}

// Port stats request
type PortStatsRequest struct {
	PortNO uint32
	pad    [4]byte // Size 4
}

// Port stats
type PortStats struct {
	PortNO     uint32
	pad        [4]byte // Size 4
	RxPackets  uint64
	TxPackets  uint64
	RxBytes    uint64
	TxBytes    uint64
	RxDropped  uint64
	TxDropped  uint64
	RxErrors   uint64
	TxErrors   uint64
	RxFrameErr uint64
	RxOverErr  uint64
	RxCRCErr   uint64
	Collisions uint64
}

// Queue stats request
type QueueStatsRequest struct {
	PortNO  uint32
	QueueID uint32
}

// Queue stats
type QueueStats struct {
	PortNO    uint32
	QueueID   uint32
	TxBytes   uint64
	TxPackets uint64
	TxErrors  uint64
}

// Group stats request
type GroupStatsRequest struct {
	GroupID uint32
	pad     [4]byte // Size 4
}

// Group stats
type GroupStats struct {
	Length      uint16
	pad         [2]byte // Size 2
	GroupID     uint32
	RefCount    uint32
	pad2        [4]byte // Size 4
	PacketCount uint64
	ByteCount   uint64
	BucketStats []BucketCounter
}

type BucketCounter struct {
	PacketCount uint64
	ByteCount   uint64
}

type GroupDescStats struct {
	Length  uint16
	Type    uint8
	pad     uint8 // change [1]byte to uint8
	GroupID uint32
	Buckets []Bucket
}

// A.3.7 Packet-Out Message
type PacketOut struct {
	Header     Header
	BufferID   uint32
	InPort     uint32
	ActionsLen uint16
	pad        [6]byte // Size 6
	Actions    []Action
	Data       buffer.Message
}

/* A.4 Asynchronous Messages */
// A.4.1 Packet-In Message
type PacketIn struct {
	Header    Header
	BufferID  uint32
	InPort    uint32
	InPhyPort uint32
	TotalLen  uint16
	Reason    uint8
	TableID   uint8
	//	Data      buffer.Message
	Data eth.Ethernet
}

// Packet-in reason
const (
	OFPRNoMatch = iota
	OFPRAction
)

// A.4.2 Flow Removed Message
type FlowRemoved struct {
	Header   Header
	Cookie   uint64
	Priority uint16
	Reason   uint8
	TableID  uint8

	DurationSec  uint32
	DurationNSec uint32

	IdleTimeout uint16
	pad         [2]byte // Size 2
	PacketCount uint64
	ByteCount   uint64
	Match       Match
}

// Flow removed reason
const (
	OFPRRIdleTimeout = iota
	OFPRRHardTimeout
	OFPRRDelete
	OFPRRGroupDelete
)

// A.4.3 Port Status Message
type PortStatus struct {
	Header Header
	Reason uint8
	pad    [7]uint8 // Size 7
	Desc   Port
}

// Port reason
const (
	OFPPRAdd = iota
	OFPPRDelete
	OFPPRModify
)

// A.4.4 Error Message
type ErrorMsg struct {
	Header Header
	Type   uint16
	Code   uint16
	Data   []uint8
}

// Error type
const (
	OFPETHelloFailed = iota
	OFPETBadRequest
	OFPETBadAction
	OFPETBadInstruction
	OFPETBadMatch
	OFPETFlowModFailed
	OFPETGroupModFailed
	OFPETPortModFailed
	OFPETTableModFailed
	OFPETQueueOPFailed
	OFPETSwitchConfigFailed
)

// Hello failed code
const (
	OFPHFCIncompatible = iota
	OFPHFCEperm
)

// Bad request code
const (
	OFPBRCBadVersion = iota
	OFPBRCBadType
	OFPBRCBadStat
	OFPBRCBadExperimenter

	OFPBRCBadSubtype
	OFPBRCEperm
	OFPBRCBadLen
	OFPBRCBufferEmpty
	OFPBRCBufferUnknown
	OFPBRCBadTableID
)

// Bad action code
const (
	OFPBACBadType = iota
	OFPBACBadLen
	OFPBACBadExperimenter
	OFPBACBadExperimenterType
	OFPBACBadOutPort
	OFPBACBadArgument
	OFPBACEperm
	OFPBACTooMany
	OFPBACBadQueue
	OFPBACBadOutGroup
	OFPBACMatchInconsistent
	OFPBACUnsupportedOrder
	OFPBACBadTag
)

// Bad instruction code
const (
	OFPBICUnknownInst = iota
	OFPBICUnsupInst
	OFPBICBadTableID
	OFPBICUnsupMetadata
	OFPBICUnsupMetadataMask
	OFPBICUnsupExpInst
)

// Bad Match code
const (
	OFPBMCBadType = iota
	OFPBMCBadLen
	OFPBMCBadTag
	OFPBMCBadDLAddrMask
	OFPBMCBadNWAddrMask
	OFPBMCBadWildcards
	OFPBMCBadField
	OFPBMCBadValue
)

// Flow mod failed code
const (
	OFPFMFCUnknown = iota
	OFPFMFCTableFull
	OFPFMFCBadTableID
	OFPFMFCOverlap
	OFPFMFCEperm
	OFPFMFCBadTimeout
	OFPFMFCBadCommand
)

// Group mod failed code
const (
	OFPGMFCGroupExists = iota
	OFPGMFCInvalidGroup
	OFPFMFCWeightUnsupported
	OFPFMFCOutOfGroups
	OFPFMFCOutOfBuckets
	OFPFMFCChainingUnsupported
	OFPFMFCWatchUnsupported
	OFPFMFCLoop
	OFPFMFCUnknownGroup
)

// Port mod failed code
const (
	OFPPMFCBadPort = iota
	OFPPMFCBadHWAddr
	OFPPMFCBadConfig
	OFPPMFCBadAdvertise
)

// Table mod failed code
const (
	OFPTMFCBadTable = iota
	OFPTMFCBadConfig
)

// Queue op failed code
const (
	OFPQOFCBadPort = iota
	OFPQOFCBadQueue
	OFPQOFCEperm
)

// Switch config failed code
const (
	OFPSCFCBadFlags = iota
	OFPSCFCBadLen
)

/* A.5 Symmetric Messages */
type ExperimenterHeader struct {
	Header       Header /*Type OFPT_VENDOR*/
	Experimenter uint32
	pad          [4]byte
}
