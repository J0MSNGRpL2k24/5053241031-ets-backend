# AI-LOG

Catat setiap penggunaan AI yang menghasilkan kode. Jika tidak memakai AI sama sekali, tulis satu baris: **tidak memakai AI**.

Log yang jujur tidak mengurangi nilai, sebanyak apa pun AI dipakai. Yang mengurangi nilai adalah log yang tidak sesuai dengan kode, atau isi log yang tidak bisa dijelaskan.

## Buat Peminjaman, Validasi, dan Race Condition (B2, B4, B7)
- Prompt: Meminta pengecekan apakah fungsi B2 (Buat Peminjaman) yang saya buat sudah aman dan sesuai kontrak ETS.
- Dipakai: Kode `handleCreatePeminjaman` yang berisi `store.mu.Lock()` dan struktur perbaikan *imports* bawaan Go di `handler_peminjaman.go`.
- Diubah: Menyesuaikan tag *struct* JSON di `models.go` (menghapus spasi setelah titik dua) karena hasil awalnya membuat *response body* di Postman menggunakan huruf kapital yang melanggar spesifikasi API.

## Lihat Peminjaman dan Deteksi Bentrok (B3, B5)
- Prompt: Meminta implementasi B3 (GET daftar dengan filter opsional dan GET by ID) serta logika perhitungan deteksi bentrok (B5).
- Dipakai: Fungsi `handleGetPeminjaman`, `handleGetPeminjamanByID`, dan blok algoritma perhitungan akumulasi `pemakaianPerSlot` di `handler_peminjaman.go`.
- Diubah: Tidak ada perubahan secara logika, algoritma validasi kapasitas stok per slot digunakan utuh karena sudah menangani status "aktif" dengan tepat.

## Ubah, Hapus Peminjaman, dan Fitur Persona B (B6, Persona B)
- Prompt: Meminta implementasi B6 (PUT dan DELETE) serta modifikasi khusus untuk Persona B (batas durasi 3 slot, limit kuota 2x per hari, dan *endpoint* sisa kuota).
- Dipakai: Fungsi `handleUpdatePeminjaman`, `handleDeletePeminjaman`, dan fungsi baru `handleGetKuota` di file `handler_peminjaman.go`.
- Diubah: Memperbaiki letak fungsi `handleGetKuota` yang sempat tidak sengaja ter-*paste* di dalam fungsi DELETE, serta merombak baris `time.Parse` di dalam kondisi `if` pada fungsi *Update* karena menyebabkan *syntax error* (`non-name on left side of :=`).

## Change Request A (Dashboard Ringkasan Harian)
- Prompt: Meminta pembuatan *endpoint* `/ringkasan` untuk mengagregasi data peminjaman harian (menghitung jumlah dan total unit slot) berdasarkan lembar CR dari dosen tanpa membuat file *handler* baru.
- Dipakai: *Struct* `RingkasanSumberDaya` dan `RingkasanResponse` di `models.go`, serta fungsi `handleGetRingkasan` di `handler_peminjaman.go`.
- Diubah: Digunakan sepenuhnya sesuai arahan AI karena sudah menggunakan pendekatan *in-memory map* dan *sorting* yang efisien tanpa harus merombak sistem penyimpanan utama (`store.go`).