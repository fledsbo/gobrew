package fermentation

import "testing"

// TestAverager tests the averager code
func TestAverager(t *testing.T) {
	a := averager{}
	ret := a.addReading(1.0)
	if ret != 1.0 {
		t.Errorf("Result should be 1.0, was %f", ret)
	}
	a.addReading(1.0)
	if ret != 1.0 {
		t.Errorf("Result should be 1.0, was %f", ret)
	}
	a.addReading(1.0)
	a.addReading(1.0)
	a.addReading(10000)
	ret = a.addReading(0)
	if ret != 1.0 {
		t.Errorf("Result should be 1.0, was %f", ret)
	}

}
