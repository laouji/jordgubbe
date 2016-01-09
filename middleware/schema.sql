CREATE TABLE review_ios (
  id INT(11) NOT NULL,
  title VARCHAR(255) NULL,
  content TEXT,
  rating TINYINT(3) NOT NULL DEFAULT 0,
  author_name VARCHAR(255) NULL,
  author_uri VARCHAR(255) NULL,
  updated DATETIME NOT NULL,
  acquired DATETIME NOT NULL,
  PRIMARY KEY (id)
);
CREATE INDEX updated_idx on review_ios(updated);
CREATE INDEX rating_idx on review_ios(rating);
