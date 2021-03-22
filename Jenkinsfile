pipeline {

  agent any

  stages {

    stage('Build And Push Docker Image') {

      environment {
        // Extract the username and password of our credentials into "DOCKER_CREDENTIALS_USR" and "DOCKER_CREDENTIALS_PSW".
        // (NOTE 1: DOCKER_CREDENTIALS will be set to "your_username:your_password".)
        // The new variables will always be YOUR_VARIABLE_NAME + _USR and _PSW.
        // (NOTE 2: You can't print credentials in the pipeline for security reasons.)
        DOCKER_CREDENTIALS = credentials('mrjosh-docker-credentials-id')
      } 

      steps {

        // Use a scripted pipeline.
        script {

          node {

            def app

            stage('Clone repository') {
              checkout scm
            }

            stage('Build image') {
              app = docker.build("homepi/homepi")
            }
            
            stage('Push image') {

              // Use the Credential ID of the Docker Hub Credentials we added to Jenkins.
              docker.withRegistry('https://registry.hub.docker.com', 'mrjosh-docker-credentials-id') {
                // Push the same image and tag it as the latest version (appears at the top of our version list).
                app.push("jenkins")
              }

            }

          }

        }

      }

    }

  }

}
