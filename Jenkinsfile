@Library('credo_jenkins_library') _

def artifactId, artifactType

pipeline {

    options {
        buildDiscarder(logRotator(numToKeepStr: '30'))
    }

    environment {
        LAMBDA_NAME = "coverage"
        GOROOT = tool name: 'CREDO GOLANG 1.11', type: 'go'
        GOPATH = "${env.WORKSPACE}"
        PATH = "${GOPATH}/bin:${GOROOT}/bin:${PATH}"
        WORKDIR = "${GOPATH}/src/bitbucket.org/credomobile/${LAMBDA_NAME}"
    }

    agent {
        node {
            label "master"
        }
    }

    stages {

        stage('Checkout') {
            steps {
                dir("${env.WORKDIR}") {
                    checkout scm
                }
            }
        }

        stage('Validate') {
            steps {
                dir("${env.WORKDIR}/infrastructure/terraform") {
                    sh 'terraform init -backend=false'
                    sh 'terraform validate -check-variables=false'
                }
                script {
                    def ansiblePlaybookExists = fileExists("infrastructure/ansible/vault-secrets.yml")
                    if (ansiblePlaybookExists) {
                        dir("${env.WORKDIR}/infrastructure/ansible") {

                            sh 'pipenv install && pipenv run ansible-playbook vault-secrets.yml --syntax-check'
                        }
                    }
                }
            }

            post {
                success {
                    echo 'Package validation passed'
                }
                failure {
                    echo 'Package validation failed'
                    slackSend channel: '#infra_chatops', color: 'danger', message: "Lambda Package validation FAILED. Job: `${env.JOB_NAME}` (<${env.BUILD_URL}|#${env.BUILD_NUMBER}>)"
                }
            }
        }

        stage('Install Dependencies') {
            steps {
                dir("${env.WORKDIR}") {
                    sh 'go version'
                    sh 'go get -u github.com/jstemmer/go-junit-report'
                    sh 'curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh'
                    sh 'dep ensure -v'
                }
            }
        }

        stage('Run Unit tests...') {
            steps {
                dir("${env.WORKDIR}") {
                    sh 'mkdir test-results'
                    sh 'go test ./... -race -cover -v 2>&1 | go-junit-report > test-results/report.xml'
                }
            }
            post {
                always {
                    dir("${env.WORKDIR}") {
                        junit 'test-results/report.xml'
                    }
                }
            }

        }

        stage('Build and Package...') {
            steps {
                dir("${env.WORKDIR}") {
                    sh "make build"
                    script {
                        artifactId   = credoArtifactId(env)
                        artifactType = env.BRANCH_NAME == 'master' ? 'release' : 'snapshot'
                    }
                }
            }
            post {
                success {
                    dir("${env.WORKDIR}") {
                        script {
                            withAWS(region: 'us-east-1') {
                                s3Upload file: "${LAMBDA_NAME}.zip", bucket: "credo-master-artifacts-lambda", path: "${LAMBDA_NAME}/${artifactId}", metadatas:["${artifactType}:true"]
                            }
                        }
                        echo 'Package uploaded to S3 successfully...'

                        sh "git tag ${artifactId}"
                        sh "git push origin ${artifactId}"
                    }
                }
                failure {
                    echo 'Build and Package step failed'
                    slackSend channel: '#infra_chatops', color: 'danger', message: "Lambda Build and Package validation FAILED. Job: `${env.JOB_NAME}` (<${env.BUILD_URL}|#${env.BUILD_NUMBER}>)"
                }
                unstable {
                    echo 'Build and Package step failed'
                    slackSend channel: '#infra_chatops', color: 'danger', message: "Lambda Build and Package validation FAILED. Job: `${env.JOB_NAME}` (<${env.BUILD_URL}|#${env.BUILD_NUMBER}>)"
                }
            }
        }

    }
    post {
        always {
            cleanWs()
            dir("${env.WORKSPACE}@libs") {
                deleteDir()
            }
            dir("${env.WORKSPACE}@tmp") {
                deleteDir()
            }
        }
    }
}