package main

import (
	"OAM/conf"
	"OAM/controllers"
	"OAM/models"
	_ "OAM/routers"
	fn "OAM/util"
	"fmt"
	"syscall"

	"errors"
	"os"
	"path/filepath"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	jsoniter "github.com/json-iterator/go"
)

func main() {
	initLog()

	go initGlobalCfg()

	storePath := beego.AppConfig.DefaultString("store_path", "")
	if storePath == "" {
		storePath, err := fn.UserHomeDir()
		if err != nil {
			panic(err)
		}
		conf.GlobalCfg.DATA_STORE_PATH = filepath.Join(storePath, "/.oam")
	} else {
		conf.GlobalCfg.DATA_STORE_PATH = storePath
	}
	if !fn.FileIsExists(conf.GlobalCfg.DATA_STORE_PATH) {
		os.MkdirAll(conf.GlobalCfg.DATA_STORE_PATH, os.ModePerm)
	}
	beego.SetStaticPath(conf.STATIC_EXT_BASE_URL, conf.GlobalCfg.DATA_STORE_PATH)

	beego.AddFuncMap("json", func(in interface{}) (string, error) {
		return jsoniter.MarshalToString(in)
	})
	beego.AddFuncMap("hasRole", func(curUser models.LoginUser, roleCode string) bool {
		return curUser.HasRole(roleCode)
	})
	beego.AddFuncMap("hasPerm", func(curUser models.LoginUser, funCode string) bool {
		return curUser.HasPermByFunCode(funCode)
	})
	beego.ErrorController(&controllers.ErrorController{})
	if !conf.GlobalCfg.IS_INSTALLED {
		beego.InsertFilter("/*", beego.BeforeRouter, controllers.LoginAuthFilter)
	} else {
		//添加登录和权限控制过滤器
		beego.InsertFilter("/*", beego.BeforeRouter, controllers.LoginAuthFilter)
		beego.InsertFilter("/*", beego.BeforeExec, controllers.AccessFilter)
	}
	beego.Run()
}

//初始配置日志组件
func initLog() {
	isDebug, _ := beego.AppConfig.Bool("debug")
	if !isDebug {
		err := logs.SetLogger(logs.AdapterFile, `{"filename":"logs/oam.log","level":6,"maxdays":10,"daily":true}`)
		if err != nil {
			panic(err)
		}
		logs.EnableFuncCallDepth(true)
		logs.Async(1e3)
		fmt.Println("初始日志文件")
	}
}

// 初始一些通用配置
func initGlobalCfg() {
	var err error
	//当前程序工作目录
	workDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	// 缓存引擎类型
	conf.GlobalCfg.CACHE_TYPE = beego.AppConfig.DefaultString("cache_type", "memory")

	// 加载RSA密钥文件
	logs.Info("加载密钥文件")
	priKeyPath := filepath.Join(workDir, "conf", "pri.pem")
	pubKeyPath := filepath.Join(workDir, "conf", "pub.pem")
	//首次生成密钥
	if !fn.FileIsExists(pubKeyPath) && !fn.FileIsExists(priKeyPath) {
		err = fn.RSAGenerateKeyFile(1024, priKeyPath, pubKeyPath)
		if err != nil {
			logs.Error("生成密钥文件失败", err)
			panic(err)
		} else {
			logs.Info("生成新的密钥文件")
		}
	}
	conf.GlobalCfg.RSA_DEFAULT_PRIVATE_KEY, err = fn.LoadPrivateKeyFromFile(priKeyPath)
	if err != nil {
		panic(errors.New("加载私钥文件失败"))
	}

	conf.GlobalCfg.RSA_DEFAULT_PUBLIC_KEY, err = fn.LoadPublicKeyFromFile(pubKeyPath)
	if err != nil {
		panic(errors.New("加载公钥文件失败"))
	}

	//初始AES密钥
	conf.GlobalCfg.SYMMETRIC_KEY = models.InitSymmetrickey()

	//checkInstall(workDir)
}

func restart() {
	fmt.Println("重启程序...")

	// 获取当前可执行文件的路径
	exePath, err := os.Executable()
	if err != nil {
		fmt.Println("无法获取可执行文件路径:", err)
		return
	}

	// 使用 syscall.Exec 来替换当前进程
	err = syscall.Exec(exePath, os.Args, os.Environ())
	if err != nil {
		fmt.Println("重启失败:", err)
	}
}

func checkInstall(workDir string) {
	installFilePath := filepath.Join(workDir, conf.INSTALL_FILE_NAME)
	if !fn.FileIsExists(installFilePath) {
		logs.Error("检测到配置初始未完成...")
		conf.GlobalCfg.IS_INSTALLED = false
	} else {
		conf.GlobalCfg.IS_INSTALLED = true
	}
}
