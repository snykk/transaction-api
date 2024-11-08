CREATE TABLE transactions (
    transaction_id uuid PRIMARY KEY,
    wallet_id uuid NOT NULL REFERENCES wallets(wallet_id) ON DELETE CASCADE, -- wallet terkait
    product_id INT REFERENCES products(product_id) ON DELETE SET NULL, -- produk yang dibeli, nullable untuk transaksi non-pembelian
    amount DECIMAL(15, 2) NOT NULL CHECK (amount > 0), -- jumlah uang yang ditransaksikan
    quantity INT CHECK (
        (product_id IS NULL AND quantity IS NULL) OR
        (product_id IS NOT NULL AND quantity > 0)
    ),
    transaction_type VARCHAR(20) NOT NULL CHECK (transaction_type IN ('deposit', 'withdraw', 'purchase')), -- tipe transaksi
    -- status VARCHAR(20) NOT NULL CHECK (status IN ('pending', 'completed', 'failed')), -- status transaksi
    -- description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

