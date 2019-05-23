#!/usr/bin/env groovy

@Library('jenkins-build-utils')_

def utils = new com.pureport.Utils()

def version = "0.2.0"
def plugin_name = "terraform-provider-pureport"

pipeline {
    agent {
      docker {
        image 'golang:1.12'
      }
    }
    options {
        disableConcurrentBuilds()
    }
    parameters {
      booleanParam(
          name: 'ACCEPTANCE_TESTS_RUN',
          defaultValue: false,
          description: 'Should we run the acceptance tests as part of run?'
          )
      booleanParam(
          name: 'ACCEPTANCE_TESTS_LOG_TO_FILE',
          defaultValue: true,
          description: 'Should debug logs be written to a separate file?'
          )
      choice(
          name: 'ACCEPTANCE_TESTS_LOG_LEVEL',
          choices: ['TRACE', 'DEBUG', 'INFO', 'WARN', 'ERROR'],
          description: 'The Terraform Debug Level'
          )
    }
    environment {
        TF_LOG              = "${params.ACCEPTANCE_TESTS_LOG_LEVEL}"
        TF_LOG_PATH         = "${params.ACCEPTANCE_TESTS_LOG_TO_FILE ? 'tf_log.log' : '' }"
        GOPATH              = "/go"
        GOCACHE             = "/tmp/go/.cache"
        PUREPORT_ENDPOINT   = "https://dev1-api.pureportdev.com"
        PUREPORT_API_KEY    = credentials('terraform-pureport-dev1-api-key')
        PUREPORT_API_SECRET = credentials('terraform-pureport-dev1-api-secret')
        GOOGLE_CREDENTIALS  = credentials('terraform-google-credentials-id')
        GOOGLE_PROJECT      = "pureport-customer1"
        GOOGLE_REGION       = "us-west2"
    }
    stages {
        stage('Configure') {
            steps {
                script {

                    plugin_name += "_v${version}"

                    // Only add the build version for the develop branch
                    if (env.BRANCH_NAME == "develop") {
                      plugin_name += "-b${env.BUILD_NUMBER}"
                    }

                }
            }
        }
        stage('Build') {
            steps {

                retry(3) {
                  sh "make"
                  sh "make plugin"
                  sh "mv terraform-provider-pureport ${plugin_name}"

                  archiveArtifacts(
                      artifacts: "${plugin_name}"
                      )
                }
            }
        }
        stage('Run Terraform Tests') {
            when {

                // This can take a long time so we may only want to do this on develop
                anyOf {
                  branch 'develop'
                  expression { return params.ACCEPTANCE_TESTS_RUN }
                }
            }
            steps {

                script {

                    // Don't fail if the test fall. Just setting this until we can get our issues
                    // resolved with the Google Provider.
                    sh "make testacc"

                    archiveArtifacts(
                        allowEmptyArchive: true,
                        artifacts: 'pureport/tf_log.log'
                        )
                }
            }
        }
        stage('Copy plugin to Nexus') {
            when {

                // This can take a long time so we may only want to do this on develop
                anyOf {
                  branch 'develop'
                  branch 'master'
                }
            }
            steps {
                script {
                    withCredentials([
                        usernamePassword(
                          credentialsId: 'nexus_credentials',
                          usernameVariable: 'nexusUsername',
                          passwordVariable: 'nexusPassword'
                          )
                    ]) {

                      def nexus_url = "https://nexus.dev.pureport.com/repository/terraform-provider-pureport/${env.BRANCH_NAME}/"

                      sh "curl -v -u ${nexusUsername}:${nexusPassword} --upload-file ${plugin_name} ${nexus_url}"

                      // Set the description text for the job
                      currentBuild.description = "Version: ${plugin_name}"

                    }
                }
            }
        }
    }
    post {
        success {
            slackSend(color: '#30A452', message: "SUCCESS: <${env.BUILD_URL}|${env.JOB_NAME}#${env.BUILD_NUMBER}>")
        }
        unstable {
            slackSend(color: '#DD9F3D', message: "UNSTABLE: <${env.BUILD_URL}|${env.JOB_NAME}#${env.BUILD_NUMBER}>")

            script {
                utils.sendUnstableEmail()
            }
        }
        failure {
            slackSend(color: '#D41519', message: "FAILED: <${env.BUILD_URL}|${env.JOB_NAME}#${env.BUILD_NUMBER}>")
            script {
                utils.sendFailureEmail()
            }
        }
    }
}