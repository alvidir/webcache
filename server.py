import logging

from flask import Flask
app = Flask(__name__)

def run(port: str):
    if not port and len(port) == 0:
        raise Exception("a valid port must be set")

    app.run(host="0.0.0.0", port=port)

@app.route("/")
def index():
    return "hello world"