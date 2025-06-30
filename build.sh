/usr/local/go/bin/go install google.golang.org/protobuf/cmd/protoc-gen-go@latest > /dev/null 2>&1
/usr/local/go/bin/go get github.com/shirou/gopsutil/v3/... > /dev/null 2>&1
/usr/local/go/bin/go get fyne.io/fyne/v2 > /dev/null 2>&1
export PATH=$PATH:/usr/local/go/bin > /dev/null 2>&1
/usr/local/go/bin/go mod vendor > /dev/null 2>&1
/usr/local/go/bin/go mod tidy > /dev/null 2>&1
protoc --go_out=. --go_opt=paths=source_relative consenbench/common/message.proto > /dev/null 2>&1
/usr/local/go/bin/go build -v -o ./consenbench/bin/bench ./consenbench/ > /dev/null 2>&1