.PHONY: install dev build start test lint clean docker-build docker-run

# Install dependencies
install:
	npm install

# Run development server
dev:
	npm run dev

# Build for production
build:
	npm run build

# Start production server
start:
	npm run start

# Run tests (if you have them)
test:
	npm run test

# Lint code
lint:
	npm run lint

# Lint code and fix issues
lint-fix:
	npm run lint:fix

# Docker operations
docker-build:
	docker build -t specialstandard-frontend .

docker-run:
	docker run -p 3000:3000 specialstandard-frontend

# Clean
clean:
	rm -rf node_modules .next out