package server

import "github.com/gin-gonic/gin"

func (s *ginServer) root(c *gin.Context) {
	url := s.buildUrl("/home")
	redirect(c, url)
}

func (s *ginServer) home(c *gin.Context) {

}