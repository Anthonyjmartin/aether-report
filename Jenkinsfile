node {
    try{
      stage('Checkout'){
          echo 'Checking out SCM'
          checkout scm
      }

      stage('Build GO image') {
        def majMinVer = sh returnStdout: true, script: """cat ./baseversion"""
        def appVer = sh returnStdout: true, script: """echo ${env.majMinVer}.${env.BUILD_ID}"""
        def testImage = docker.build("aether-report:${env.appVer}", "--build-arg version=${env.appVer}")
      }

      stage('Test'){
        testImage.inside("-v $PWD:/go/src/gitlab.com/anthony.j.martin/aether-report"){
          //List all our project files with 'go list ./... | grep -v /vendor/ | grep -v github.com | grep -v golang.org'
          def paths = sh 'go list ./... | grep -v /vendor/ | grep -v github.com | grep -v golang.org'

          echo 'Vetting'

          sh """cd $GOPATH/src && go tool vet ${paths}"""

          echo 'Linting'
          sh """cd $GOPATH/src && golint ${paths}"""

          echo 'Testing'
          sh """cd $GOPATH/src && go test -cover ${paths}"""
        }
      }

      stage('Build'){
        testImage.inside("-v $PWD:/go/src/gitlab.com/anthony.j.martin/aether-report"){
          echo 'Building Executable'

          //Produced binary is $GOPATH/src/cmd/project/project
          sh """cd $GOPATH/src/gitlab.com/anthony.j.martin/aether-report && go build -a -v -o build/linux/amd64/aether-report -ldflags '-X main.version=${env.appVer}' cmd/aether-report/main.go"""
        }
        def rpmImage = docker.build("aether-report:${env.appVer}", "-f ./Dockerfile.RPM --build-arg version=${env.appVer}")
        rpmImage.inside("-v $PWD:/go/src/gitlab.com/anthony.j.martin/aether-report"){
          sh """cd $GOPATH/src/gitlab.com/anthony.j.martin/aether-report && go-bin-rpm generate -a amd64 -o build/linux/amd64/aether-report-${env.appVer}.rpm --version ${env.appVer}"""
        }
      }
    }catch (e) {
        // If there was an exception thrown, the build failed
        currentBuild.result = "FAILED"
    }
}
