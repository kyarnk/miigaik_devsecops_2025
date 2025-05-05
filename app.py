from flask import Flask, request, render_template
import sqlite3
from datetime import datetime
import pytz
import re
import os

app = Flask(__name__, template_folder='templates')

def init_db():
    conn = sqlite3.connect(':memory:')
    cursor = conn.cursor()
    cursor.execute("CREATE TABLE tasks (id INTEGER PRIMARY KEY, task TEXT, is_secret INTEGER DEFAULT 0)")
    tasks = [
        (1, 'Найти уязвимость', 0),
        (2, 'Пофиксить...а что фиксить...', 0),
        (3, 'Потраить чужие приложения', 0),
        (4, 'Написать приложение', 0),
        (5, 'Поднять куб?', 0),
        (6, 'TIME IS KEY', 0),
        (7, 'OR - это ошибка, а может и нет.', 0),
        (8, 'Время поможет', 0),
        (9, 'SecSence - лучшая команда!', 0),
        (10, 'V2hhdCBmbGFnIGNpcGhlciBzaG91bGQgSSB1c2U/', 0),
        (11, 'PO4TI_FL@G: {SYNT - J3YYQ0AR}', 1),
        (12, 'дедушка палит - даже не пытайся', 0),
        (13, '▶︎ •၊၊||၊|။||||| 0:10 Кишлак, Хаски - Громко', 0),
        (14, 'Савостин - легенда', 0),
        (15, 'admin? Куда мы лезем ▄︻デ══━一💥ඞඞඞඞඞඞඞඞඞඞ', 0),
    ]
    cursor.executemany("INSERT INTO tasks VALUES (?, ?, ?)", tasks)
    conn.commit()
    return conn

def is_vulnerable():
    moscow_time = datetime.now(pytz.timezone('Europe/Moscow'))
    hour = moscow_time.hour

    if 6 <= hour < 7:
        return True
    elif 7 <= hour < 9:
        return False
    elif 9 <= hour < 10:
        return True
    elif 10 <= hour < 12:
        return False
    elif 12 <= hour < 13:
        return True
    elif 13 <= hour < 15:
        return False
    elif 15 <= hour < 16:
        return True
    elif 16 <= hour < 18:
        return False
    elif 18 <= hour < 19:
        return True
    elif 19 <= hour < 21:
        return False
    elif 21 <= hour < 22:
        return True
    elif 22 <= hour < 24:
        return False
    else:
        return False

def build_query(user_filter, vulnerable):
    if not user_filter.strip():
        return "SELECT * FROM tasks WHERE id = -1", []

    if vulnerable:
        # Преднамеренная SQL-инъекция
        return f"SELECT * FROM tasks WHERE task LIKE '%{user_filter}%'", None

    if not re.match(r'^[a-zA-Zа-яА-Я0-9 ]+$', user_filter):
        return "SELECT * FROM tasks WHERE id = -1", []

    return "SELECT * FROM tasks WHERE task LIKE ? AND is_secret = 0", [f"%{user_filter}%"]

@app.route('/dedushka')
def index():
    search_query = request.args.get('search', '')
    conn = init_db()
    cursor = conn.cursor()

    vulnerable = is_vulnerable()
    try:
        query, params = build_query(search_query, vulnerable)
        if params is None:
            cursor.execute(query)  # инъекция
        else:
            cursor.execute(query, params)
        tasks = cursor.fetchall()
    except:
        tasks = []
    finally:
        conn.close()

    return render_template("index.html", tasks=tasks, search_query=search_query, vulnerable=vulnerable)



if __name__ == "__main__":
    app.run(host="0.0.0.0", port=5000)
