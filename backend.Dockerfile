# Sample Backend Dockerfile, DELETE ME AFTERWARDS
FROM oven/bun:1 AS base
WORKDIR /usr/src/app
COPY ./sample_backend/ .
RUN bun install
ENV PORT=3000
EXPOSE $PORT
CMD ["bun", "index.ts"]
