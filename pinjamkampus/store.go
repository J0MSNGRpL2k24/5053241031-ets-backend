package main

import "sync"

type MemoryStore struct {
	mu               sync.Mutex
	SumberDaya       map[int]SumberDaya  
	Peminjaman       map[int]Peminjaman  
	nextSumberDayaID int
	nextPeminjamanID int
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		SumberDaya:       make(map[int]SumberDaya), 
		Peminjaman:       make(map[int]Peminjaman), 
		nextSumberDayaID: 1,
		nextPeminjamanID: 1,
	}
}