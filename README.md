### Сборщик практическх работ

Веб интерфес из одной страницы для отправки выполненых практических работ преподавателю

Используемые фреймворки:
1. Flask
2. Jinja2
3. Bootstrap 5

Установка:

1. Клонируйте данный репозиторий:
```bash
git clone https://github.com/quigon189/scavenger.git
```

2. Установите `python` и `uv` любым удобным спсобом:
```bash
pip install uv #один из варинатов установки uv
```

3. Выполните `uv sync` для установки всех необходимых зависимостей:
```bash
cd scavenger
uv sync
```

4. Установите в переменную окружения DUMP абсолютный путь до каталога:
```bash
export DUMP=/home/user/DUMP #Linux

set DUMP="C:\\DUMP" #Windows
```

Для запуска:
```bash
uv run flask --debug run #запуск flask для отладки

uv run gunicorn --reload --log-level debug -b 0.0.0.0:1234 wsgi:app #запуск через gunicorn
```
