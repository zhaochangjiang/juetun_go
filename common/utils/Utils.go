package utils

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"io"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"reflect"
	"strings"
)

/**
* 将一个结构体数据转换为Map数据
* @author karl.zhao<zhaocj2009@hotmail.com>
* @date 2017/08/15
* @return  map[string]interface{}
*
 */
func Struct2Map(obj interface{}) map[string]interface{} {

	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Name] = v.Field(i).Interface()
	}
	return data
}

/**
* 判断一个值是否在一个数组或者切片中
* @author karl.zhao<zhaocj2009@hotmail.com>
* @date 2017/08/15
* @param p interface{} 数值
* @param arr []interface{} 数组或切片值
* @return boolean
*
 */
func InArrayOrSlice(p interface{}, arr []interface{}) bool {
	for _, v := range arr {
		if v == p {
			return true
		}
	}
	return false

}

/**
* 字符串截取函数
* @author karl.zhao<zhaocj2009@hotmail.com>
* @date 2017/08/15
* @param str string 被截取的字符串
* @param length int 截取字符串的长度
* @return string
 */
func Substr(str string, start, length int) string {
	rs := []rune(str)
	rl := len(rs)
	end := 0
	//如果开始为负数，则从字符串尾部开始算
	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length

	if start > end {
		start, end = end, start
	}

	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}
	return string(rs[start:end])
}

/**
* 将一个未知数据类型转换为字符串格式，如果转换失败，则输出错误信息
* @author karl.zhao<zhaocj2009@hotmail.com>
* @date 2017/08/15
* @param p interface{} 转换的数据
* @return string,error
 */
func ConvertInterfaceToString(p interface{}) (string, error) {
	var c string
	var e error
	switch p.(type) {
	case string:
		c = p.(string)
	default:
		e = errors.New("you send params type must be string")
	}
	return c, e
}

/**
* 判断KEY 是否存在于Map中
* @author karl.zhao<zhaocj2009@hotmail.com>
* @date 2017/08/15
* @param key string 数值
* @param map[string]interface{} map内容
* @return bool
 */
func Isset(key string, mapContent map[string]interface{}) bool {
	if _, ok := mapContent[key]; ok {
		return true
	}
	return false
}

/**
* 获得本机MAC地址
* @author karl.zhao<zhaocj2009@hotmail.com>
* @date 2017/08/15
* @return *[]string
 */
func getMacAddress() *[]string {

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

/**
* strings.TrimRight切割 参数为 permitcontroller,permit时异常的自编方法
* @author karl.zhao<zhaocj2009@hotmail.com>
* @date 2017/09/05
*
 */
func TrimRight(s string, cutset string) string {
	if s == "" || cutset == "" {
		return s
	}
	maxi := strings.LastIndex(s, cutset)
	rs := []rune(s)
	return string(rs[0:maxi])
}

//切片的头追加数据
func SliceUnshift(oSlice []interface{}, content interface{}) *[]interface{} {
	slice := []interface{}{content}
	slice = append(slice, oSlice...)
	return &slice
}

//切片的头追加数据
func SliceUnshiftString(oSlice []string, content string) *[]string {
	slice := []string{content}
	slice = append(slice, oSlice...)
	return &slice
}

//切片的头追加数据
func SliceUnshiftInt(oSlice []int, content int) *[]int {
	slice := []int{content}
	slice = append(slice, oSlice...)
	return &slice
}

//生成32位md5字串
func GetMd5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

//生成Guid字串
func GetGuid() string {
	b := make([]byte, 48)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return GetMd5String(base64.URLEncoding.EncodeToString(b))
}

//获得当前的路径
func getCurrentPath() string {
	s, err := exec.LookPath(os.Args[0])
	checkErr(err)
	i := strings.LastIndex(s, "\\")
	path := string(s[0 : i+1])
	return path
}

//获得当前目录下的文件列表
func GetFileList(directory string) []string {

	dir_list, e := ioutil.ReadDir(directory)

	//如果有错误，则抛出
	if e != nil {
		checkErr(e)
	}
	data := make([]string, 0, 1)
	for _, v := range dir_list {
		name := v.Name()
		data = append(data, name)
	}
	return data
}

//错误处理
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
