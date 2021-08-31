import os

from flask import Flask
from constants import SERVER_PORT

app = Flask(__name__)

@app.route("/")
def index():
    return "hello world"

if __name__ == "__name__":
    app.run()