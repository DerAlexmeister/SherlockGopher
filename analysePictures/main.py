from sherlockanalysepictures.getmetadata import *
from flask import Flask
from flask_script import Manager
from flask_script import Server
from celery import Celery

app = Flask(__name__)
app.config.update(
    CELERY_BROKER_URL='redis://localhost:6379',
    CELERY_RESULT_BACKEND='redis://localhost:6379'
)

celery = make_celery(app)
manager = Manager(app)

def make_celery(app):
    celery = Celery(
        app.import_name,
        backend=app.config['CELERY_RESULT_BACKEND'],
        broker=app.config['CELERY_BROKER_URL']
    )
    celery.conf.update(app.config)

    class ContextTask(celery.Task):
        def __call__(self, *args, **kwargs):
            with app.app_context():
                return self.run(*args, **kwargs)

    celery.Task = ContextTask
    return celery

@celery.task()
def getImageInIntervall():
    DownloadImage()

class CreateDbTable(Server):
    def __call__(self, app, *args, **kwargs):
        databaseCreateTable()
        return Server.__call__(self, app, *args, **kwargs)

manager.add_command('runserver', CreateDbTable(host='0.0.0.0', port=8203))

if __name__ == "__main__":
    manager.run()