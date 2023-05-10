package main

import (
	"encoding/hex"
	"testing"
)

type Results struct {
	reqType string
	x       uint32
	y       uint32
}

func TestDeserializeMsg(t *testing.T) {
	samples := map[string]Results{
		"490000303900000065": {"I", 12345, 101},
		"490000303a00000066": {"I", 12346, 102},
		"490000303b00000064": {"I", 12347, 100},
		"490000a00000000005": {"I", 40960, 5},
		"510000300000004000": {"Q", 12288, 16384},
	}
	for h, r := range samples {
		sample, err := hex.DecodeString(h)
		if err != nil {
			t.Errorf("Invalid test sample %v. Invalid hex sample", h)
		}
		if len(sample) != 9 {
			t.Errorf("Invalid test sample %v. Invalid hex sample", h)
		}
		l, x, y, err := deserializeMsg(sample)
		if err != nil {
			t.Error("deserializeMsg failed when it shouldn't: ", err)
		}
		if *l != r.reqType || *x != r.x || *y != r.y {
			t.Errorf("Wrong deserialze function (%s,%d,%d) should be equal to (%s,%d,%d)", *l, *x, *y, r.reqType, r.x, r.y)
		}
	}
}

func TestMeansToAnEnd02(t *testing.T) {
}
