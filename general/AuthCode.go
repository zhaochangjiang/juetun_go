package general

import (
	"github.com/astaxie/beego"
)

type AuthCode struct {
}

func (this *AuthCode) GetImage(controller *beego.Controller) {
	//w http.ResponseWriter, r *http.Request
	//var ccaptcha *gocaptcha.Captcha
	//	configFile := flag.String("c", "gocaptcha.conf", "the config file")
	//	captcha, err := gocaptcha.CreateCaptchaFromConfigFile(*configFile)
	//	w := this.Ctx.ResponseWriter
	//	r := *this.Ctx.Request
	//	key := r.FormValue("key")
	//	if len(key) >= 0 {
	//		cimg, err := ccaptcha.GetImage(key)
	//		log.Println("err", err)
	//		if nil == err {
	//			w.Header().Add("Content-Type", "image/png")
	//			png.Encode(w, cimg)
	//		} else {
	//			log.Printf("show image error:%s", err.Error())
	//			w.WriteHeader(500)
	//		}
	//	}

	//	log.Printf("[cmd:showimage][remote_addr:%s][key:%s]", r.RemoteAddr, key)
}
func (this *AuthCode) Validate(controller *beego.Controller) bool {
	return true
}
