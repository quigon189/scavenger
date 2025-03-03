from flask_wtf import FlaskForm
from wtforms import StringField, FileField
from wtforms.validators import DataRequired

class SendForm(FlaskForm):
    sender_name = StringField('ФИО', validators=[DataRequired()])
    sender_group = StringField("Группа", validators=[DataRequired()])

    number = StringField("Номер и название работы", validators=[DataRequired()])

    send_file = FileField("Выберите файл", validators=[DataRequired()])

