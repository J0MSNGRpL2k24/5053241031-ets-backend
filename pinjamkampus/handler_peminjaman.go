package main

import "net/http"

func handleCreatePeminjaman(store *MemoryStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: Implementasi B2 (Buat Peminjaman)
	}
}

func handleGetPeminjaman(store *MemoryStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: Implementasi B3 (Lihat daftar Peminjaman)
	}
}

func handleGetPeminjamanByID(store *MemoryStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: Implementasi B3 (Lihat detail Peminjaman)
	}
}

func handleUpdatePeminjaman(store *MemoryStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: Implementasi B6 (Ubah Peminjaman)
	}
}

func handleDeletePeminjaman(store *MemoryStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: Implementasi B6 (Hapus Peminjaman)
	}
}