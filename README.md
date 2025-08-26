# Software Fall 2025 Template

You are free to delete, remove, and shape the layout of your repository in whichever way you like. The provided setup files are there to use at your own convenience.

## Table of Contents

1. [File Structure](#file-structure)
2. [Deployments](#deployments)

### File-Structure

The template repository is laid out as follows below.

```bash
├── .github # Place workflows here
│   ├── pull_request_template.md
│   └── workflows
│       └── backend-deploy.yml
├── backend # Backend source code
├── backend.Dockerfile # Backend dockerfile
├── CONTRIBUTING.md # Contribution documentation for engineers
├── frontend # Frontend source code
├── LICENSE
├── README.md
└── sample_backend # DELETE THIS, PURELY FOR TESTING PURPOSES ONLY
    ├── .gitignore
    ├── bun.lock
    ├── index.ts
    ├── package.json
    ├── README.md
    └── tsconfig.json
```

### Deployments

A sample deployment will be found in the workflows folder. Each workflow calls a reusable workflow upstream to deploy container images to digital ocean.

Please provide the following repository secrets to utilize the deployment workflow:

- `DO_TOKEN` _For authentication into the digital ocean container registry_

You are not _required_ to utilize this workflow, indeed you can delete the provided workflow if you would like to employ your own deployment structure.

Additionally you will need to modify the caller script to provide a repo name of your choice. An example is given to you here:

```yaml
jobs:
  deploy:
    uses: GenerateNU/shiperate/.github/workflows/backend-deploy.yml@main # Do not change this
    with:
      context: . # Current working directory
      dockerfile: backend.Dockerfile # Path to the dockerfile relative to the current working directory
      repo: registry.digitalocean.com/gen-sw-fall-2025/test # The last slash is your repository name
      tag: latest # Whatever tag you want, but make sure your deployment platform is setup to listen to the necessary tags.
    secrets: inherit
```

Good Luck!
