package records

import (
	"testing"
)

func TestCalculateBinarySizeExtendedSwitchFlow(t *testing.T) {
	rec := ExtendedSwitchFlow{
		SourceVlan:          1234,
		SourcePriority:      15,
		DestinationVlan:     4321,
		DestinationPriority: 1,
	}

	size := rec.calculateBinarySize()
	if size != 16 {
		t.Errorf("expected\n%+#v\n, got\n%+#v", 76, size)
	}
}
