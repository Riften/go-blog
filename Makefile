.PHONY:blog
.PHONY:blog.exe
.PHONY:installer.exe
blog.exe:
	go build -o blog.exe
blog:
	go build -o blog
installer.exe:
	go build -o installer.exe installer/main.go