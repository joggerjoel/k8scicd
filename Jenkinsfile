pipeline {
    agent any
    environment {
        registry = "joggerjoel/k8scicd"
        GOCACHE = "/tmp"
    }
    stages {
        stage('Build') {
            agent { 
                docker { 
                    image 'golang' 
                }
            }
            steps {
                // Create our project directory.
                sh 'cd ${GOPATH}/src'
                sh 'mkdir -p ${GOPATH}/src/hello-world'
                // Copy all files in our Jenkins workspace to our project directory.
                echo 'pwd'
                sh 'pwd'
                sh 'ls ..'
                
                sh 'cp -r ${WORKSPACE}/* ${GOPATH}/src/hello-world'
                // Build the app.
                sh 'rm -f go.mod'
                sh 'go mod init hello-world'
                sh 'go mod tidy'
                sh 'go build'              
                
            }     
        }
        stage('Test') {
            agent { 
                docker { 
                    image 'golang' 
                }
            }
            steps {                 
                // Create our project directory.
                sh 'cd ${GOPATH}/src'
                sh 'mkdir -p ${GOPATH}/src/hello-world'
                // Copy all files in our Jenkins workspace to our project directory.                
                sh 'cp -r ${WORKSPACE}/* ${GOPATH}/src/hello-world'
                echo 'pwd'
                sh 'pwd'
                sh 'ls ..'

                // Remove cached test results.
                sh 'go clean -cache'
                // Run Unit Tests.
                sh 'go test ./... -v -short'            
            }
        }
        stage('Publish') {
            environment {
                registryCredential = 'dockerhub'
            }
            steps{
                script {
                    def appimage = docker.build registry + ":$BUILD_NUMBER"
                    docker.withRegistry( '', registryCredential ) {
                        appimage.push()
                        appimage.push('latest')
                    }
                }
            }
        }
        stage ('Deploy') {
            steps {
                script{
                    def image_id = registry + ":$BUILD_NUMBER"
                    sh "ansible-playbook  playbook.yml --extra-vars \"image_id=${image_id}\""
                }
            }
        }
    }
}
