import sys
import os
import openai
import base64
import requests

# OpenAI API Key из переменной окружения
api_key = os.getenv("OPENAI_API_KEY")

# Проверяем, установлен ли api_key
if not api_key:
    print("Ошибка: переменная окружения OPENAI_API_KEY не установлена")
    sys.exit(1)

# Получаем image_path и promt_text из аргументов командной строки
if len(sys.argv) != 3:
    print("Использование: python ai.py <image_path> <promt_text>")
    sys.exit(1)

image_path = sys.argv[1]
promt_text = sys.argv[2]

# Функция для кодирования изображения
def encode_image(image_path):
    with open(image_path, "rb") as image_file:
        return base64.b64encode(image_file.read()).decode('utf-8')

# Получаем base64 строку
base64_image = encode_image(image_path)

headers = {
    "Content-Type": "application/json",
    "Authorization": f"Bearer {api_key}"
}

payload = {
    "model": "gpt-4o",
    "messages": [
        {
            "role": "user",
            "content": [
                {
                    "type": "text",
                    "text": promt_text
                },
                {
                    "type": "image_url",
                    "image_url": {
                        "url": f"data:image/jpeg;base64,{base64_image}"
                    }
                }
            ]
        }
    ],
    "max_tokens": 300
}

response = requests.post("https://api.openai.com/v1/chat/completions", headers=headers, json=payload)

answer = response.json()

# Проверяем статус код
if response.status_code != 200:
    # API вернул ошибку
    print("Ошибка от OpenAI API:", answer.get('error', 'Неизвестная ошибка'))
    sys.exit(1)

# Продолжаем, если 'choices' существует
if 'choices' in answer:
    print(answer['choices'][0]['message']['content'])
else:
    print("Неожиданная структура ответа API:", answer)
    sys.exit(1)