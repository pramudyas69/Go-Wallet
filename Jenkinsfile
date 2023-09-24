pipeline {
    agent any
    tools {
        go 'Go1.14'
    }
    environment {
        GO111MODULE = 'on'
        CGO_ENABLED = '0'
        GOPATH = "${JENKINS_HOME}/workspace/${JOB_NAME}/builds/${BUILD_NUMBER}"
    }
    stages {
        stage("Prepare Environment") {
            steps {
                script {
                    // Instalasi make jika belum terinstal
                    def makeInstalled = sh(script: 'which make', returnStatus: true)
                    if (makeInstalled != 0) {
                        echo "Installing make..."
                        sh 'sudo apt-get update && sudo apt-get install -y make'
                    }
                }
            }
        }
        stage("unit-test") {
            steps {
                sh 'make unit-tests'
            }
        }
    }
}
