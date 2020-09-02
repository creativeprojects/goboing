package main

import "testing"

// quick test to make sure the byte slice is not copied but only references in the two lists
// and this is right, there's only one allocation per cycle:
// BenchmarkByteSliceAllocation-12    	22749994	        45.8 ns/op	     160 B/op	       1 allocs/op
func BenchmarkByteSliceAllocation(b *testing.B) {
	b.ReportAllocs()
	list1 := make([][]byte, b.N)
	list2 := make([][]byte, b.N)
	for i := 0; i < b.N; i++ {
		buffer := make([]byte, 100)
		list1[i] = buffer
		list2[i] = buffer
	}
}
