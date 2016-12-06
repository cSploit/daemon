package controllers

import (
	"github.com/cSploit/daemon/models"
	"github.com/cSploit/daemon/views"
	"github.com/gin-gonic/gin"
	"gopkg.in/oleiade/reflections.v1"
	"reflect"
	"strings"
)

func init() {
	j := models.Job{}

	fields, _ := reflections.Fields(j)

	for _, f := range fields {
		if tag, e := reflections.GetFieldTag(j, f, "gorm"); e == nil && strings.Contains(tag, "many2many") {
			v, _ := reflections.GetField(j, f) // []Host
			t := reflect.TypeOf(v).Elem()      // Host

			jobRelationships = append(jobRelationships, t)
		}
	}
}

// contains Types detected as affected entities
var jobRelationships = make([]reflect.Type, 0)

var JobController = Controller{
	EntityName: "job",
	Index:      jobIndex,
	Show:       jobShow,
}

func jobIndex(c *gin.Context) {
	var id uint64
	var found []models.Job

	db := models.GetDbInstance()

	// if a entity_id is available in the URL use it to restrict the searched jobs

	for _, eType := range jobRelationships {
		name := strings.ToLower(eType.Name())
		if getId(c, name, &id) == nil {
			entityPtr := reflect.New(eType).Interface()
			if err := reflections.SetField(entityPtr, "ID", uint(id)); err != nil {
				db.Error = err
				goto fetched
			}

			ass := db.Model(entityPtr).Association("Jobs")

			if ass.Error != nil {
				db.Error = ass.Error
				goto fetched
			}

			db.Error = ass.Find(&found).Error
			goto fetched
		}
	}

	db = db.Find(&found)

fetched:

	renderView(c, views.JobIndex, found, db)
}

func jobShow(c *gin.Context) {
	var id uint64
	var j models.Job

	if fetchId(c, "job", &id) != nil {
		return
	}

	db := models.GetDbInstance().Find(&j, id)

	renderView(c, views.JobShow, j, db)
}
