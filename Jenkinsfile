pipeline {
    agent  any 

    stages {
        stage('Test Repo Vuln') {
                
            agent { label "trivy" }

            steps {
                sh 'trivy fs --timeout 14m --severity HIGH,CRITICAL .'
            }
        }

        stage('Build Image') {

            agent { label "docker"}

            steps {
                script {
                        withDockerRegistry(credentialsId: 'dockerhub') {
                               sh "docker build -t huan271/tickeria:latest ."
                        }
                }
            }
        }

        stage('Docker Image Scan') {

            agent { label "trivy"}

            steps {
                sh "trivy image --format table -o trivy-image-report.html huan271/tickeria:latest "
            }
        }

        stage('Push docker registry') {

            agent { label "docker"}

            steps {
                script {
                        withDockerRegistry(credentialsId: 'dockerhub') {
                               sh "docker push huan271/tickeria:latest"
                        }
                }
            }
        }
    }
}