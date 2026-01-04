pipeline {
    agent any

    environment {
        DOCKER_HUB = credentials('dockerhub')
        REPO = 'https://github.com/kvasianovych/kbot'
        BRANCH = 'develop'
    }

    parameters {
        choice(
            name: 'OS',
            choices: ['linux', 'darwin', 'windows'],
            description: 'Target operating system'
        )
        choice(
            name: 'ARCH',
            choices: ['amd64', 'arm64'],
            description: 'Target architecture'
        )
        booleanParam(
            name: 'SKIP_TESTS',
            defaultValue: false,
            description: 'Skip running tests'
        )
        booleanParam(
            name: 'SKIP_LINT',
            defaultValue: false,
            description: 'Skip running linter'
        )
    }

    stages {

        stage('clone') {
            steps {
                echo 'Clone Repository'
                git branch: "${BRANCH}", url: "${REPO}"
            }
        }

        stage('test') {
            steps {
                echo 'Testing started'
                sh "make test"
            }
        }

        stage('build') {
            steps {
                echo "Building binary started"
                sh "make build"
            }
        }

        stage('image') {
            steps {
                echo "Building image started"
                sh "make image"
            }
        }

        stage('login to docker') {
            steps {
                sh 'echo $DOCKER_HUB_PSW | docker login -u $DOCKER_HUB_USR --password-stdin'
            }
        }

        stage('push image') {
            steps {
              sh "make push"
            }
        }
    }
}
