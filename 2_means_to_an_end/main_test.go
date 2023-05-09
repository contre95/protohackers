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
	sample1, _ := hex.DecodeString("490000303900000065") //  I 12345 101
	sample2, _ := hex.DecodeString("490000303a00000066") //  I 12346 102
	sample3, _ := hex.DecodeString("490000303b00000064") //I12347100
	sample4, _ := hex.DecodeString("490000a00000000005") //I409605
	sample5, _ := hex.DecodeString("510000300000004000") //Q1228816384

    samples := map[string]Results{sample1: Results{"I"}, sample2, sample3, sample4, sample5}
	for i, sample := range samples {
		if len(sample) != 9 {
			t.Errorf("Invalid test sample %d. Invalid hex sample", i)
		}
	}
	l, x, y, err := deserializeMsg(sample)
	if err != nil {
		t.Error("deserializeMsg failed when it shouldn't: ", err)
	}
	if *l != "I" || *y != 101 || *x != 12345 {
		t.Errorf("Wrong deserialze function (%s,%d,%d) should be equal to (I,12345,101)", *l, *x, *y)
	}
}

// func TestMeansToAnEnd02(t *testing.T) {
// 	sample, err := hex.DecodeString("49 00 00 30 39 00 00 00 65")
// 	db := map[uint32]uint32{}
// 	_, err = meansToAnEnd02(sample, 0, db)
// 	if err != nil {
// 		t.Error("Test failed")
// 	}
// 	// fmt.Println(string(resp))
// }
