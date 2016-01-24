CREATE TABLE review (
  id INT(11) NOT NULL,
  title VARCHAR(255) NULL,
  content TEXT,
  rating TINYINT(3) NOT NULL DEFAULT 0,
  device_type TINYINT(3) NOT NULL,
  device_name VARCHAR(255) NULL,
  author_name VARCHAR(255) NULL,
  author_uri VARCHAR(255) NULL,
  created DATETIME NULL,
  updated DATETIME NOT NULL,
  acquired DATETIME NOT NULL,
  PRIMARY KEY (id, device_type)
);
CREATE INDEX updated_idx on review(updated);
CREATE INDEX rating_idx on review(rating);
