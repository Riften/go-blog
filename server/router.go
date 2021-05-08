package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func (s *ginServer) initRouter() {
	if s.cfg == nil {
		log.Fatal("no running config when init server router")
	}

	s.router = gin.Default()
	if s.cfg.RequestOutput() {
		s.router.Use(outPutInfo)
	}

	s.router.GET("", s.root)
	s.router.GET("/home", s.home)

	cmdGroup := s.router.Group("/cmd", assertLocalhost)
	{
		cmdGroup.POST("markdown/render", s.renderMd)
	}
}

func outPutInfo(c *gin.Context) {
	out := c.Request.Host + ": "+ c.Request.Method + " - " + c.Request.RequestURI// + "\n"
	log.Debug(out)
	c.Next()
}

func redirect(g *gin.Context, url string) {
	log.Debug("redirect to: ", url)
	g.Redirect(http.StatusMovedPermanently, url)
}

// assertLocalhost middleware is used before any request that only allow localhost.
// Such requests are from command line most times.
func assertLocalhost(c *gin.Context) {
	remote := strings.Split(c.Request.RemoteAddr, ":")[0]
	if remote == "127.0.0.1" {
		c.Next()
	} else {
		log.Warn("Receive an request not allowed for host other than localhost from ",
			remote, ": ", c.Request.Method, " - ", c.Request.RequestURI)
		c.String(http.StatusBadRequest, "This request is only allowed for localhost.")
		c.Abort()
	}
}
