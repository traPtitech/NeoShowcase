FROM ghcr.io/astral-sh/uv:python3.13-bookworm-slim@sha256:531f855bda2c73cd6ef67d56b733b357cea384185b3022bd09f05e002cd144ca
WORKDIR /app
COPY . .
RUN uv sync --locked
CMD ["uv", "run", "main.py"]
