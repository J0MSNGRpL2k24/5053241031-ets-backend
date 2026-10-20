package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

func handleCreatePeminjaman(store *MemoryStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var req Peminjaman
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			sendError(w, http.StatusBadRequest, "Body request tidak valid")
			return
		}

		
		req.Peminjam = strings.TrimSpace(req.Peminjam)
		if req.Peminjam == "" || len(req.Peminjam) > 50 {
			sendError(w, http.StatusBadRequest, "Peminjam wajib diisi maksimal 50 karakter")
			return
		}

		
		if _, err := time.Parse("2006-01-02", req.Tanggal); err != nil {
			sendError(w, http.StatusBadRequest, "Format tanggal harus YYYY-MM-DD yang valid")
			return
		}

		
		if req.SlotMulai < 1 || req.SlotSelesai > 8 || req.SlotMulai > req.SlotSelesai {
			sendError(w, http.StatusBadRequest, "Slot harus antara 1-8 dan slot mulai <= slot selesai")
			return
		}

		
		store.mu.Lock()
		defer store.mu.Unlock() 

		sd, exists := store.SumberDaya[req.SumberDayaID]
		if !exists {
			sendError(w, http.StatusBadRequest, "Sumber daya dengan ID tersebut tidak ada")
			return
		}

		if req.Jumlah < 1 || req.Jumlah > sd.Stok {
			sendError(w, http.StatusBadRequest, "Jumlah tidak boleh nol atau melebihi stok sumber daya")
			return
		}

		// --- IMPLEMENTASI B5: DETEKSI BENTROK ---
		// Hitung akumulasi pemakaian unit untuk setiap slot pada tanggal dan sumber daya yang sama
		pemakaianPerSlot := make(map[int]int)

		for _, p := range store.Peminjaman {
			// Hanya peminjaman "aktif" yang menahan slot
			if p.Status == "aktif" && p.Tanggal == req.Tanggal && p.SumberDayaID == req.SumberDayaID {
				// Tambahkan jumlah pemakaian ke slot yang direntang peminjaman ini
				for s := p.SlotMulai; s <= p.SlotSelesai; s++ {
					pemakaianPerSlot[s] += p.Jumlah
				}
			}
		}

		// Periksa rentang slot peminjaman baru ini
		isBentrok := false
		for s := req.SlotMulai; s <= req.SlotSelesai; s++ {
			// Jika pemakaian saat ini ditambah jumlah baru melebihi stok, tolak!
			if pemakaianPerSlot[s]+req.Jumlah > sd.Stok {
				isBentrok = true
				break
			}
		}

		if isBentrok {
			sendError(w, http.StatusConflict, "Bentrok: Kapasitas sumber daya tidak cukup pada slot tersebut")
			return
		}
		


		req.Status = "aktif" 
		req.ID = store.nextPeminjamanID
		store.Peminjaman[req.ID] = req
		store.nextPeminjamanID++

		
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(req)
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