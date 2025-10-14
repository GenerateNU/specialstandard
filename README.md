# The Special Standard

An all-in-one digital curriculum platform designed for school districts and speech-language pathologists to simplify therapy management, progress tracking, and compliance. For students, it provides interactive and engaging activities that make learning fun and effective.

## 🎯 Overview

The Special Standard revolutionizes speech-language therapy in educational settings by providing:

- **For Therapists**: Streamlined therapy management, automated progress tracking, and compliance reporting
- **For Students**: Interactive, gamified learning experiences that adapt to individual needs

## ✨ Features

- **Digital Curriculum Management**: Comprehensive library of evidence-based therapy activities
- **Progress Tracking**: Real-time monitoring of student progress with detailed analytics
- **Compliance Tools**: Automated IEP goal tracking and report generation
- **Interactive Activities**: Engaging, game-based learning experiences for students
- **Multi-Platform Support**: Web-based platform accessible on any device
- **Secure Data Management**: FERPA and HIPAA compliant data handling

Here's a modular section you can drop into your README:

## 🛠️ Development

### Quick Start

Start all services with Docker and hot reload:

```bash
make dev
```

Open [http://localhost:3000](http://localhost:3000) to view the application.

### Common Commands

| Command            | Description                        |
| ------------------ | ---------------------------------- |
| `make dev`         | Start all services with hot reload |
| `make test`        | Run all tests                      |
| `make lint`        | Check code quality                 |
| `make lint-fix`    | Auto-fix linting issues            |
| `make docker-down` | Stop all services                  |
| `make docker-logs` | View service logs                  |

Run `make help` to see all available commands.

## 📁 File Structure

```bash
├── .github                      # GitHub specific files
│   ├── pull_request_template.md # PR template
│   └── workflows               # CI/CD workflows
│       ├── backend-ci.yml      # Backend CI pipeline
│       └── backend-deploy.yml  # Backend deployment
├── backend                     # Backend source code (Go)
├── backend.Dockerfile          # Backend dockerfile
├── frontend                    # Frontend source code
├── docker-compose.yml          # Docker compose configuration
├── CONTRIBUTING.md            # Contribution guidelines
├── LICENSE                    # Project license
└── README.md                  # This file
```
