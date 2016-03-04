package records

/*import (
	"io"
)

type HostDescriptionCounter struct {
	HostnameLen uint32
	Hostname []byte `lengthLookUp:"HostnameLen"`
	UUID [16]byte
	MachineType uint32
	OSName uint32
	OSReleaseLen uint32
	OSRelease []byte `lengthLookUp:"OSReleaseLen"`
}

// RecordName returns the Name of this flow record
func (f HostDescriptionCounter) RecordName() string {
	return "HostDescriptionCounter"
}

// RecordType returns the ID of the sflow flow record
func (f HostDescriptionCounter) RecordType() int {
	return TypeHostDescriptionCounterRecord
}

func (f HostDescriptionCounter) Encode(w io.Writer) error {
	var err error

	return err
}*/

//enum machine_type {
//unknown = 0,
//other   = 1,
//x86     = 2,
//x86_64  = 3,
//ia64    = 4,
//sparc   = 5,
//alpha   = 6,
//powerpc = 7,
//m68k    = 8,
//mips    = 9,
//arm     = 10,
//hppa    = 11,
//s390    = 12
//}
//
///* The os_name enumeration may be expanded over time.
//   Applications receiving sFlow must be prepared to receive
//   host_descr structures with unknown machine_type values.
//
//   The authoritative list of machine types will be maintained
//   at www.sflow.org */
//
//enum os_name {
//unknown   = 0,
//other     = 1,
//
//
//
//FINAL                           sFlow.org                       [Page 7]
//
//FINAL                     sFlow Host Structures                July 2010
//
//
//linux     = 2,
//windows   = 3,
//darwin    = 4,
//hpux      = 5,
//aix       = 6,
//dragonfly = 7,
//freebsd   = 8,
//netbsd    = 9,
//openbsd   = 10,
//osf       = 11,
//solaris   = 12
//}
//
///* Physical or virtual host description */
///* opaque = counter_data; enterprise = 0; format = 2000 */
//struct host_descr {
//string hostname<64>;       /* hostname, empty if unknown */
//opaque uuid<16>;           /* 16 bytes binary UUID, empty if unknown */
//machine_type machine_type; /* the processor family */
//os_name os_name;           /* Operating system */
//string os_release<32>;     /* e.g. 2.6.9-42.ELsmp,xp-sp3, empty if unknown */
//}
