package httpServer

import (
	"archive/zip"
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"time"
)

var localdir string

func Start(port, username, password, dir string) error {
	UserName = username
	Password = password
	localdir = dir
	if inf, err := os.Stat(localdir); err != nil && os.IsNotExist(err) {
		err := os.MkdirAll(localdir, os.ModePerm)
		if err != nil {
			fmt.Printf("创建文件[%s]错误:%s", localdir, err.Error())
		}
	} else {
		fmt.Printf("%+v", inf)
	}
	app := gin.Default()
	/*app.Group("", func(context *gin.Context) {
		switch context.Request.Method {
		case http.MethodGet:
			get(context)
		case http.MethodPost:
			fallthrough
		case http.MethodPut:
			put(context)
		case http.MethodDelete:
			delete(context)

		}
	})*/
	app.NoRoute(func(context *gin.Context) {
		switch context.Request.Method {
		case http.MethodGet:
			get(context)
		case http.MethodPost:
			post(context)
		case http.MethodPut:
			put(context)
		case http.MethodDelete:
			delete(context)
		}
		if !context.Writer.Written() {
			context.Writer.Write([]byte{0})
		}
		//if !context.IsAborted() {
		//context.AbortWithStatus(http.StatusOK)
		//}

	})
	/*app.Use(author, fileSecurity)
	app.GET("", get)
	app.DELETE("", delete)
	app.PUT("", put)
	app.POST("", post)*/
	err := app.Run(port)
	if err != nil {
		return err
	}

	return nil
}

var cstZone = time.FixedZone("CST", 8*3600)

func formatSize(size int64) string {
	switch {
	case size > 1024:
		//K
		return fmt.Sprintf("%.2fK", float64(size/1024))
	case size > 1024*1024:
		//M
		return fmt.Sprintf("%.2fM", float64(size/1024*1024))
	case size > 1024*1024*1024:
		return fmt.Sprintf("%.2fG", float64(size/1024*1024*1024))
	//G
	default:
		return fmt.Sprintf("%dB", size)
	}
}

func get(c *gin.Context) {
	rquri, err := url.PathUnescape(c.Request.RequestURI)
	if err != nil {
		rquri = c.Request.RequestURI
	}
	ph := path.Join(localdir, rquri)
	info, err := os.Stat(ph)
	if err != nil {
		c.String(http.StatusOK, "目录不存在")
		return
	}
	/*		tmp := `<html>
			<head><title>Index of /</title></head>
			<body>
			<h1>Index of /</h1><hr><pre><a href="../">../</a>
			<a href="%E4%BE%9B%E9%9C%80/">供需/</a>                                                23-Jan-2019 10:37       -
			<a href="%E5%95%86%E5%AE%B6%E8%BF%9B%E9%A9%BB/">商家进驻/</a>                                              13-Mar-2019 17:02       -
			<a href="%E5%A4%A7%E6%95%B0%E6%8D%AE/">大数据/</a>                                               22-Jan-2019 10:57       -
			<a href="%E9%9C%80%E6%B1%82/">需求/</a>                                                02-Jan-2019 14:13       -
			<a href="ind">ind</a>                                                29-Nov-2018 10:58       0
			</pre><hr></body>
			</html>`
	*/
	if info.IsDir() {
		s := fmt.Sprintf(`<html>
<head><title>Index of /</title></head>
<body>
<h1>Index of %s</h1><hr><pre><form style="display:flex" id="upload-form" action="./" method="post" enctype="multipart/form-data" >
<a style="margin-left:20px;margin-right:50px;" href="../">../</a>
<label style="margin-left:20px;margin-right:20px;display:flex;" ><input type="checkbox" value="true" name="zip"/>自动解压zip</label>
<input type="file" id="upload" name="file" />
<input type="submit" value="上传" />

</form><br/>`, rquri)
		if c.Request.RequestURI == "/" {
			s = `<html>
<head><title>Index of /</title></head>
<body>
<h1>Index of /</h1><hr><pre><form style="display:flex" id="upload-form" action="./" method="post" enctype="multipart/form-data" >
<label style="margin-left:20px;margin-right:20px;display:flex;"><input value="true"  type="checkbox" name="zip"/>自动解压zip</label>
　　　<input type="file" id="upload" name="file" />
　　　<input type="submit" value="上传" />
</form>`

		}
		fs, err := ioutil.ReadDir(ph)
		if err != nil {
			c.String(http.StatusInternalServerError, "服务器错误")
			return
		}
		for _, f := range fs {
			if f.IsDir() {
				s += fmt.Sprintf(`<a href="%s/">%s/</a>                                                %s       -<br/>`, f.Name(), f.Name(), f.ModTime().In(cstZone).String())
			} else {
				s += fmt.Sprintf(`<a href="%s">%s</a>                                                %s       %s<br/>`, f.Name(), f.Name(), f.ModTime().In(cstZone).String(), formatSize(f.Size()))
			}
		}
		s += `</pre><hr></body>
</html>`
		c.Data(http.StatusOK, "text/html;charset=utf-8", []byte(s))
	} else {
		//返回内容
		f, err := os.Open(ph)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, &ErrorMessage{Message: "读取上传文件失败"})
			return
		}
		defer f.Close()
		fi, _ := f.Stat()
		http.ServeContent(c.Writer, c.Request, fi.Name(), fi.ModTime(), f)
		//c.File(ph)
	}

}
func delete(c *gin.Context) {
	ph := c.Request.RequestURI
	ph = path.Join(localdir, ph)
	fn, err := os.Stat(ph)
	if err != nil {
		c.Abort()
		return
	}
	if fn.IsDir() {
		os.RemoveAll(ph)
	} else {
		os.Remove(ph)
	}
}
func put(c *gin.Context) {
	rquri, err := url.PathUnescape(c.Request.RequestURI)
	if err != nil {
		rquri = c.Request.RequestURI
	}
	ph := path.Join(localdir, rquri)
	f, err := c.FormFile("file")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, &ErrorMessage{Message: "读取上传文件失败"})
		return
	}
	fs, err := f.Open()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, &ErrorMessage{Message: "读取上传文件失败"})
		return
	}
	defer fs.Close()
	ff, err := os.OpenFile(ph, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.ModePerm)
	if err != nil {
		fmt.Printf("创建文件错误:%s", err.Error())
		if os.IsNotExist(err) {
			os.MkdirAll(path.Dir(ph), os.ModePerm)
			ff, err = os.OpenFile(ph, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.ModePerm)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, &ErrorMessage{Message: "服务器创传文件失败"})
				return
			}
			defer ff.Close()
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, &ErrorMessage{Message: "服务器创传文件失败"})
			return
		}
	} else {
		defer ff.Close()
	}
	if _, err = io.Copy(ff, fs); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, &ErrorMessage{Message: "服务器保存文件失败"})
		return
	}
}
func post(c *gin.Context) {
	f, err := c.FormFile("file")
	if err != nil {
		c.Data(http.StatusInternalServerError, "text/html;charset=utf-8", []byte(`请选择文件<br/>`))
		return
	}

	z := c.PostForm("zip")
	if z == "true" {
		/*ef, err := ioutil.TempFile("", f.Filename)
		if err != nil {
			c.Data(http.StatusInternalServerError, "text/html;charset=utf-8", []byte(`保存文件失败`))
			return
		}
		defer os.Remove(ef.Name())*/
		fs, err := f.Open()
		if err != nil {
			c.Data(http.StatusInternalServerError, "text/html;charset=utf-8", []byte(`保存文件失败`))
			return
		}
		defer fs.Close()
		/*_, err = io.Copy(ef, fs)
		if err != nil {
			c.Data(http.StatusInternalServerError, "text/html;charset=utf-8", []byte(`保存文件失败`))
			return
		}*/
		//解压文件
		rd, err := zip.NewReader(fs, f.Size)
		if err != nil {
			c.Data(http.StatusInternalServerError, "text/html;charset=utf-8", []byte(`保存文件失败`))
			return
		}
		for _, e := range rd.File {
			name := e.Name
			if e.FileHeader.NonUTF8 {
				name, err = DecodeGBK(name)
				if err != nil {
					c.Data(http.StatusInternalServerError, "text/html;charset=utf-8", []byte(`解压文件失败，未知编码`))
					return
				}
			}
			if e.FileInfo().IsDir() {
				os.MkdirAll(path.Join(localdir, c.Request.RequestURI, name), os.ModePerm)
			} else {
				ff, err := os.OpenFile(path.Join(localdir, c.Request.RequestURI, name), os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.ModePerm)
				if err != nil {
					c.Data(http.StatusInternalServerError, "text/html;charset=utf-8", []byte(`保存文件失败`))
					return
				}
				r, err := e.Open()
				if err != nil {
					c.Data(http.StatusInternalServerError, "text/html;charset=utf-8", []byte(`保存文件失败`))
					ff.Close()
					return
				}
				_, err = io.Copy(ff, r)
				if err != nil {
					c.Data(http.StatusInternalServerError, "text/html;charset=utf-8", []byte(`保存文件失败`))
					r.Close()
					ff.Close()
					return
				}
				ff.Close()
				r.Close()
			}
		}
	} else {
		fn := filepath.Join(localdir, c.Request.RequestURI, f.Filename)
		ff, err := os.OpenFile(fn, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.ModePerm)
		if err != nil {
			c.Data(http.StatusInternalServerError, "text/html;charset=utf-8", []byte(`保存文件失败`))
			return
		}
		fs, err := f.Open()
		if err != nil {
			c.Data(http.StatusInternalServerError, "text/html;charset=utf-8", []byte(`保存文件失败`))
			return
		}
		defer fs.Close()
		_, err = io.Copy(ff, fs)
		if err != nil {
			c.Data(http.StatusInternalServerError, "text/html;charset=utf-8", []byte(`保存文件失败`))
			return
		}
	}
	c.Redirect(http.StatusMovedPermanently, c.Request.RequestURI)
}

func DecodeGBK(s string) (string, error) {
	I := bytes.NewReader([]byte(s))
	O := transform.NewReader(I, simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(O)
	if e != nil {
		return "", e
	}
	return string(d), nil
}

type ErrorMessage struct {
	Message string `json:"message"`
}
