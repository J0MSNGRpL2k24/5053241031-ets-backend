# PROSES

## 1. Persona dan masalahnya
(2-3 kalimat dengan kata-katamu sendiri: siapa pengguna persona kamu, masalah apa yang diselesaikan fitur khasnya.)
 


Mahasiswa. masalah yang diselesaikan fitur khasnya adalah beberapa orang memborong alat sehingga yang lain tidak kebagian. Fitur khasnya adalah durasi maksimal 3 slot, kuota 2 peminjaman per orang per hari, dan endpoint untuk melihat sisa kuota.

## 2. Rancangan struktur data
(Tulis struct dan slice/map apa yang kamu simpan, dan kenapa dipilih begitu. Rancang ini sebelum menulis kode, lalu catat kalau berubah di tengah jalan.)

struct SumberDaya {
    ID   int    `json:"id"`
    Nama string `json:"nama"`
}

struct Peminjaman {
    ID           int    `json:"id"`
    SumberDayaID int    `json:"sumber_daya_id"`
    Peminjam     string `json:"peminjam"`
    Tanggal      string `json:"tanggal"`
    SlotMulai    int    `json:"slot_mulai"`
    SlotSelesai  int    `json:"slot_selesai"`
    Jumlah       int    `json:"jumlah"`
    Status       string `json:"status"`
} 

saya memilih struct dan slice/map ini karena saya ingin menyimpan data peminjaman dan sumber daya secara terstruktur. Struct SumberDaya menyimpan informasi tentang sumber daya yang tersedia, sedangkan struct Peminjaman menyimpan informasi tentang peminjaman yang dilakukan oleh mahasiswa. Slice digunakan untuk menyimpan daftar sumber daya dan peminjaman, sehingga memudahkan dalam pengelolaan data. Map digunakan untuk mengakses data dengan cepat berdasarkan ID, sehingga meningkatkan efisiensi dalam pencarian dan manipulasi data.


## 3. Satu keputusan desain
(Satu keputusan yang kamu ambil, pilihan lain yang kamu pertimbangkan, dan alasanmu memilih yang ini. Contoh: kenapa cek bentrok per slot, bukan per rentang.) 

saya memilih untuk menggunakan struktur clean architecture dengan memisahkan handler, store, dan model. Pilihan lain yang saya pertimbangkan adalah menggunakan monolithic architecture, namun saya memilih clean architecture karena untuk separation of concerns pada kebutuhan untuk peminjaman dan untuk sumber daya. Saya juga cek bentrok per slot karena peminjaman yang membuat kapasitas suatu slot terlampaui ditolak dengan 409, sehingga lebih mudah untuk mengelola dan memvalidasi peminjaman berdasarkan slot yang tersedia. 


## 4. Bug atau kendala (minimal dua)
### Bug 1
- **Gejala:** Respons JSON dari API mengembalikan *key* dengan awalan huruf kapital (contoh: `"Nama"`, `"Jenis"`, `"Stok"`) yang melanggar kontrak spesifikasi. Selain itu, input POST sering kali gagal divalidasi karena nilai *field* dianggap kosong oleh sistem.
- **Cara menemukan:** Saat menguji *endpoint* `POST /peminjaman` dan `POST /sumber-daya` di Postman, objek JSON yang dikembalikan memiliki format huruf besar, bukan format *snake_case* atau huruf kecil.
- **Penyebab dan solusi:** Penyebabnya adalah kesalahan penulisan (kelebihan spasi) pada *struct tag* JSON di file `models.go` (contoh: `` `json: "nama"` ``). Bahasa Go sangat sensitif dan otomatis mengabaikan tag tersebut jika ada spasi. Solusinya adalah menghapus spasi setelah titik dua menjadi `` `json:"nama"` `` sehingga *mapping* JSON berjalan sempurna.

### Bug 2
- **Gejala:** Server gagal di-*compile* atau dijalankan melalui terminal dengan munculnya pesan *error* fatal berbunyi `non-name req.Peminjam == "" || _ on left side of :=` dan `cannot use _ as value or type`.
- **Cara menemukan:** Terlihat langsung pada *output* terminal saat mengeksekusi perintah `go run .` setelah mencoba menggabungkan beberapa baris validasi input (termasuk validasi fitur khas Persona B dan fitur *Update*).
- **Penyebab dan solusi:** Penyebabnya adalah mencoba mendeklarasikan variabel baru menggunakan *short assignment* `:=` (seperti `_, err := time.Parse(...)`) di dalam satu baris deretan kondisi logika `||` pada blok `if`. Hal ini dilarang oleh sintaks Go. Solusinya adalah memisahkan validasi `time.Parse` menjadi blok `if` tersendiri agar mematuhi aturan sintaksis *compiler* Go.

## 5. Change request
Saya memutuskan untuk menggunakan map peminjamanPerSD untuk menghitung munculnya peminjaman dalam 1 kali loop 0(P) lalu menggabungkan dengan data sumber daya yang setelah itu di sort 0(S log S). Pendekatan aggregate memori ini sangat cepat karena tidak membebani komputasi dengan nested loop (0(P*S)), sehingga ringkasan per tanggal dapat ditampilkan dengan lebih cepat dan efisien untuk dashboard frontend.

