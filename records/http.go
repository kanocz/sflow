package records

import (
	"io"
)

// HTTP Request Types
const (
	HTTPOther   = 0
	HTTPOptions = 1
	HTTPGet     = 2
	HTTPHead    = 3
	HTTPPost    = 4
	HTTPPut     = 5
	HTTPDelete  = 6
	HTTPTrace   = 7
	HTTPConnect = 8
)

// HTTPRequestFlow - TypeHTTPRequestFlowRecord
type HTTPRequestFlow struct {
	Method       uint32
	Protocol     uint32 /* HTTP protocol version: Encoded as major_number * 1000 + minor_number. e.g. HTTP1.1 is encoded as 1001 */
	URILen       uint32
	URI          []byte `xdr:"opaque,lengthField=URILen"` /* URI exactly as it came from the client */
	HostLen      uint32
	Host         []byte `xdr:"opaque,lengthField=HostLen"` /* Host value from request header */
	RefererLen   uint32
	Referer      []byte `xdr:"opaque,lengthField=RefererLen"` /* Referer value from request header */
	UserAgentLen uint32
	UserAgent    []byte `xdr:"opaque,lengthField=UserAgentLen"` /* User-Agent value from request header */
	XFFLen       uint32
	XFF          []byte `xdr:"opaque,lengthField=XFFLen"` /* X-Forwarded-For value from request header */
	AuthUserLen  uint32
	AuthUser     []byte `xdr:"opaque,lengthField=AuthUserLen"` /* RFC 1413 identity of user*/
	MimeTypeLen  uint32
	MimeType     []byte `xdr:"opaque,lengthField=MimeTypeLen"` /* Mime-Type of response */
	ReqBytes     uint64 /* Content-Length of request */
	RespBytes    uint64 /* Content-Length of response */
	Duration     uint32 /* duration of the operation (in microseconds) */
	Status       uint32 /* HTTP status code */
}

// RecordName returns the Name of this flow record
func (f HTTPRequestFlow) RecordName() string {
	return "HTTPRequestFlow"
}

// RecordType returns the ID of the sflow flow record
func (f HTTPRequestFlow) RecordType() int {
	return TypeHTTPRequestFlowRecord
}

func (f HTTPRequestFlow) Encode(w io.Writer) error {
	var err error

	return err
}

// HTTPCounters - TypeHTTPCounterRecord
type HTTPCounter struct {
	MethodOptionCount  uint32
	MethodGetCount     uint32
	MethodHeadCount    uint32
	MethodPostCount    uint32
	MethodPutCount     uint32
	MethodDeleteCount  uint32
	MethodTraceCount   uint32
	MethodConnectCount uint32
	MethodOtherCount   uint32
	Status1XXCount     uint32
	Status2XXCount     uint32
	Status3XXCount     uint32
	Status4XXCount     uint32
	Status5XXCount     uint32
	StatusOtherCount   uint32
}

// RecordName returns the Name of this flow record
func (f HTTPCounter) RecordName() string {
	return "HTTPCounter"
}

// RecordType returns the ID of the sflow flow record
func (f HTTPCounter) RecordType() int {
	return TypeHTTPCounterRecord
}

func (f HTTPCounter) Encode(w io.Writer) error {
	var err error

	return err
}

// ExtendedProxyRequest - TypeHTTPExtendedProxyFlowRecord
type ExtendedProxyRequestFlow struct {
	//string<255> uri;           /* URI in request to downstream server */
	//string<64>  host;          /* Host in request to downstream server */
}
