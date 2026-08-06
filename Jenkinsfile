pipeline {
    agent {
        label 'golang'
    }

    stages {
        stage('Checkout') {
            steps {
                sh 'ls -alth'
            }
        }

        stage('Go Version') {
            steps {
                sh 'go version'
            }
        }
    }
}