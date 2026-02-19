-- Add product_image_url column to projects table
ALTER TABLE projects ADD COLUMN IF NOT EXISTS product_image_url TEXT;
