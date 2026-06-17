INSERT INTO roles (name, description) VALUES
('SUPER_ADMIN', 'Pemilik sistem dengan akses penuh'),
('ADMIN', 'Admin yang mengelola layanan dan order'),
('CUSTOMER', 'Pelanggan yang melakukan pemesanan jasa');

INSERT INTO service_categories (name, slug, description) VALUES
('Jasa Website', 'jasa-website', 'Layanan pembuatan website company profile, landing page, dan aplikasi web'),
('Joki Tugas', 'joki-tugas', 'Layanan bantuan pengerjaan tugas sesuai kebutuhan customer'),
('Perbaikan Bug', 'perbaikan-bug', 'Layanan debugging dan perbaikan error aplikasi'),
('Konsultasi Project', 'konsultasi-project', 'Layanan konsultasi struktur project, database, dan deployment');