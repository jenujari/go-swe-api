pipeline {
    agent any

    triggers {
        githubPush()
    }

    environment {
        IMAGE_NAME = "docker.io/jhon5456/swegrpc"
        REGISTRY_CREDENTIALS = "docker-creds"
    }

    stages {

        stage('Checkout Code') {
            steps {
                checkout scm
            }
        }

        stage('Build for Tests') {
            steps {
                sh '''
                    echo "Building test image..."
                    podman run --rm -v .:/workspace -w /workspace docker.io/bufbuild/buf generate proto
                    podman compose -f compose.yaml build test_sweapi
                '''
            }
        }

        stage('Run Linting and Tests') {
            steps {
                sh '''
                    echo "Running golangci-lint..."
                    podman compose -f compose.yaml run --entrypoint sh --rm test_sweapi -c "go tool golangci-lint run -v"

                    echo "Running tests..."
                    podman compose -f compose.yaml \
                        run --entrypoint sh --rm test_sweapi -c \
                        "go test ./... -v -coverpkg=./... -coverprofile=coverage.out \
                        && grep -v '\\.pb\\.go' coverage.out > cov.tmp \
                        && mv cov.tmp coverage.out \
                        && go tool cover -func=coverage.out | tee coverage-summary.txt" \
                        && go tool cover -html=coverage.out -o coverage.html
                '''
            }
        }

        stage('Validate branch is main before build') {
            steps {
                script {
                    def branch = env.GIT_BRANCH?.replace("origin/", "")
                    env.BRANCH_NAME = branch

                    echo "Triggered branch: ${branch}"

                    if (branch != "main") {

                        currentBuild.result = 'NOT_BUILT'

                        echo """
                        Skipping futher build steps in pipeline.
                        Push detected on non-main branch:
                        ${branch}
                        """

                        return
                    }

                    echo "Main branch detected. Continuing pipeline..."
                }
            }
        }

        stage('Generate Commit SHA') {
            steps {
                script {
                    SHORT_SHA = sh(
                        script: 'git rev-parse --short HEAD',
                        returnStdout: true
                    ).trim()

                    env.IMAGE_TAG = SHORT_SHA

                    echo "Image tag set to: $IMAGE_TAG"
                }
            }
        }

        stage('Build Podman Image') {
            when {
                expression {
                    env.BRANCH_NAME == 'main'
                }
            }
            steps {
                sh """
                    podman build -t $IMAGE_NAME:$SHORT_SHA .
                    podman tag $IMAGE_NAME:$SHORT_SHA $IMAGE_NAME:latest
                """
            }
        }

        stage('Docker Login') {
            when {
                expression {
                    env.BRANCH_NAME == 'main'
                }
            }
            steps {
                withCredentials([
                    usernamePassword(
                        credentialsId: "${REGISTRY_CREDENTIALS}",
                        usernameVariable: 'DOCKER_USER',
                        passwordVariable: 'DOCKER_PASS'
                    )
                ]) {
                    sh 'echo $DOCKER_PASS | podman login docker.io -u $DOCKER_USER --password-stdin'
                }
            }
        }

        stage('Push Docker Image') {
            when {
                expression {
                    env.BRANCH_NAME == 'main'
                }
            }
            steps {
                sh """
                    podman push $IMAGE_NAME:$SHORT_SHA
                    podman push $IMAGE_NAME:latest
                """
            }
        }

        stage('Update Kubernetes Deployment') {
            when {
                expression {
                    env.BRANCH_NAME == 'main'
                }
            }
            steps {
                sh """
                    kubectl set image deployment/sweapi \
                    sweapi=$IMAGE_NAME:latest

                    kubectl rollout status deployment/sweapi
                """
            }
        }
    }

    post {
        success {
            echo 'Deployment successful!'
        }

        failure {
            echo 'Pipeline failed!'
        }

        notBuilt {
            echo 'Pipeline skipped for non-main branch'
        }

        always {
            echo 'Cleaning workspace...'
            archiveArtifacts artifacts: '''
              coverage.out,
              coverage-summary.txt,
              coverage.html
            ''',
            fingerprint: true

            publishHTML([
                allowMissing: false,
                alwaysLinkToLastBuild: true,
                keepAll: true,
                reportDir: '.',
                reportFiles: 'coverage.html',
                reportName: 'Go Coverage Report'
            ])

            recordCoverage(
                sourceCodeRetention: 'EVERY_BUILD',
                tools: [[parser: 'GO_COV', pattern: 'coverage.out']]
            )

            cleanWs()
        }
    }
}
