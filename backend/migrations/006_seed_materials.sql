-- +migrate Up

-- Удаление уникального ограничения на name, если оно существует
ALTER TABLE materials
    DROP CONSTRAINT IF EXISTS unique_name;

-- Добавление уникального ограничения на article
ALTER TABLE materials
    ADD CONSTRAINT unique_article UNIQUE (article);

INSERT INTO materials (name, article, category, image_url, description, purchase_location)
VALUES
    ('Обои Бежевый однотон', '60350-03', 'Обои', '/uploads/oboi1.jpg', 'Бежевый однотон', 'ОБОИ, Каслинская, 5'),
    ('Обои Синие листья', '60349-05', 'Обои', '/uploads/oboi2.jpg', 'Синие листья', 'ОБОИ, Каслинская, 5'),
    ('Обои Серая штукатурка', '60662-07', 'Обои', '/uploads/oboi3.jpg', 'Серая штукатурка', 'ОБОИ, Каслинская, 5'),
    ('Обои Бежевый животные саванна', '555028', 'Обои', '/uploads/oboi4.jpg', 'Бежевый животные саванна', 'ОБОИ, Каслинская, 5'),
    ('Линолеум Бежевый доска', 'ARNOLD 1', 'Линолеум', '/uploads/linonium.jpg', 'Бежевый доска', 'ЛинДвор, Каслинская, 5'),
    ('Линолеум Дуб доска', 'ДУБ СТЕЛЛА', 'Линолеум', '/uploads/linolium2.jpg', 'Дуб доска', 'ЛинДвор, Каслинская, 5'),
    ('Линолеум Темный Дуб доска', 'ДУБ КОВАЛЬСК', 'Линолеум', '/uploads/linolium3.jpg', 'Темный Дуб доска', 'ЛинДвор, Каслинская, 5'),
    ('Обои серый фон мелкие цветы', '588257', 'Обои', '/uploads/oboi5.jpg', 'серый фон мелкие цветы', 'ОБОИ, Каслинская, 5'),
    ('Обои белый фон мелкие цветы', '588251', 'Обои', '/uploads/oboi6.jpg', 'белый фон мелкие цветы', 'ОБОИ, Каслинская, 5'),
    ('Обои Крупные зеленые листья', '588173', 'Обои', '/uploads/oboi7.jpg', 'Крупные зеленые листья', 'ОБОИ, Каслинская, 5'),
    ('Обои зеленый однотон с текстурой', '588183', 'Обои', '/uploads/oboi8.jpg', 'зеленый однотон с текстурой', 'ОБОИ, Каслинская, 5'),
    ('Обои серый однотон с текстурой', '588189', 'Обои', '/uploads/oboi9.jpg', 'серый однотон с текстурой', 'ОБОИ, Каслинская, 5'),
    ('Обои белый однотон с текстурой', '588182', 'Обои', '/uploads/oboi10.jpg', 'белый однотон с текстурой', 'ОБОИ, Каслинская, 5'),
    ('Обои серый джинсовый однотон', '287397', 'Обои', '/uploads/oboi11.jpg', 'серый джинсовый однотон', 'ОБОИ, Каслинская, 5'),
    ('Обои светло-зеленый фисташковый однотон', '287393', 'Обои', '/uploads/oboi12.jpg', 'светло-зеленый фисташковый однотон', 'ОБОИ, Каслинская, 5'),
    ('Ламинат Дуб Тамми белый ламинат под доску', '60628', 'Ламинат', '/uploads/laminat1.jpg', 'Дуб Тамми белый ламинат под доску', 'ЛинДвор, Каслинская, 5'),
    ('Ламинат Дуб Мадейра холодный коричневый под доску', '60621', 'Ламинат', '/uploads/laminat2.jpg', 'Дуб Мадейра холодный коричневый под доску', 'ЛинДвор, Каслинская, 5'),
    ('Ламинат Дуб Монтуар теплый коричневый под доску', '35007', 'Ламинат', '/uploads/laminat3.jpg', 'Дуб Монтуар теплый коричневый под доску', 'ЛинДвор, Каслинская, 5')
ON CONFLICT (article) DO NOTHING;
