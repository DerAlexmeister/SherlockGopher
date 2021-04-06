import os
from sherlockanalysepictures.getmetadata import *
from flask import Flask
from flask_script import Manager
from flask_script import Server
from celery import Celery

app = Flask(__name__)

# prepare redis and flask url
def init():
    add = readFromENV("FLASKA_URL", "0.0.0.0")
    uri = "redis://" + add + ":6379"
    return uri

# allows docker usage
def readFromENV(key, defaultVal):
    value = os.environ[key]
    if value == "":
        return defaultVal
    return value

val = init()

app.config.update(
    CELERY_BROKER_URL=val,
    CELERY_RESULT_BACKEND=val
)

# create celery app
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

celery = make_celery(app)
manager = Manager(app)

# background task to download images
@celery.task(name='tasks.getImageInIntervall')
def getImageInIntervall():
    DownloadImage()

class CreateDbTable(Server):
    def __call__(self, app, *args, **kwargs):
        databaseCreateTable()
        return Server.__call__(self, app, *args, **kwargs)

# create db on startup
manager.add_command('runserver', CreateDbTable(host='0.0.0.0', port=8203))

if __name__ == "__main__":
    manager.run()