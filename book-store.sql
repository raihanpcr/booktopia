-- ==============================
-- TABLE: users
-- ==============================
CREATE TABLE users (
    user_id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    role ENUM('admin', 'pembeli') DEFAULT 'pembeli',
    saldo DECIMAL(12, 2) DEFAULT 0.00,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ==============================
-- TABLE: books
-- ==============================
CREATE TABLE books (
    book_id INT PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(200) NOT NULL,
    author VARCHAR(100),
    publisher VARCHAR(100),
    year_published YEAR,
    category VARCHAR(100),
    price DECIMAL(10, 2) NOT NULL,
    status ENUM('available', 'unavailable') DEFAULT 'available',
    is_donation_only BOOLEAN DEFAULT FALSE,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ==============================
-- TABLE: transactions
-- ==============================
CREATE TABLE transactions (
    transaction_id INT PRIMARY KEY AUTO_INCREMENT,
    user_id INT,
    transaction_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    total_amount DECIMAL(12, 2) NOT NULL,
    status ENUM('pending', 'completed', 'cancelled') DEFAULT 'pending',
    FOREIGN KEY (user_id) REFERENCES users(user_id)
);

-- ==============================
-- TABLE: transaction_details
-- ==============================
CREATE TABLE transaction_details (
    detail_id INT PRIMARY KEY AUTO_INCREMENT,
    transaction_id INT,
    book_id INT,
    quantity INT NOT NULL,
    price_per_unit DECIMAL(10, 2) NOT NULL,
    FOREIGN KEY (transaction_id) REFERENCES transactions(transaction_id),
    FOREIGN KEY (book_id) REFERENCES books(book_id)
);

-- ==============================
-- TABLE: top_up
-- ==============================
CREATE TABLE top_up (
    top_up_id INT PRIMARY KEY AUTO_INCREMENT,
    user_id INT,
    amount DECIMAL(12, 2) NOT NULL,
    method VARCHAR(50),
    top_up_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    status ENUM('pending', 'success', 'failed') DEFAULT 'pending',
    FOREIGN KEY (user_id) REFERENCES users(user_id)
);

-- ==============================
-- TABLE: ebook_gift_log
-- ==============================
CREATE TABLE ebook_gift_log (
    gift_id INT PRIMARY KEY AUTO_INCREMENT,
    donor_id INT NOT NULL,
    recipient_email VARCHAR(100) NOT NULL,
    book_id INT NOT NULL,
    gift_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    message TEXT,
    status ENUM('pending', 'accepted', 'declined') DEFAULT 'pending',
    recipient_user_id INT DEFAULT NULL,
    FOREIGN KEY (donor_id) REFERENCES users(user_id),
    FOREIGN KEY (book_id) REFERENCES books(book_id),
    FOREIGN KEY (recipient_user_id) REFERENCES users(user_id)
);
