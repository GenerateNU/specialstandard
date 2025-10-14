module.exports = {
  // Generate API functions only
  api: {
    output: {
      mode: 'tags',
      target: 'src/lib/api',
      client: 'axios',
      override: {
        mutator: {
          path: 'src/lib/api/apiClient.ts',
          name: 'customAxios',
        },
        naming: 'camelCase',
      },
      prettier: true,
    },
    input: {
      target: '../backend/api/openapi.yaml',
    },
  },
}
