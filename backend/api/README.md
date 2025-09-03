# API Documentation

This directory contains the OpenAPI spec for The Special Standard API.

- `openapi.yaml`: OpenAPI 3.0 specification
- Swagger UI: http://localhost:8080/swagger (should autoroute to http://localhost:8080/swagger/index.html) when running the backend

To update the OpenAPI spec:
1. Edit `openapi.yaml`
2. Restart the backend to see changes in Swagger UI. Changes made to `openapi.yaml` will not be automatically reflected in Swagger UI, as the spec is served statically at http://localhost:8080/api/openapi.yaml

## Why we are enforcing this contract
The OpenAPI spec acts as a single source of truth for our API. It is a formal way to describe API endpoints, request/response schemas, and error details, allowing anyone to understand and use our API without looking at the actual code. It also helps verify that we are meeting project requirements.

Here are some additional reasons:
- Pretty documentation via Swagger
- Type-safe APIs. When we begin frontend development, there is a lot of tooling we can use to create type-safe APIs. 
    - https://openapi-ts.dev
    - https://github.com/ferdikoomen/openapi-typescript-codegen

You can read the official OpenAPI specification here: [OpenAPI Specification](https://swagger.io/specification/). 
# Tips and best practices
- Keep the spec updated whenever endpoints or data schemas change to ensure consistency across frontend, backend, and documentation
- You can import the OpenAPI spec into Postman for easy testing setup
- You can use the Swagger UI to view endpoint details and make test requests directly from the browser
