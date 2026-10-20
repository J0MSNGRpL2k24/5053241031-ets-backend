package main

import (
	"encoding/json"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

// Helper agar membalas error selalu dalam format JSON (Aturan Umum)
func sendError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}


func handleCreateSumberDaya(store *MemoryStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var req SumberDaya
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			sendError(w, http.StatusBadRequest, "Body request tidak valid")
			return
		}

		
		req.Nama = strings.TrimSpace(req.Nama)
		if req.Nama == "" || len(req.Nama) > 50 {
			sendError(w, http.StatusBadRequest, "Nama wajib diisi dan maksimal 50 karakter")
			return
		}
		if req.Jenis != "ruang" && req.Jenis != "alat" {
			sendError(w, http.StatusBadRequest, "Jenis hanya boleh ruang atau alat")
			return
		}
		if req.Stok < 1 {
			sendError(w, http.StatusBadRequest, "Stok minimal 1")
			return
		}

		
		store.mu.Lock()
		req.ID = store.nextSumberDayaID
		store.SumberDaya[req.ID] = req
		store.nextSumberDayaID++
		store.mu.Unlock()

		
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(req)
	}
}

// Handler GET /sumber-daya
func handleGetSumberDaya(store *MemoryStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		list := []SumberDaya{}

		store.mu.Lock()
		
		for _, sd := range store.SumberDaya {
			list = append(list, sd)
		}
		store.mu.Unlock()

		
		sort.Slice(list, func(i, j int) bool {
			return list[i].ID < list[j].ID
		})

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(list)
	}
}


func handleGetSumberDayaByID(store *MemoryStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		
		idStr := r.PathValue("id")
		
		
		id, err := strconv.Atoi(idStr)
		if err != nil {
			sendError(w, http.StatusBadRequest, "ID harus berupa angka")
			return
		}

		store.mu.Lock()
		sd, exists := store.SumberDaya[id]
		store.mu.Unlock()

		
		if !exists {
			sendError(w, http.StatusNotFound, "Sumber daya tidak ditemukan")
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(sd)
	}
}