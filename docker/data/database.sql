CREATE TABLE IF NOT EXISTS category (
    id INT PRIMARY KEY NOT NULL AUTO_INCREMENT,
    name VARCHAR(255),
    description TEXT,
    created_at DATETIME,
    updated_at DATETIME
);

CREATE TABLE IF NOT EXISTS image (
    id INT PRIMARY KEY NOT NULL AUTO_INCREMENT,
    name VARCHAR(255),
    slug VARCHAR(255),
    description TEXT,
    created_at DATETIME,
    updated_at DATETIME,
    id_category INT, 
    FOREIGN KEY (id_category) REFERENCES category(id)                      
);

CREATE TABLE IF NOT EXISTS tag (
    id INT PRIMARY KEY NOT NULL AUTO_INCREMENT,
    name VARCHAR(255),
    created_at DATETIME,
    updated_at DATETIME
);

CREATE TABLE IF NOT EXISTS image_tag (
    image_id INT,
    tag_id INT,
    FOREIGN KEY (image_id) REFERENCES image(id),
    FOREIGN KEY (tag_id) REFERENCES tag(id)
);

