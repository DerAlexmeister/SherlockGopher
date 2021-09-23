import pytz
from sherlockanalysepictures.getmetadata import *
from flask import Flask
from flask_script import Manager
from flask_script import Server
from flask_apscheduler import APScheduler

class Config:
    SCHEDULER_API_ENABLED = True

app = Flask(__name__)
app.config.from_object(Config())

scheduler = APScheduler()
scheduler.init_app(app)

manager = Manager(app)

class CreateDbTable(Server):

    @scheduler.task("interval", id="image", seconds=20, timezone=pytz.UTC)
    def getImageInIntervall():
        DownloadImage()
    
    def __call__(self, app, *args, **kwargs):
        databaseDeleteTable()
        databaseCreateTable()
        scheduler.start()
        return Server.__call__(self, app, *args, **kwargs)

# create db on startup
manager.add_command('runserver', CreateDbTable(host='0.0.0.0', port=8203))

if __name__ == "__main__":
    manager.run()