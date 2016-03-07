package records

import (
	"net"
	"testing"
)

func TestCalculateBinarySizeExtendedRouterFlow(t *testing.T) {
	rec := ExtendedRouterFlow{
		NextHopType: 2,
		NextHop:     net.ParseIP("2001:0db8:ac10:fe01::"),
		SrcMask:     23,
		DstMask:     1,
	}

	size := rec.BinarySize()
	if size != 28 {
		t.Errorf("expected\n%+#v\n, got\n%+#v", 76, size)
	}
}
