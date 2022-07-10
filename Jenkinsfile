pipeline {
    agent any
    tools {
        go 'go-1.18'
    }
    environment {
        GOPATH = "${WORKSPACE}/go"
    }
    stages {
        stage('Prepare') {
            steps {
                sh 'go mod download'
            }
        }
        stage('Test') {
            steps {
                sh 'go test -v | /home/flohoss/go-junit-report > report.xml'
            }
        }
        stage('Integration') {
            steps {
                junit 'report.xml'
            }
        }
    }
}