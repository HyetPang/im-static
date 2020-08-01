node{
    stage('checkout'){
        git credentialsId: 'af2302ff-5927-4f27-bc8e-ace3600763b2', url: 'git@github.com/zengyu2020:im/im-static.git'
    }
    stage('complie'){
        def root = tool name: 'Go 1.14.3', type: 'go'
        withEnv(["GOROOT=${root}","PATH+GO=${root}/bin","GO111MODULE=on","CGO_ENABLED=0","GOPROXY=https://goproxy.cn,direct","GOPRIVATE=github.com/zengyu2020"]){
            sh 'go env'
            sh 'go mod download'
            sh 'go build -ldflags "-w -s" -o im-static-server main.go'
        }
    }
    stage('deploy'){
        sh 'supervisorctl restart im-static-server'
    }
}