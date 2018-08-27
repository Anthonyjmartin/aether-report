node {
    try{
      stage('Checkout'){
          echo 'Checking out SCM'
          checkout scm
      }
      env.majMinVer = sh (returnStdout: true, script: "./version.sh").trim()
      env.appVer = sh (returnStdout: true, script: "echo ${env.majMinVer}.${env.BUILD_ID}").trim()
      def testImage = docker.build("aether-report", "--build-arg version='${env.appVer}' .")

      testImage.inside("-v $WORKSPACE/:/go/src/gitlab.com/anthony.j.martin/aether-report -w /go/src/gitlab.com/anthony.j.martin/aether-report") {
        // sh "pwd; echo; ls"
        // stage('Test'){
          // sh "pwd; echo; ls"
          //List all our project files with 'go list ./... | grep -v /vendor/ | grep -v github.com | grep -v golang.org'
        env.paths = sh (returnStdout: true, script: 'cd /go/src/gitlab.com/anthony.j.martin/aether-report && go list ./... | grep -v /vendor/ | grep -v github.com | grep -v golang.org')

        echo 'Vetting'

        sh "cd /go/src && go tool vet ${env.paths}"

        echo 'Linting'
        sh "cd /go/src && golint ${env.paths}"

        echo 'Testing'
        sh "cd /go/src && go test -cover ${env.paths}"
        // }
      }

      testImage.inside("-v .:/go/src/gitlab.com/anthony.j.martin/aether-report -w /go/src/gitlab.com/anthony.j.martin/aether-report") {
        stage('Build Executable'){
          //Produced binary is $GOPATH/src/cmd/project/project
          sh "cd /go/src/gitlab.com/anthony.j.martin/aether-report && go build -a -v -o build/linux/amd64/aether-report -ldflags '-X main.version=${env.appVer}' cmd/aether-report/main.go"
        }
      }
      def rpmImage = docker.build("aether-report-rpm", "--build-arg version=${env.appVer} ./Dockerfile.RPM")
      rpmImage.inside("-v $WORKSPACE:/go/src/gitlab.com/anthony.j.martin/aether-report"){
        sh "cd $GOPATH/src/gitlab.com/anthony.j.martin/aether-report && go-bin-rpm generate -a amd64 -o build/linux/amd64/aether-report-${env.appVer}.rpm --version ${env.appVer}"
      }
    }catch (e) {
        // If there was an exception thrown, the build failed
        echo "Caught: ${e}"
        currentBuild.result = "FAILED"
    }
}
