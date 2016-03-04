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
	//HTTP_method method;        /* method */
	//version protocol;          /* HTTP protocol version */
	//string<255> uri;           /* URI exactly as it came from the client */
	//string<64> host;           /* Host value from request header */
	//string<255> referer;       /* Referer value from request header */
	//string<128> useragent;     /* User-Agent value from request header */
	//string<64> xff;            /* X-Forwarded-For value from request header */
	//string<32> authuser;       /* RFC 1413 identity of user*/
	//string<64> mime-type;      /* Mime-Type of response */
	//unsigned hyper req_bytes;  /* Content-Length of request */
	//unsigned hyper resp_bytes; /* Content-Length of response */
	//unsigned int uS;           /* duration of the operation (in microseconds) */
	//int status;                /* HTTP status code */
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
