package console

import (
	"github.com/gingerxman/eel"
	"github.com/gingerxman/eel/utils"
	"html/template"
	"bytes"
	"github.com/gingerxman/eel/router"
	"github.com/gingerxman/eel/config"
)

type Console struct {
	eel.RestResource
}

func (this *Console) Resource() string {
	return "console.console"
}

func (this *Console) GetParameters() map[string][]string {
	return map[string][]string{
		"GET":    []string{},
	}
}

func (this *Console) Get(ctx *eel.Context) {
	path, _ := utils.SearchFile("service_console.html", "./static")
	t, _ := template.ParseFiles(path)
	
	var bufferWriter bytes.Buffer
	resources := router.Resources()
	
	serviceName := config.ServiceConfig.String("appname")
	t.Execute(&bufferWriter, map[string]interface{}{
		"Resources": resources,
		"Name": serviceName,
	})
	
	ctx.Response.Content(bufferWriter.Bytes(), "text/html; charset=utf-8")
}
