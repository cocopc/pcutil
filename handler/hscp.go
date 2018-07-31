package handler

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/sdbaiguanghe/glog"
)

// Base 文件结构体
type Base struct {
	Path       string
}

// SetSavePath 设置文件存储路径
func (b *Base) SetSavePath(path string) {
	b.Path = path
}

// HSCPHandler 通过调用本地的scp命令拷贝文件
func (b *Base) HSCPHandler(w http.ResponseWriter, r *http.Request) {
	// parse 
	keys, ok := r.URL.Query()["dsthostdir"]
	if !ok || len(keys[0]) < 1 {
			glog.Infoln("Url Param 'dsthostdir' is missing")
			w.Write([]byte("Url Param 'dsthostdir' is missing\n"))
			return
	}
	// Query()["key"] will return an array of items, 
	// we only want the single item.
	dstHostDir := keys[0]


	r.ParseMultipartForm(32 << 20)
	//curl localhost:1031/upload -F file=@halt.sh
	file, handle, err := r.FormFile("file")
	if err != nil {
		glog.Errorln("FormFile错误", err)
		return
	}
	defer file.Close()

	// 格式化时间 2006-01-02 15:04:05
	dstDir := b.Path + time.Now().Format("20060102") + "/"
	os.MkdirAll(dstDir, 0777)
	dstFile := dstDir + handle.Filename
	df, err := os.OpenFile(dstFile, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		glog.Errorln("打开文件写入错误", err)
		return
	}

	defer df.Close()

	io.Copy(df, file)
	w.Write([]byte("Write file to tmpDir ok.\n"))
	execCommand := "scp " + dstFile + " " + dstHostDir
	w.Write([]byte("执行命令: " + execCommand + "\n"))
	result, err := execShell(execCommand)

	if err != nil {
		glog.Errorln("调用系统命令错误", result, err)
		w.Write([]byte("调用系统命令错误\n"))
		w.Write([]byte(result))
	} else {
		w.Write([]byte(result))
		w.Write([]byte("scp拷贝完毕"))
	}

}

//阻塞式的执行外部shell命令的函数,等待执行完毕并返回标准输出
func execShell(s string) (string, error) {

	//函数返回一个*Cmd，用于使用给出的参数执行name指定的程序
	cmd := exec.Command("/bin/bash", "-c", s)
	glog.Infoln("执行的命令:", cmd.Args)
	//读取io.Writer类型的cmd.Stdout，再通过bytes.Buffer(缓冲byte类型的缓冲器)将byte类型转化为string类型(out.String():这是bytes类型提供的接口)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	cmd.Stdout = &out
	//Run执行c包含的命令，并阻塞直到完成。  这里stdout被取出
	err := cmd.Run()
	if err != nil {
		glog.Infoln("执行命令错误:", stderr.String())
		return stderr.String(), err
	}
	return out.String(), nil
}

//错误处理函数
