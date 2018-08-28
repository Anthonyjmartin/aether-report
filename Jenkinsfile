node('compose') {
    try{
      stage('Checkout'){
          echo 'Checking out SCM'
          checkout scm
      }
      sh "sed -i 's/!release!/${env.BUILD_ID}/g' rpm.json"
      env.appVer = sh (returnStdout: true, script: "./version.sh").trim()
      env.PDIR = "/go/src/gitlab.com/anthony.j.martin/aether-report/"

      docker.image('golang:1.10').inside('-v ${WORKSPACE}:$PDIR') {
        stage('Build initial image'){
          sh "cd ${env.PDIR} && go get github.com/golang/dep/cmd/dep"
          sh "cd ${env.PDIR} && go get github.com/golang/lint/golint"
          sh "cd ${env.PDIR} && go get github.com/tebeka/go2xunit"
          sh "cd ${env.PDIR} && dep ensure"
        }
        stage('Test and Lint'){
          sh "cd ${env.PDIR} && go test -cover ./..."
          sh "cd ${env.PDIR} && go tool vet ./"
          sh "cd ${env.PDIR} && golint ./..."
        }
        stage('Build Executable'){
          sh "cd ${env.PDIR} && go build -a -v -o build/linux/amd64/aether-report -ldflags '-X main.version=${env.appVer}' cmd/aether-report/main.go"
        }
      }
      docker.image('registry.gitlab.com/anthony.j.martin/aether-report/rpm-builder:latest').inside('-v ${WORKSPACE}:$PDIR') {
        stage('Build RPM'){
          sh "cd ${env.PDIR} && go-bin-rpm generate -a amd64 -o build/linux/amd64/aether-report-${env.appVer}.rpm --version ${env.appVer}; exit 0"
        }
      }
    }catch (e) {
        echo "Caught: ${e}"
        currentBuild.result = "FAILED"
    }
}
