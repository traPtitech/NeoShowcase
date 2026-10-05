// Simulates a static site build (e.g. an SPA) without any dependencies.
const fs = require("node:fs");

fs.mkdirSync("dist", { recursive: true });
fs.writeFileSync(
  "dist/index.html",
  "<!doctype html>\n<html>\n  <body>Hello, World!</body>\n</html>\n",
);
