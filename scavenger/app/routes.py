from flask import render_template, flash, redirect, url_for
from werkzeug.utils import secure_filename

import os

base_dir = os.path.abspath(os.path.dirname(__file__))

base_dir = os.path.join(base_dir, "dump")

from app import app
from app.forms import SendForm

@app.route('/', methods=["GET"])
@app.route('/index', methods=["GET", "POST"])
def index():

    form = SendForm()

    if form.validate_on_submit():

        group_dir = form.sender_group.data.replace(" ","_").lower()
        stud_dir = form.sender_name.data.replace(" ","_").lower()
        job_dir = form.number.data.replace(" ","_").lower()

        full_path_dir = os.path.join(base_dir, group_dir, stud_dir, job_dir)

        if not os.path.exists(full_path_dir):
            os.makedirs(full_path_dir)

        if not form.send_file.data:
            flash('Файл отсутствует')
            return redirect(url_for('index'))

        n = len(os.listdir(full_path_dir))

        ext = secure_filename(form.send_file.data.filename).split(".")[1]

        form.send_file.data.save(os.path.join(full_path_dir, f"{n+1}.{ext}"))

        flash("Ваша работа отправлена")
    
    elif form.is_submitted():

        flash("Необходимо заполнить все поля")

    return render_template("index.html", title="Сборщик практических работ", form=form)