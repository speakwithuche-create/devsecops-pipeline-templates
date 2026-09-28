from app import create_app


def client():
    return create_app().test_client()


def test_health():
    assert client().get("/healthz").json == {"status": "ok"}


def test_version_from_env(monkeypatch):
    monkeypatch.setenv("APP_VERSION", "1.2.3")
    assert client().get("/version").json == {"version": "1.2.3"}


def test_security_headers_on_every_response():
    for path in ("/healthz", "/version", "/missing"):
        r = client().get(path)
        assert r.headers["X-Content-Type-Options"] == "nosniff"
        assert r.headers["Cache-Control"] == "no-store"
