-- Migration for taxes table
-- Run this SQL script to create the taxes table

-- Create taxes table
CREATE TABLE IF NOT EXISTS taxes (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    nama_pajak varchar(100) NOT NULL,
    tipe_pajak varchar(10) NOT NULL CHECK (tipe_pajak IN ('sc', 'pb1')),
    presentase decimal(5,2) NOT NULL CHECK (presentase >= 0 AND presentase <= 100),
    deskripsi text,
    status varchar(20) DEFAULT 'active' CHECK (status IN ('active', 'inactive')),
    prioritas integer DEFAULT 0,
    created_at timestamptz DEFAULT NOW(),
    updated_at timestamptz DEFAULT NOW()
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_taxes_status ON taxes(status);
CREATE INDEX IF NOT EXISTS idx_taxes_prioritas ON taxes(prioritas);
CREATE INDEX IF NOT EXISTS idx_taxes_tipe_pajak ON taxes(tipe_pajak);

-- Add comments
COMMENT ON TABLE taxes IS 'Table untuk menyimpan data pajak (PB1, Service Charge, dll)';
COMMENT ON COLUMN taxes.tipe_pajak IS 'Tipe pajak: sc (Service Charge) atau pb1 (Pajak Barang dan Jasa 1)';
COMMENT ON COLUMN taxes.presentase IS 'Persentase pajak (0-100)';
COMMENT ON COLUMN taxes.status IS 'Status pajak: active atau inactive';
COMMENT ON COLUMN taxes.prioritas IS 'Urutan prioritas penerapan pajak (semakin tinggi semakin prioritas)';

-- Insert sample data (optional)
INSERT INTO taxes (nama_pajak, tipe_pajak, presentase, deskripsi, status, prioritas)
VALUES 
    ('PB1', 'pb1', 10.00, 'Pajak Barang dan Jasa 1', 'active', 1),
    ('Service Charge', 'sc', 5.00, 'Biaya layanan', 'active', 2)
ON CONFLICT DO NOTHING;

-- Verify
SELECT * FROM taxes ORDER BY prioritas DESC, nama_pajak ASC;
