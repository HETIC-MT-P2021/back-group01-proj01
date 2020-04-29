/*
    This file is used by the docker-compose build command to build the mysql db
    
    Tables:
    * category : stores categories (id, name, desc, creation, update)
    * image : stores images (id, name, desc, type, creation, update, category ID)
    * tag : stores tags (id, name, creation date)
    * image_tag : links images to tags by ids (Many to Many relation)
*/

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
    slug VARCHAR(255) UNIQUE,
    description TEXT,
    type VARCHAR(10),
    created_at DATETIME,
    updated_at DATETIME,
    category_id INT, 
    FOREIGN KEY (category_id) 
        REFERENCES category(id)   
        ON DELETE CASCADE 
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
    FOREIGN KEY (image_id) 
        REFERENCES image(id),
    FOREIGN KEY (tag_id) 
        REFERENCES tag(id)
);


/*
    Starter sample data
    
    3 categories with ids 1 to 3
    7 images metadata with ids 1 to 7
    11 tags with ids 1 to 11, linked to images
*/

INSERT INTO `category` (`id`, `name`, `description`, `created_at`, `updated_at`) VALUES
('1', 'holidays', 'a collection of holidays pictures', '2020-04-28 19:25:05', '2020-04-28 19:25:05'),
('2', 'animals', 'a collection of cute animals', '2020-04-28 19:29:18', '2020-04-28 19:29:18'),
('3', 'cars', 'a collection of cars', '2020-04-28 19:29:25', '2020-04-28 19:29:25');

INSERT INTO `image` (`id`, `name`, `slug`, `description`, `type`, `created_at`, `updated_at`, `category_id`) VALUES
('1', 'trip to Cancun', 'goo1u5d89r', 'what a wonderful trip ! ', '', '2020-04-28 19:25:49', '2020-04-28 19:25:49', '1'),
('2', 'trip to Tahiti', '7w0anxy08y', 'what a wonderful trip ! ', '', '2020-04-28 19:26:02', '2020-04-28 19:26:02', '1'),
('3', 'car', 'v028zdr051', 'car is fast ', '', '2020-04-28 19:29:47', '2020-04-28 19:30:33', '3'),
('4', 'grey car', '9vyx4y1t7k', 'car is fast ', '', '2020-04-28 19:30:02', '2020-04-28 19:30:46', '3'),
('5', 'car', '8paa447pfk', 'car is fast ', '', '2020-04-28 19:30:17', '2020-04-28 19:30:55', '3'),
('6', 'cat', '2oqn4u7hhx', 'A cute cat', '', '2020-04-28 19:33:52', '2020-04-28 19:33:52', '2'),
('7', 'dog', '9hjtv67dpk', 'A cute dog', '', '2020-04-28 19:35:06', '2020-04-28 19:35:06', '2');

INSERT INTO `tag` (`id`, `name`, `created_at`, `updated_at`) VALUES
('1', 'pool', '2020-04-28 19:25:49', '2020-04-28 19:25:49'),
('2', 'holidays', '2020-04-28 19:25:49', '2020-04-28 19:25:49'),
('3', 'sun', '2020-04-28 19:25:49', '2020-04-28 19:25:49'),
('4', 'car', '2020-04-28 19:29:47', '2020-04-28 19:29:47'),
('5', 'fast', '2020-04-28 19:29:47', '2020-04-28 19:29:47'),
('6', 'red', '2020-04-28 19:29:47', '2020-04-28 19:29:47'),
('7', 'grey', '2020-04-28 19:30:02', '2020-04-28 19:30:02'),
('8', 'cat', '2020-04-28 19:33:52', '2020-04-28 19:33:52'),
('9', 'cute', '2020-04-28 19:33:52', '2020-04-28 19:33:52'),
('10', 'love', '2020-04-28 19:33:52', '2020-04-28 19:33:52'),
('11', 'dog', '2020-04-28 19:35:06', '2020-04-28 19:35:06');

INSERT INTO `image_tag` (`image_id`, `tag_id`) VALUES
('1', '1'),
('1', '2'),
('1', '3'),
('2', '1'),
('2', '2'),
('2', '3'),
('3', '4'),
('3', '5'),
('3', '6'),
('4', '4'),
('4', '5'),
('4', '7'),
('5', '4'),
('5', '5'),
('5', '7'),
('6', '8'),
('6', '9'),
('6', '10'),
('7', '11'),
('7', '9'),
('7', '10');

