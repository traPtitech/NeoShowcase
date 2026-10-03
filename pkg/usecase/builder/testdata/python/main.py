import os
from http.server import BaseHTTPRequestHandler, HTTPServer


class Handler(BaseHTTPRequestHandler):
    def do_GET(self):
        self.send_response(200)
        self.end_headers()
        self.wfile.write(b"Hello, World!\n")


HTTPServer(("", int(os.environ.get("PORT", "8080"))), Handler).serve_forever()
