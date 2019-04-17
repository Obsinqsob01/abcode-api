package controllers

import (
	"abcode/models"
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
)

//  QuizController operations for Quiz
type QuizController struct {
	beego.Controller
}

type LessonRelated struct {
	Lesson string `json:"lesson"`
}

// URLMapping ...
func (c *QuizController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
	c.URLFor("/lesson/:id", c.GetOneByLessonId)
}

// Post ...
// @Title Post
// @Description create Quiz
// @Param	body		body 	models.Quiz	true		"body for Quiz content"
// @Success 201 {int} models.Quiz
// @Failure 403 body is empty
// @router / [post]
func (c *QuizController) Post() {
	var v models.Quiz
	var lr LessonRelated

	json.Unmarshal(c.Ctx.Input.RequestBody, &v)
	json.Unmarshal(c.Ctx.Input.RequestBody, &lr)

	idLesson, _ := strconv.ParseInt(lr.Lesson, 10, 64)
	v.Lesson, _ = models.GetLessonById(idLesson)

	if _, err := models.AddQuiz(&v); err == nil {
		c.Ctx.Output.SetStatus(201)
		c.Data["json"] = v
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// GetOne ...
// @Title Get One
// @Description get Quiz by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Quiz
// @Failure 403 :id is empty
// @router /:id [get]
func (c *QuizController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 0, 64)
	v, err := models.GetQuizById(id)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = v
	}
	c.ServeJSON()
}

// GetOneByLessonId ...
// @Title Get One
// @Description get Tema by tema id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Lesson
// @Failure 403 :id is empty
// @router /lesson/:id [get]
func (c *QuizController) GetOneByLessonId() {
	idLesson := c.Ctx.Input.Param(":id")
	idL, _ := strconv.ParseInt(idLesson, 0, 64)

	v, err := models.GetQuizByLessonId(idL)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = v
	}
	c.ServeJSON()
}

// GetAll ...
// @Title Get All
// @Description get Quiz
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Quiz
// @Failure 403
// @router / [get]
func (c *QuizController) GetAll() {
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

	l, err := models.GetAllQuiz(query, fields, sortby, order, offset, limit)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = l
	}
	c.ServeJSON()
}

// Put ...
// @Title Put
// @Description update the Quiz
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Quiz	true		"body for Quiz content"
// @Success 200 {object} models.Quiz
// @Failure 403 :id is not int
// @router /:id [put]
func (c *QuizController) Put() {
	var lr LessonRelated
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 0, 64)
	v := models.Quiz{Id: id}

	json.Unmarshal(c.Ctx.Input.RequestBody, &v)
	json.Unmarshal(c.Ctx.Input.RequestBody, &lr)

	idLesson, _ := strconv.ParseInt(lr.Lesson, 10, 64)
	v.Lesson, _ = models.GetLessonById(idLesson)

	if err := models.UpdateQuizById(&v); err == nil {
		c.Data["json"] = "OK"
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// Delete ...
// @Title Delete
// @Description delete the Quiz
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *QuizController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 0, 64)
	if err := models.DeleteQuiz(id); err == nil {
		c.Data["json"] = "OK"
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}
