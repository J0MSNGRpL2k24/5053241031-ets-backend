package main

import (
	"encoding/json"
	"net/http"
	"sort"
	"strconv"
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

		// PERSONA B: Durasi maksimal 3 slot
		if (req.SlotSelesai - req.SlotMulai + 1) > 3 {
			sendError(w, http.StatusBadRequest, "Durasi peminjaman maksimal 3 slot")
			return
		}

		store.mu.Lock()
		defer store.mu.Unlock()

		// PERSONA B: Kuota maksimal 2 peminjaman per hari
		jumlahPinjamHariIni := 0
		for _, p := range store.Peminjaman {
			if strings.EqualFold(p.Peminjam, req.Peminjam) && p.Tanggal == req.Tanggal && p.Status == "aktif" {
				jumlahPinjamHariIni++
			}
		}

		if jumlahPinjamHariIni >= 2 {
			sendError(w, http.StatusBadRequest, "Kuota peminjaman Anda (2x per hari) sudah habis")
			return
		}

		sd, exists := store.SumberDaya[req.SumberDayaID]
		if !exists {
			sendError(w, http.StatusBadRequest, "Sumber daya dengan ID tersebut tidak ada")
			return
		}

		if req.Jumlah < 1 || req.Jumlah > sd.Stok {
			sendError(w, http.StatusBadRequest, "Jumlah tidak boleh nol atau melebihi stok sumber daya")
			return
		}

		// B5: Deteksi Bentrok
		pemakaianPerSlot := make(map[int]int)
		for _, p := range store.Peminjaman {
			if p.Status == "aktif" && p.Tanggal == req.Tanggal && p.SumberDayaID == req.SumberDayaID {
				for s := p.SlotMulai; s <= p.SlotSelesai; s++ {
					pemakaianPerSlot[s] += p.Jumlah
				}
			}
		}

		isBentrok := false
		for s := req.SlotMulai; s <= req.SlotSelesai; s++ {
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

// B3: Lihat Daftar Peminjaman
func handleGetPeminjaman(store *MemoryStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		q := r.URL.Query()
		qTanggal, qStatus, qSDIDStr := q.Get("tanggal"), q.Get("status"), q.Get("sumber_daya_id")

		var qSDID int
		if qSDIDStr != "" {
			id, err := strconv.Atoi(qSDIDStr)
			if err != nil {
				sendError(w, http.StatusBadRequest, "Filter sumber_daya_id harus berupa angka")
				return
			}
			qSDID = id
		}

		store.mu.Lock()
		list := []Peminjaman{}
		for _, p := range store.Peminjaman {
			if (qTanggal == "" || p.Tanggal == qTanggal) &&
				(qStatus == "" || p.Status == qStatus) &&
				(qSDID == 0 || p.SumberDayaID == qSDID) {
				list = append(list, p)
			}
		}
		store.mu.Unlock()

		sort.Slice(list, func(i, j int) bool { return list[i].ID < list[j].ID })
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(list)
	}
}

// B3: Lihat Detail Satu Peminjaman
func handleGetPeminjamanByID(store *MemoryStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			sendError(w, http.StatusBadRequest, "ID harus berupa angka")
			return
		}

		store.mu.Lock()
		p, exists := store.Peminjaman[id]
		store.mu.Unlock()

		if !exists {
			sendError(w, http.StatusNotFound, "Peminjaman tidak ditemukan")
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(p)
	}
}

// B6: Ubah Peminjaman
// B6: Ubah Peminjaman
func handleUpdatePeminjaman(store *MemoryStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			sendError(w, http.StatusBadRequest, "ID harus berupa angka")
			return
		}

		var req Peminjaman
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			sendError(w, http.StatusBadRequest, "Body request tidak valid")
			return
		}

		// VALIDASI DIPISAH AGAR TIDAK ERROR SYNTAX
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

		lama, exists := store.Peminjaman[id]
		if !exists {
			sendError(w, http.StatusNotFound, "Peminjaman tidak ditemukan")
			return
		}

		sd, sdExists := store.SumberDaya[req.SumberDayaID]
		if !sdExists || req.Jumlah < 1 || req.Jumlah > sd.Stok {
			sendError(w, http.StatusBadRequest, "Sumber daya tidak valid atau jumlah melebihi stok")
			return
		}

		// B5: Deteksi Bentrok PUT (abaikan ID sendiri)
		pemakaianPerSlot := make(map[int]int)
		for _, p := range store.Peminjaman {
			if p.ID != id && p.Status == "aktif" && p.Tanggal == req.Tanggal && p.SumberDayaID == req.SumberDayaID {
				for s := p.SlotMulai; s <= p.SlotSelesai; s++ {
					pemakaianPerSlot[s] += p.Jumlah
				}
			}
		}

		isBentrok := false
		for s := req.SlotMulai; s <= req.SlotSelesai; s++ {
			if pemakaianPerSlot[s]+req.Jumlah > sd.Stok {
				isBentrok = true
				break
			}
		}

		if isBentrok {
			sendError(w, http.StatusConflict, "Bentrok: Kapasitas sumber daya tidak cukup pada slot tersebut")
			return
		}

		lama.SumberDayaID, lama.Peminjam, lama.Tanggal = req.SumberDayaID, req.Peminjam, req.Tanggal
		lama.SlotMulai, lama.SlotSelesai, lama.Jumlah = req.SlotMulai, req.SlotSelesai, req.Jumlah
		store.Peminjaman[id] = lama

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(lama)
	}
}

// B6: Hapus Peminjaman
func handleDeletePeminjaman(store *MemoryStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			sendError(w, http.StatusBadRequest, "ID harus berupa angka")
			return
		}

		store.mu.Lock()
		defer store.mu.Unlock()

		if _, exists := store.Peminjaman[id]; !exists {
			sendError(w, http.StatusNotFound, "Peminjaman tidak ditemukan")
			return
		}

		delete(store.Peminjaman, id)
		w.WriteHeader(http.StatusNoContent)
	}
}

// PERSONA B: Endpoint sisa kuota
func handleGetKuota(store *MemoryStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		peminjam := strings.TrimSpace(r.URL.Query().Get("peminjam"))
		tanggal := r.URL.Query().Get("tanggal")

		if peminjam == "" || tanggal == "" {
			sendError(w, http.StatusBadRequest, "Query peminjam dan tanggal wajib diisi")
			return
		}

		store.mu.Lock()
		terpakai := 0
		for _, p := range store.Peminjaman {
			if strings.EqualFold(p.Peminjam, peminjam) && p.Tanggal == tanggal && p.Status == "aktif" {
				terpakai++
			}
		}
		store.mu.Unlock()

		sisa := 2 - terpakai
		if sisa < 0 {
			sisa = 0
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]int{"sisa_kuota": sisa})
	}
}