module.exports = {
  petstore: {
    output: {
      mode: "tags",
      target: "src/lib/api/api.ts",
      client: "react-query",
      mock: false,
    },
    input: {
      target: "../backend/api/openapi.yaml",
      validation: false,
    },
  },
};
