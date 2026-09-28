"""Minimal service used to exercise the Python and container templates."""

import os

from flask import Flask, jsonify


def create_app() -> Flask:
    app = Flask(__name__)

    @app.get("/healthz")
    def healthz():
        return {"status": "ok"}

    @app.get("/version")
    def version():
        return jsonify(version=os.getenv("APP_VERSION", "dev"))

    @app.after_request
    def security_headers(resp):
        resp.headers["X-Content-Type-Options"] = "nosniff"
        resp.headers["Cache-Control"] = "no-store"
        return resp

    return app
