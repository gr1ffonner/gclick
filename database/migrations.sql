CREATE TABLE material (
    uuid UUID PRIMARY KEY,
    material_type VARCHAR(50) NOT NULL CHECK (material_type IN ('статья', 'видеоролик', 'презентация')),
    publication_status VARCHAR(50) NOT NULL CHECK (publication_status IN ('архивный', 'активный')),
    title VARCHAR(255) NOT NULL,
    content TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = now();
   RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_material_updated_at
BEFORE UPDATE ON material
FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_column();
