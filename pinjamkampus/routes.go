// routes.go
package main

import "net/http"

func NewHandler() http.Handler {
	mux := http.NewServeMux()
	store := NewMemoryStore()

	// Endpoint sumber-daya
	mux.HandleFunc("POST /sumber-daya", handleCreateSumberDaya(store))
	mux.HandleFunc("GET /sumber-daya", handleGetSumberDaya(store))
	mux.HandleFunc("GET /sumber-daya/{id}", handleGetSumberDayaByID(store))

	// Endpoint peminjaman
	mux.HandleFunc("POST /peminjaman", handleCreatePeminjaman(store))
	mux.HandleFunc("GET /peminjaman", handleGetPeminjaman(store))
	mux.HandleFunc("GET /peminjaman/{id}", handleGetPeminjamanByID(store))
	
	
	mux.HandleFunc("PUT /peminjaman/{id}", handleUpdatePeminjaman(store))
	mux.HandleFunc("DELETE /peminjaman/{id}", handleDeletePeminjaman(store))


	mux.HandleFunc("GET /kuota", handleGetKuota(store)) //endpoint untuk melihat sisa kuota peminjaman per hari

	return mux
}