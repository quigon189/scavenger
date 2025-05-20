from flask_wtf import FlaskForm
from wtforms import StringField, FileField, SelectField
from wtforms.validators import DataRequired

class SendForm(FlaskForm):
    sender_name = StringField('ФИО', validators=[DataRequired()], render_kw={'class': 'form-control', 'required': True})
    sender_group = SelectField("Группа", validators=[DataRequired()], render_kw={'class': 'form-select', 'required': True})
    number = StringField("Номер и название работы", validators=[DataRequired()], render_kw={'class': 'form-control', 'required': True})

    send_file = FileField("Выберите файл", validators=[DataRequired()], render_kw={'class': 'form-control', 'required': True})

