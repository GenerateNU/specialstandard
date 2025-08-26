const server = Bun.serve({
  port: 3000,
  fetch(request) {
    return new Response("Sample backend is up and running!");
  },
});

console.log(`Listening on ${server.url}`);
