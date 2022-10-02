BINARY="ssh-tool"

default:   # 这下面都要以tab开头  4个空格不行
	@echo "build the ${BINARY}"
	@GOOS=linux GOARCH=amd64 go build -o  build/${BINARY}.linux  -tags=jsoniter
	@go build -o  build/${BINARY}.mac  -tags=jsoniter
	@echo "build done."