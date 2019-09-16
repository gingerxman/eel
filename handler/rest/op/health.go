package op

import (
	"github.com/gingerxman/eel"
	"github.com/gingerxman/eel/config"
	"io/ioutil"
	"os"
	"time"
)

type Health struct {
	eel.RestResource
}

func (this *Health) Resource() string {
	return "op.health"
}

func (this *Health) SkipAuthCheck() bool {
	return true
}

func (this *Health) GetParameters() map[string][]string {
	return map[string][]string{
		"GET":    []string{},
	}
}

func (this *Health) Get(ctx *eel.Context) {
	content, err := ioutil.ReadFile("./image.version")
	if err != nil {
		//panic(err)
		content = []byte("no_image")
	}
	
	eelMode := os.Getenv("EEL_RUNMODE")
	k8sEnv := os.Getenv("_K8S_ENV")
	now := time.Now().Format("2006-01-02 15:04:05")
	serviceName := config.ServiceConfig.String("SERVICE_NAME")
	
	ctx.Response.JSON(eel.Map{
		"service":   serviceName,
		"is_online": true,
		"time":      now,
		"image":     string(content),
		"mode": eelMode,
		"k8s_env": k8sEnv,
	})
}
