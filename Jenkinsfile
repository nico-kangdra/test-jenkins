pipeline {
    agent {label 'docker'}
    
    options {
        timeout(time: 1, unit: 'HOURS')
        timestamps()
        buildDiscarder(logRotator(numToKeepStr: '10'))
    }
    
    stages {
        stage('Load & Validate Pipeline Config') {
            steps {
                script {
                    try {
                        def config = readYaml file: 'pipeline.yaml'
                        env.PIPELINE_CONFIG = groovy.json.JsonOutput.toJson(config)
                        echo "✓ Pipeline config loaded successfully"
                        echo "Stages: ${config.stages.collect { it.name }.join(', ')}"
                    } catch (Exception e) {
                        error("Failed to load pipeline.yaml: ${e.message}")
                    }
                }
            }
        }
        
        stage('Execute Pipeline Stages') {
            steps {
                script {
                    def config = readYaml file: 'pipeline.yaml'
                    
                    config.stages.each { stageDef ->
                        executePipelineStage(stageDef)
                    }
                }
            }
        }
    }
    
    post {
        always {
            echo "Pipeline execution completed"
            cleanWs() 
        }
        failure {
            echo "❌ Pipeline failed - check logs above"
        }
        success {
            echo "✓ Pipeline completed successfully"
        }
    }
}

// Reusable function untuk execute individual stage
void executePipelineStage(Map stageDef) {
    try {
        timeout(time: stageDef.timeout ?: 30, unit: 'MINUTES') {
            retry(stageDef.retry ?: 0) {
                echo "\n═══════════════════════════════════════════"
                echo "Stage: ${stageDef.name}"
                
                // Check if stage should be skipped
                if (stageDef.when && !evaluateWhen(stageDef.when)) {
                    echo "⊘ Skipped (when condition not met)"
                    return
                }
                
                // Setup environment variables
                def envList = []
                if (stageDef.env) {
                    stageDef.env.each { envEntry ->
                        if (envEntry instanceof String) {
                            envList.add(envEntry)
                        } else if (envEntry instanceof Map) {
                            envEntry.each { k, v -> envList.add("${k}=${v}") }
                        }
                    }
                }
                
                // Run natively on agent (without Docker container)
                withEnv(envList) {
                    sh """
                        set -e
                        ${stageDef.command}
                    """
                }
                echo "✓ ${stageDef.name} completed successfully\n"
            }
        }
    } catch (Exception e) {
        if (stageDef.allowFailure) {
            echo "⚠ ${stageDef.name} failed but marked as allowFailure: ${e.message}"
            return
        }
        error("Stage '${stageDef.name}' failed: ${e.message}")
    }
}

// Helper function untuk evaluate 'when' conditions
boolean evaluateWhen(Map when) {
    if (when.branch) {
        return env.BRANCH_NAME == when.branch
    }
    if (when.tag) {
        return env.TAG_NAME != null && env.TAG_NAME != ''
    }
    if (when.build != null) {
        return when.build == true
    }
    if (when.expression) {
        return when.expression.call()
    }
    return true
}