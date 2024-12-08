-- +migrate Up

ALTER TABLE materials
ADD CONSTRAINT unique_name UNIQUE (name);

INSERT INTO materials (name, category, image_url, description, purchase_location)
VALUES
    ('Обои Бежевый однотон', 'Обои', '/uploads/oboi1.jpg', 'Бежевый однотон', 'ОБОИ, Каслинская, 5'),
    ('Обои Синие листья', 'Обои', '/uploads/oboi2.jpg', 'Синие листья', 'ОБОИ, Каслинская, 5'),
    ('Обои Серая штукатурка', 'Обои', '/uploads/oboi3.jpg', 'Серая штукатурка', 'ОБОИ, Каслинская, 5'),
    ('Обои Бежевый животные саванна', 'Обои', '/uploads/oboi4.jpg', 'Бежевый животные саванна', 'ОБОИ, Каслинская, 5'),
    ('Линолеум Бежевый доска', 'Линолеум', '/uploads/linonium.jpg', 'Бежевый доска', 'ЛинДвор, Каслинская, 5'),
    ('Линолеум Дуб доска', 'Линолеум', '/uploads/linolium2.jpg', 'Дуб доска', 'ЛинДвор, Каслинская, 5'),
    ('Линолеум Темный Дуб доска', 'Линолеум', '/uploads/linolium3.jpg', 'Темный Дуб доска', 'ЛинДвор, Каслинская, 5')
ON CONFLICT (name) DO NOTHING;
