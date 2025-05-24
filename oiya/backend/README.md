OIYA (OJEK INDONESIA RAYA)
TEKNOLOGI DAN ARSITEKTUR
1.	Backend		: Golang dengan echo untuk API yang cepat dan scalable.
2.	Frontend		: Bisa pakai React.
3.	Database		: MySQL untuk data pengguna, transaksi, dan trip history.
4.	Lokasi & Maps	: Leaflet.js dan OpenStreetMap untuk penentuan lokasi dan navigasi.
5.	Real-time		: WebSocket atau MQTT untuk update posisi driver secara real-time.


CHECK LIST FRONTEND FILE
a.	Penumpang
- Login:	Form nomor + OTP
- Daftar	Form data pengguna baru
- Halaman Utama	Tombol pesan layanan + peta
- Form Pesan Driver	Input lokasi + tombol konfirmasi
- Tracking Pesanan	Peta + info driver
- Pembayaran	Konfirmasi pembayaran tunai / QRIS
- Riwayat & Rating	Daftar trip + tombol beri rating
- Chat	Penumpang <--> driver

b.	Driver
- Login Driver	Form nomor HP + OTP untuk driver login
- Register Driver	Form pendaftaran data lengkap driver
- Dashboard Driver	Halaman utama driver, status & navigasi order
- Terima Order	Halaman/komponen untuk terima atau tolak order
- Tracking Pengantaran	Tracking perjalanan antar penumpang
- Riwayat Trip Driver	Daftar perjalanan + rating dari penumpang
- Top Up Saldo	Form isi saldo membership driver
- Chat	driver <--> penumpang

c.	Admin 
- Halaman Admin	Keterangan
- Login Admin	Form login untuk admin
- Dashboard Admin	Ringkasan statistik, grafik, dan quick action
- Manajemen User	Daftar dan kelola pelanggan & driver
- Manajemen Order	Pantau, edit, dan kelola order & status
- Manajemen Pembayaran	Monitoring transaksi dan status pembayaran
- Laporan & Statistik	Grafik dan laporan perjalanan, pendapatan
- Pengaturan Sistem	Setting umum sistem, tarif, dan konfigurasi
- Chat	driver <--> admin
- Broadcast message	Admin Penumpangdriver
- Tool untuk pasang iklan banner di homepage penumpang	•	Add banner di homepage penumpang 
•	Link banner ke jajano.id
•	Statistic view 

Backend dan database silahkan dibuat menyesuaikan dengan front end. Jangan ada yg di hardcode. Ikuti pola coding yg bagus
