package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/SevvyP/items/internal/db"
	"github.com/SevvyP/items/pkg"
	"github.com/gin-gonic/gin"
)

func (s *Server) HandleCreateItem(c *gin.Context) {
	decoder := json.NewDecoder(c.Request.Body)
	var item pkg.Item
    err := decoder.Decode(&item)
    if err != nil {
        log.Println(err)
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
    }
	if item.Name == "" || item.Description == "" {
		log.Println("create request missing name or description")
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}
	err = s.db.CreateItem(item, c)
	if err != nil {
		log.Println(err)
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, item)
}

func (s *Server) HandleGetItem(c *gin.Context) {
	if c.Param("name") == "" {
		log.Println("create request missing name or description")
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}
	item, err := s.db.GetItem(c.Param("name"), c)
	if err != nil {
		if err.Error() == db.ErrNotFound {
			log.Println(err)
			c.Writer.WriteHeader(http.StatusNotFound)
			return
		}
		log.Println(err)
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, item)
}

func (s *Server) HandleUpdateItem(c *gin.Context) {
	decoder := json.NewDecoder(c.Request.Body)
	var item pkg.Item
    err := decoder.Decode(&item)
    if err != nil {
        log.Println(err)
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
    }
	if item.Name == "" || item.Description == "" {
		log.Println("create request missing name or description")
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}
	err = s.db.UpdateItem(item, c)
	if err != nil {
		log.Println(err)
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, item)
}

func (s *Server) HandleDeleteItem(c *gin.Context) {
	if c.Param("name") == "" {
		log.Println("create request missing name or description")
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}
	item, err := s.db.DeleteItem(c.Param("name"), c)
	if err != nil {
		if err.Error() == db.ErrNotFound {
			log.Println(err)
			c.Writer.WriteHeader(http.StatusNotFound)
			return
		}
		log.Println(err)
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, item)
}