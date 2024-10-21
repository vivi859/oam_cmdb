del OAM.exe
del oam.tar.gz
del oam.zip
go env -w CGO_ENABLED=0

if "%1" == "" (
    bee pack -exp=screenshot:logs:.vscode:.git:README.md:tests:build.bat:oam.zip:oam.tar.gz:Dockerfile:conf/pri.pem:conf/pub.pem -a=oam -be GOOS=linux -be GOARCH=amd64
)
if "%1" == "w" (
    bee pack -f=zip -exp=screenshot:logs:.vscode:.git:README.md:tests:build.bat:oam.zip:oam.tar.gz:Dockerfile:conf/pri.pem:conf/pub.pem -a=oam -be GOOS=windows -be GOARCH=amd64
)
