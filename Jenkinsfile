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
        stage("unit-test") {
            steps {
                sh 'make unit-tests'
            }
        }
    }
}
