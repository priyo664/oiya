-- Buat database (jika belum ada)
CREATE DATABASE IF NOT EXISTS oiya_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE oiya_db;

-- Tabel users (penumpang & admin)
CREATE TABLE users (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    phone VARCHAR(20) NOT NULL UNIQUE,
    email VARCHAR(100),
    password TEXT,
    role ENUM('passenger', 'admin') NOT NULL DEFAULT 'passenger',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Tabel drivers
CREATE TABLE drivers (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    phone VARCHAR(20) NOT NULL UNIQUE,
    email VARCHAR(100),
    password TEXT,
    vehicle_type VARCHAR(50),
    vehicle_number VARCHAR(30),
    photo TEXT,
    status ENUM('available', 'busy', 'offline') DEFAULT 'offline',
    balance DECIMAL(12,2) DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Tabel trips
CREATE TABLE trips (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    passenger_id BIGINT,
    driver_id BIGINT,
    origin TEXT,
    destination TEXT,
    status ENUM('pending', 'accepted', 'in_progress', 'completed', 'cancelled') DEFAULT 'pending',
    price DECIMAL(12,2),
    rating INT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (passenger_id) REFERENCES users(id) ON DELETE SET NULL,
    FOREIGN KEY (driver_id) REFERENCES drivers(id) ON DELETE SET NULL
);

-- Tabel payments
CREATE TABLE payments (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    trip_id BIGINT,
    amount DECIMAL(12,2),
    method ENUM('cash', 'qris'),
    status ENUM('pending', 'completed', 'failed') DEFAULT 'pending',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (trip_id) REFERENCES trips(id) ON DELETE CASCADE
);

-- Tabel chats
CREATE TABLE chats (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    sender_role ENUM('passenger', 'driver', 'admin'),
    sender_id BIGINT,
    receiver_role ENUM('passenger', 'driver', 'admin'),
    receiver_id BIGINT,
    message TEXT,
    sent_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Tabel banners (untuk admin pasang iklan)
CREATE TABLE banners (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    image_url TEXT,
    link_url TEXT,
    views INT DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
