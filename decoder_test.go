package main

import "testing"

const sample = `/firstGame/Android/anzhi002/0.1.0.60_0.0.0.0.2/210/3ff4e3f1c3a71cd99bceeb891577b2fc2/anzhi_201501171417022I8S29P7dT/17/21537/MissionCompleted/12_3`

func TestDecode_func(t *testing.T) {
	v := decode_func(sample)
	t.Log(v)
}
