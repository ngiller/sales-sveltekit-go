-- Migration: Add disc_value column to quotation and quotation_master tables
-- Kolom ini menyimpan nilai nominal discount (dalam Rupiah)

ALTER TABLE `quotation`
  ADD COLUMN `disc_value` float NOT NULL DEFAULT 0 AFTER `disc`;

ALTER TABLE `quotation_master`
  ADD COLUMN `disc_value` float NOT NULL DEFAULT 0 AFTER `disc`;

-- Update existing rows: calculate disc_value from existing disc percentage
UPDATE `quotation` SET `disc_value` = `total` * `disc` / 100 WHERE `disc` IS NOT NULL AND `disc` > 0 AND `total` IS NOT NULL AND `total` > 0;
UPDATE `quotation_master` SET `disc_value` = `total` * `disc` / 100 WHERE `disc` IS NOT NULL AND `disc` > 0 AND `total` IS NOT NULL AND `total` > 0;
