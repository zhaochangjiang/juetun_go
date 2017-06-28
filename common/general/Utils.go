package general

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"io"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"strings"
)

type UtilsInterface interface {
	GetFileList(directory string) []string
}

type Utils struct {
}

//获得本机MAC地址
func (this *Utils) getMacAddress() *[]string {

	interfaces, err := net.Interfaces()
	if err != nil {
		panic("Error : " + err.Error())
	}
	mac := make([]string, 0)
	for _, inter := range interfaces {
		if len(inter.HardwareAddr) > 0 {
			mac = append(mac, inter.HardwareAddr.String()) //获取本机MAC地址
		}
	}
	return &mac
}

//切片的头追加数据
func (this *Utils) Slice_unshift(oSlice []interface{}, content interface{}) *[]interface{} {
	slice := []interface{}{content}
	slice = append(slice, oSlice...)
	return &slice
}

//生成32位md5字串
func (this *Utils) GetMd5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

//生成Guid字串
func (this *Utils) GetGuid() string {
	b := make([]byte, 48)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return this.GetMd5String(base64.URLEncoding.EncodeToString(b))
}

//获得当前的路径
func (this *Utils) getCurrentPath() string {
	s, err := exec.LookPath(os.Args[0])
	this.checkErr(err)
	i := strings.LastIndex(s, "\\")
	path := string(s[0 : i+1])
	return path
}

//获得当前目录下的文件列表
func (this *Utils) GetFileList(directory string) []string {

	dir_list, e := ioutil.ReadDir(directory)

	//如果有错误，则抛出
	if e != nil {
		this.checkErr(e)
	}
	data := make([]string, 0, 1)
	for _, v := range dir_list {
		name := v.Name()
		data = append(data, name)
	}
	return data
}

//错误处理
func (this *Utils) checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
