package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/astaxie/beego/orm"
)

type Quiz struct {
	Id           int64   `orm:"auto"`
	Content      string  `orm:"size(128)"`
	Answer1      string  `orm:"size(128)"`
	Answer2      string  `orm:"size(128)"`
	Answer3      string  `orm:"size(128)"`
	WhichCorrect int64   `orm:"null`
	Lesson       *Lesson `orm:"rel(fk)"`
}

func init() {
	orm.RegisterModel(new(Quiz))
}

// AddQuiz insert a new Quiz into database and returns
// last inserted Id on success.
func AddQuiz(m *Quiz) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetQuizById retrieves Quiz by Id. Returns error if
// Id doesn't exist
func GetQuizById(id int64) (v *Quiz, err error) {
	o := orm.NewOrm()
	v = &Quiz{Id: id}
	if err = o.QueryTable(new(Quiz)).Filter("Id", id).RelatedSel().One(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetQuizByLessonId retrieves Quiz by Lesson Id. Returns error if
// Id doesn't exist
func GetQuizByLessonId(idL int64) (v *Quiz, err error) {
	o := orm.NewOrm()
	v = &Quiz{}

	v.Lesson, _ = GetLessonById(idL)

	if err = o.QueryTable(new(Quiz)).Filter("lesson_id", idL).RelatedSel().One(v); err == nil {
		return v, nil
	}

	return nil, err
}

// GetAllQuiz retrieves all Quiz matches certain condition. Returns empty list if
// no records exist
func GetAllQuiz(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Quiz))
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		qs = qs.Filter(k, v)
	}
	// order by:
	var sortFields []string
	if len(sortby) != 0 {
		if len(sortby) == len(order) {
			// 1) for each sort field, there is an associated order
			for i, v := range sortby {
				orderby := ""
				if order[i] == "desc" {
					orderby = "-" + v
				} else if order[i] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
			qs = qs.OrderBy(sortFields...)
		} else if len(sortby) != len(order) && len(order) == 1 {
			// 2) there is exactly one order, all the sorted fields will be sorted by this order
			for _, v := range sortby {
				orderby := ""
				if order[0] == "desc" {
					orderby = "-" + v
				} else if order[0] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return nil, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return nil, errors.New("Error: unused 'order' fields")
		}
	}

	var l []Quiz
	qs = qs.OrderBy(sortFields...).RelatedSel()
	if _, err = qs.Limit(limit, offset).All(&l, fields...); err == nil {
		if len(fields) == 0 {
			for _, v := range l {
				ml = append(ml, v)
			}
		} else {
			// trim unused fields
			for _, v := range l {
				m := make(map[string]interface{})
				val := reflect.ValueOf(v)
				for _, fname := range fields {
					m[fname] = val.FieldByName(fname).Interface()
				}
				ml = append(ml, m)
			}
		}
		return ml, nil
	}
	return nil, err
}

// UpdateQuiz updates Quiz by Id and returns error if
// the record to be updated doesn't exist
func UpdateQuizById(m *Quiz) (err error) {
	o := orm.NewOrm()
	v := Quiz{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteQuiz deletes Quiz by Id and returns error if
// the record to be deleted doesn't exist
func DeleteQuiz(id int64) (err error) {
	o := orm.NewOrm()
	v := Quiz{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Quiz{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
