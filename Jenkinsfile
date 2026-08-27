pipeline {
    agent  none 

    stages {
        stage('Test trivy') {
                
            agent { label "trivy" }

            steps {
                sh 'ls -alth'
                sh 'trivy fs --timeout 14m --severity HIGH,CRITICAL .'
            }
        }

        stage('Test golang') {

            agent { label "golang"}

            steps {
                sh 'go version'
                sh 'ls -alth'
            }
        }

        stage('Test docker registry') {

            agent { label "docker"}

            steps {
                script {
                        withDockerRegistry(credentialsId: 'dockerhub') {
                               sh "docker build -t huan271/tickeria:latest ."
                            }
                        }
            }
        }
    }
}