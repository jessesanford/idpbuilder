# CI/CD Integration Examples

This document provides examples for integrating idpbuilder push into various CI/CD platforms.

## GitHub Actions

### Basic Workflow

```yaml
name: Build and Push

on:
  push:
    branches: [ main ]

jobs:
  push:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v2

      - name: Build image
        run: docker build -t myapp:${{ github.sha }} .

      - name: Push to registry
        env:
          IDPBUILDER_REGISTRY_USER: ${{ secrets.REGISTRY_USERNAME }}
          IDPBUILDER_REGISTRY_PASSWORD: ${{ secrets.REGISTRY_TOKEN }}
        run: idpbuilder push --quiet
```

### Multi-Platform Build and Push

```yaml
name: Multi-Platform Push

on:
  release:
    types: [published]

jobs:
  push:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Build multi-arch image
        run: |
          docker buildx create --use
          docker buildx build --platform linux/amd64,linux/arm64 \
            -t myapp:${{ github.ref_name }} .

      - name: Push image
        env:
          IDPBUILDER_REGISTRY_USER: ${{ secrets.REGISTRY_USERNAME }}
          IDPBUILDER_REGISTRY_PASSWORD: ${{ secrets.REGISTRY_TOKEN }}
        run: |
          idpbuilder push --timeout 15m
```

## GitLab CI

### Basic Pipeline

```yaml
stages:
  - build
  - push

build:
  stage: build
  script:
    - docker build -t $CI_REGISTRY_IMAGE:$CI_COMMIT_SHA .

push:
  stage: push
  variables:
    IDPBUILDER_REGISTRY_USER: $CI_REGISTRY_USER
    IDPBUILDER_REGISTRY_PASSWORD: $CI_REGISTRY_PASSWORD
  script:
    - idpbuilder push --quiet
  only:
    - main
```

### With Manual Approval

```yaml
push_production:
  stage: push
  variables:
    IDPBUILDER_REGISTRY_USER: $PROD_REGISTRY_USER
    IDPBUILDER_REGISTRY_PASSWORD: $PROD_REGISTRY_TOKEN
  script:
    - idpbuilder push --timeout 10m --verbose
  when: manual
  only:
    - tags
```

## Jenkins

### Declarative Pipeline

```groovy
pipeline {
    agent any

    environment {
        REGISTRY_CREDS = credentials('registry-credentials')
    }

    stages {
        stage('Build') {
            steps {
                sh 'docker build -t myapp:${BUILD_NUMBER} .'
            }
        }

        stage('Push') {
            steps {
                withCredentials([usernamePassword(
                    credentialsId: 'registry-credentials',
                    usernameVariable: 'IDPBUILDER_REGISTRY_USER',
                    passwordVariable: 'IDPBUILDER_REGISTRY_PASSWORD'
                )]) {
                    sh 'idpbuilder push --quiet'
                }
            }
        }
    }

    post {
        failure {
            echo 'Push failed - check logs'
        }
    }
}
```

### Scripted Pipeline with Retry

```groovy
node {
    stage('Build') {
        sh 'docker build -t myapp:latest .'
    }

    stage('Push') {
        withEnv([
            "IDPBUILDER_REGISTRY_USER=${env.REGISTRY_USER}",
            "IDPBUILDER_REGISTRY_PASSWORD=${env.REGISTRY_TOKEN}"
        ]) {
            retry(3) {
                sh 'idpbuilder push --retry-attempts 5'
            }
        }
    }
}
```

## CircleCI

### Basic Configuration

```yaml
version: 2.1

jobs:
  build-and-push:
    docker:
      - image: cimg/base:stable
    steps:
      - checkout
      - setup_remote_docker

      - run:
          name: Build image
          command: docker build -t myapp:${CIRCLE_SHA1} .

      - run:
          name: Push image
          command: |
            export IDPBUILDER_REGISTRY_USER=${REGISTRY_USER}
            export IDPBUILDER_REGISTRY_PASSWORD=${REGISTRY_TOKEN}
            idpbuilder push --quiet

workflows:
  build-deploy:
    jobs:
      - build-and-push:
          filters:
            branches:
              only: main
```

## Azure Pipelines

### Basic Pipeline

```yaml
trigger:
  branches:
    include:
      - main

pool:
  vmImage: 'ubuntu-latest'

steps:
  - task: Docker@2
    inputs:
      command: build
      Dockerfile: '**/Dockerfile'
      tags: $(Build.BuildId)

  - script: |
      export IDPBUILDER_REGISTRY_USER=$(registryUsername)
      export IDPBUILDER_REGISTRY_PASSWORD=$(registryPassword)
      idpbuilder push --quiet
    displayName: 'Push to registry'
    env:
      registryUsername: $(REGISTRY_USERNAME)
      registryPassword: $(REGISTRY_PASSWORD)
```

## Best Practices for CI/CD

### 1. Use Quiet Mode

Reduce log noise in CI/CD:

```bash
idpbuilder push --quiet
```

### 2. Set Appropriate Timeouts

Adjust for CI/CD environment:

```bash
idpbuilder push --timeout 10m
```

### 3. Handle Failures Gracefully

```yaml
- name: Push with error handling
  run: |
    if ! idpbuilder push --quiet; then
      echo "::error::Push failed"
      exit 1
    fi
```

### 4. Use Secrets for Credentials

Never hardcode credentials:

```yaml
env:
  IDPBUILDER_REGISTRY_USER: ${{ secrets.REGISTRY_USER }}
  IDPBUILDER_REGISTRY_PASSWORD: ${{ secrets.REGISTRY_TOKEN }}
```

### 5. Implement Retry Logic

For flaky network conditions:

```bash
idpbuilder push --retry-attempts 10 --timeout 15m
```

## See Also

- [Basic Examples](basic-push.md)
- [Advanced Examples](advanced-push.md)
- [Environment Variables Reference](../reference/environment-vars.md)
