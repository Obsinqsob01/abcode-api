package controllers

import (
	"abcode/models"
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
)

//  TemaController operations for Tema
type TemaController struct {
	beego.Controller
}

// CourseRelated
type CourseRelated struct {
	Course string `json:"course"`
}

// URLMapping ...
func (c *TemaController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
	c.URLFor("/:courseId/:id", c.GetOneByCourseId)
	c.URLFor("/all/:courseId", c.GetAllByCourseId)
}

// Post ...
// @Title Post
// @Description create Tema
// @Param	body		body 	models.Tema	true		"body for Tema content"
// @Success 201 {int} models.Tema
// @Failure 403 body is empty
// @router / [post]
func (c *TemaController) Post() {
	var v models.Tema
	var cr CourseRelated

	json.Unmarshal(c.Ctx.Input.RequestBody, &cr)
	json.Unmarshal(c.Ctx.Input.RequestBody, &v)

	idCourse, _ := strconv.ParseInt(cr.Course, 10, 64)
	v.Course, _ = models.GetCourseById(idCourse)

	if _, err := models.AddTema(&v); err == nil {
		c.Ctx.Output.SetStatus(201)
		c.Data["json"] = v
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// GetOne ...
// @Title Get One
// @Description get Tema by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Tema
// @Failure 403 :id is empty
// @router /:id [get]
func (c *TemaController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 0, 64)
	v, err := models.GetTemaById(id)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = v
	}
	c.ServeJSON()
}

// GetOneByCourseId ...
// @Title Get One
// @Description get Tema by course id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Tema
// @Failure 403 :id is empty
// @router /:courseId/:id [get]
func (c *TemaController) GetOneByCourseId() {
	idTema := c.Ctx.Input.Param(":id")
	idT, _ := strconv.ParseInt(idTema, 0, 64)

	idCourse := c.Ctx.Input.Param(":courseId")
	idC, _ := strconv.ParseInt(idCourse, 0, 64)

	v, err := models.GetTemaByCourseId(idT, idC)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = v
	}
	c.ServeJSON()
}

// GetAll ...
// @Title Get All
// @Description get Tema
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Tema
// @Failure 403
// @router / [get]
func (c *TemaController) GetAll() {
	var fields []string
	var sortby []string
	var order []string
	var query = make(map[string]string)
	var limit int64 = 10
	var offset int64

	// fields: col1,col2,entity.col3
	if v := c.GetString("fields"); v != "" {
		fields = strings.Split(v, ",")
	}
	// limit: 10 (default is 10)
	if v, err := c.GetInt64("limit"); err == nil {
		limit = v
	}
	// offset: 0 (default is 0)
	if v, err := c.GetInt64("offset"); err == nil {
		offset = v
	}
	// sortby: col1,col2
	if v := c.GetString("sortby"); v != "" {
		sortby = strings.Split(v, ",")
	}
	// order: desc,asc
	if v := c.GetString("order"); v != "" {
		order = strings.Split(v, ",")
	}
	// query: k:v,k:v
	if v := c.GetString("query"); v != "" {
		for _, cond := range strings.Split(v, ",") {
			kv := strings.SplitN(cond, ":", 2)
			if len(kv) != 2 {
				c.Data["json"] = errors.New("Error: invalid query key/value pair")
				c.ServeJSON()
				return
			}
			k, v := kv[0], kv[1]
			query[k] = v
		}
	}

	l, err := models.GetAllTema(query, fields, sortby, order, offset, limit)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = l
	}
	c.ServeJSON()
}

// GetAllByCourseId ...
// @Title Get All By Course Id
// @Description get Tema
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Tema
// @Failure 403
// @router /all/:courseId [get]
func (c *TemaController) GetAllByCourseId() {
	var fields []string
	var sortby []string
	var order []string
	var query = make(map[string]string)
	var limit int64 = 10
	var offset int64

	// fields: col1,col2,entity.col3
	if v := c.GetString("fields"); v != "" {
		fields = strings.Split(v, ",")
	}
	// limit: 10 (default is 10)
	if v, err := c.GetInt64("limit"); err == nil {
		limit = v
	}
	// offset: 0 (default is 0)
	if v, err := c.GetInt64("offset"); err == nil {
		offset = v
	}
	// sortby: col1,col2
	if v := c.GetString("sortby"); v != "" {
		sortby = strings.Split(v, ",")
	}
	// order: desc,asc
	if v := c.GetString("order"); v != "" {
		order = strings.Split(v, ",")
	}
	// query: k:v,k:v
	if v := c.GetString("query"); v != "" {
		for _, cond := range strings.Split(v, ",") {
			kv := strings.SplitN(cond, ":", 2)
			if len(kv) != 2 {
				c.Data["json"] = errors.New("Error: invalid query key/value pair")
				c.ServeJSON()
				return
			}
			k, v := kv[0], kv[1]
			query[k] = v
		}
	}

	idStr := c.Ctx.Input.Param(":courseId")
	id, _ := strconv.ParseInt(idStr, 0, 64)

	l, err := models.GetAllTemaByCourseId(query, fields, sortby, order, offset, limit, id)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = l
	}
	c.ServeJSON()
}

// Put ...
// @Title Put
// @Description update the Tema
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Tema	true		"body for Tema content"
// @Success 200 {object} models.Tema
// @Failure 403 :id is not int
// @router /:id [put]
func (c *TemaController) Put() {
	var cr CourseRelated
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 0, 64)
	v := models.Tema{Id: id}

	json.Unmarshal(c.Ctx.Input.RequestBody, &cr)
	json.Unmarshal(c.Ctx.Input.RequestBody, &v)

	idCourse, _ := strconv.ParseInt(cr.Course, 10, 64)
	v.Course, _ = models.GetCourseById(idCourse)

	if err := models.UpdateTemaById(&v); err == nil {
		c.Data["json"] = "OK"
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// Delete ...
// @Title Delete
// @Description delete the Tema
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *TemaController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 0, 64)
	if err := models.DeleteTema(id); err == nil {
		c.Data["json"] = "OK"
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}
