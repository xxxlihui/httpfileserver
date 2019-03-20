package main

import (
	"gopkg.in/urfave/cli.v2"
	"httpfileserver/httpServer"
	"os"
)

func main() {
	var port string
	var username string
	var password string
	var dir string
	app := &cli.App{
		Name:    "http文件服务器",
		Version: "1.0",
		Usage:   "一个简单的http文件服务器",
		Authors: []*cli.Author{{Name: "lhn", Email: "550124023@qq.com"}},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "port",
				Aliases:     []string{"p"},
				Usage:       "服务监听的端口",
				Value:       ":80",
				Destination: &port,
			},
			&cli.StringFlag{
				Name:        "username",
				Aliases:     []string{"n", "name"},
				Usage:       "用户名称,httpServer base认证",
				Value:       "sa",
				Destination: &username,
			},
			&cli.StringFlag{
				Name:        "password",
				Aliases:     []string{"pwd"},
				Usage:       "用户密码, httpServer base认证",
				Destination: &password,
				Value:       "123",
			},
			&cli.StringFlag{
				Name:        "directory",
				Aliases:     []string{"dir", "d"},
				Usage:       "文件的根目录",
				Destination: &dir,
				Value:       "~/httpserverdir/",
			},
		},
		Action: func(context *cli.Context) error {
			/*http.ListenAndServe(port, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Printf("url:%s",r.RequestURI)
			}))*/
			err := httpServer.Start(port, username, password, dir)
			return err
		},
	}
	app.Run(os.Args)
}
