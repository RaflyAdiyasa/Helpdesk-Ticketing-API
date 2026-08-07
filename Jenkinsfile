pipeline {
    agent  none 

    stages {
        stage('Test trivy') {
                
            agent { label "trivy" }

            steps {
                sh 'ls -alth'
                sh 'trivy fs -f json -o result.json .'
            }
        }

        stage('Test golang') {
            steps {
                sh 'go version'
                sh 'ls -alth'
            }
        }
    }
}