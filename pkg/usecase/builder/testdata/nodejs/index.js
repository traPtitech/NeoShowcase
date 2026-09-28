const http = require("node:http");

const port = process.env.PORT || 8080;
http
  .createServer((_req, res) => {
    res.end("Hello, World!\n");
  })
  .listen(port);
