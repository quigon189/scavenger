Сборщик практическх работ

Для установки:

    python -m venv venv
    pip install -r reqirements.txt

Для запуска

    gunicorn -b 0.0.0.0:1234 wsgi:app