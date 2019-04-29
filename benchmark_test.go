package goradix

import (
	"log"
	"math/rand"
	"runtime"
	"testing"
	"unsafe"
)

const keycount = 1000000 // 1M
// const keycount = 10000000 // 10M; nodes: 418,558,156
// const keycount = 50000000 // 50M; nodes: 1,976,725,806
// const keycount = 100000000 // 100M
// const keycount = 200000000 // 200M

var checkValues, insertValues []uint64
var m map[uint64]struct{}
var bravo *Bravo
var charlie *Charlie

func init() {
	checkValues = make([]uint64, keycount)
	for i := 0; i < keycount; i++ {
		checkValues[i] = rand.Uint64()
	}

	insertValues = make([]uint64, keycount)
	for i := 0; i < keycount; i++ {
		insertValues[i] = rand.Uint64()
	}
}

func getHeapAlloc() uint64 {
	runtime.GC()
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	return ms.HeapAlloc
}

func BenchmarkMap(b *testing.B) {
	if m == nil {
		log.Printf("building map")
		m = make(map[uint64]struct{})
		for i := 0; i < keycount; i++ {
			m[insertValues[i]] = struct{}{}
		}
		b.ResetTimer()
	}
	for i := 0; i < b.N; i++ {
		_ = m[checkValues[i%keycount]]
	}
}

func BenchmarkBravo(b *testing.B) {
	if bravo == nil {
		log.Printf("building bravo")
		bravo = NewBravo()
		for i := 0; i < keycount; i++ {
			bravo.Store(insertValues[i])
		}
		log.Printf("keycount: %d; bravo.nodes: %d; footprint: %d", keycount, bravo.Count, bravo.Count*16)
		b.ResetTimer()
	}
	for i := 0; i < b.N; i++ {
		_ = bravo.Load(checkValues[i%keycount])
	}
}

func BenchmarkCharlie4(b *testing.B) {
	b.Skip()
	bits := uint8(4)
	if charlie == nil {
		log.Printf("building charlie")
		charlie = NewCharlie(bits)
		for i := 0; i < keycount; i++ {
			charlie.Store(insertValues[i])
		}
		log.Printf("keycount: %d; charlie.nodes: %d; footprint: %d", keycount, charlie.Count, charlie.Sizeof())
		b.ResetTimer()
	}
	for i := 0; i < b.N; i++ {
		_ = charlie.Load(checkValues[i%keycount])
	}
	charlie = nil
}

func BenchmarkCharlie8(b *testing.B) {
	bits := uint8(8)
	if charlie == nil {
		log.Printf("building charlie")
		charlie = NewCharlie(bits)
		for i := 0; i < keycount; i++ {
			charlie.Store(insertValues[i])
		}
		log.Printf("keycount: %d; charlie.nodes: %d; footprint: %d", keycount, charlie.Count, charlie.Sizeof())
		b.ResetTimer()
	}
	for i := 0; i < b.N; i++ {
		_ = charlie.Load(checkValues[i%keycount])
	}
	charlie = nil
}

func (charlie *Charlie) Sizeof() uint64 {
	var f *cnode
	sizePointer := uint64(unsafe.Sizeof(f))
	sizePointers := charlie.childCount * sizePointer
	sizeNode := uint64(unsafe.Sizeof(charlie.head)) + sizePointers
	sizeNodes := charlie.Count * sizeNode
	return uint64(unsafe.Sizeof(charlie)) + sizeNodes
}

func BenchmarkCharlie16(b *testing.B) {
	bits := uint8(16)
	if charlie == nil {
		log.Printf("building charlie")
		charlie = NewCharlie(bits)
		for i := 0; i < keycount; i++ {
			charlie.Store(insertValues[i])
		}
		log.Printf("keycount: %d; charlie.nodes: %d; footprint: %d", keycount, charlie.Count, charlie.Sizeof())
		b.ResetTimer()
	}
	for i := 0; i < b.N; i++ {
		_ = charlie.Load(checkValues[i%keycount])
	}
	charlie = nil
}
