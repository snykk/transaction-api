CREATE TABLE wallets (
  wallet_id uuid PRIMARY KEY,
  user_id uuid REFERENCES users(user_id),
  balance DECIMAL(15, 2) NOT NULL CHECK (balance >= 0), -- saldo tidak boleh negatif
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

