package records

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
	/*
			  unsigned int method_option_count;
		  unsigned int method_get_count;
		  unsigned int method_head_count;
		  unsigned int method_post_count;
		  unsigned int method_put_count;
		  unsigned int method_delete_count;
		  unsigned int method_trace_count;
		  unsigned int method_connect_count;
		  unsigned int method_other_count;
		  unsigned int status_1XX_count;
		  unsigned int status_2XX_count;
		  unsigned int status_3XX_count;
		  unsigned int status_4XX_count;
		  unsigned int status_5XX_count;
		  unsigned int status_other_count;
	*/
}

// ExtendedProxyRequest - TypeHTTPExtendedProxyFlowRecord
type ExtendedProxyRequestFlow struct {
	//string<255> uri;           /* URI in request to downstream server */
	//string<64>  host;          /* Host in request to downstream server */
}
