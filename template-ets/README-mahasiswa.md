# PinjamKampus API

- Nama dan NRP:
- Persona:
- Cara menjalankan: `go run .` lalu buka http://localhost:8080

## Fitur
- [ ] B1 Sumber daya
- [ ] B2 Buat peminjaman
- [ ] B3 Lihat peminjaman (daftar, filter, detail)
- [ ] B4 Validasi dan error JSON
- [ ] B5 Deteksi bentrok (409)
- [ ] B6 Ubah dan hapus
- [ ] B7 Aman untuk permintaan bersamaan
- [ ] Fitur khas persona
- [ ] Change request

## Yang belum selesai


B2	Buat peminjaman	Satu peminjaman = satu sumber daya, satu tanggal, rentang slot, dan jumlah unit
B3	Lihat peminjaman	Daftar (dengan filter), dan detail satu peminjaman
B4	Validasi dan error	Semua input diperiksa; error selalu berbentuk JSON dengan kode status yang tepat
B5	Deteksi bentrok	Peminjaman yang membuat kapasitas suatu slot terlampaui ditolak dengan 409
B6	Ubah dan hapus	PUT mengganti data peminjaman, DELETE menghapusnya; keduanya tetap menjaga aturan bentrok
B7	Aman untuk permintaan bersamaan	Banyak permintaan yang datang bersamaan tidak boleh merusak data atau membuat dua peminjaman menempati slot yang sama
