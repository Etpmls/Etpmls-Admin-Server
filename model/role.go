package model

import (
	"Etpmls-Admin-Server/core"
	"Etpmls-Admin-Server/database"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"time"
)

type Role struct {
	ID        uint `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
	Name string `json:"name"`
	Remark string `json:"remark"`
	Users []User `gorm:"many2many:role_users" json:"users"`
	Permissions []Permission `gorm:"many2many:role_permissions;" json:"permissions"`
}


// Create Role
// 创建角色
type ApiRoleCreate struct {
	ID        uint `json:"-"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-"`
	Name string `json:"name" binding:"required" validate:"max=30"`
	Remark string `json:"remark"`
	Permissions []Permission `gorm:"many2many:role_permissions;" json:"permissions"`
}
func (this *Role) RoleCreate(c *gin.Context, j ApiRoleCreate) (error) {
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		type Role ApiRoleCreate
		form := Role(j)

		// Insert Data
		result := tx.Create(&form)
		if result.Error != nil {
			core.LogError.Output(core.MessageWithLineNum(result.Error.Error()))
			return result.Error
		}

		// Role Create Event for module
		r, err := this.Role_InterfaceToRole(form)
		if err != nil {
			core.LogError.Output(core.MessageWithLineNum(err.Error()))
			return err
		}
		select {
		case core.Event.Event_RoleCreate <- core.EventObject{
			Context: c,
			Content: r,
		}:
		case <- time.After(time.Second * 3):
		}

		return nil
	})

	return err
}


// Modify role
// 修改角色
type ApiRoleEdit struct {
	ID        uint `json:"id" binding:"required" validate:"min=1"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-"`
	Name string `json:"name" binding:"required" validate:"max=30"`
	Remark string `json:"remark"`
	Permissions []Permission `gorm:"many2many:role_permissions;" json:"permissions"`
}
func (this *Role) RoleEdit(c *gin.Context, j ApiRoleEdit) (error) {
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		type Role ApiRoleEdit
		form := Role(j)
		result := tx.Save(&form)
		if result.Error != nil {
			core.LogError.Output(core.MessageWithLineNum(result.Error.Error()))
			return result.Error
		}

		// Role Edit Event for module
		r, err := this.Role_InterfaceToRole(form)
		if err != nil {
			core.LogError.Output(core.MessageWithLineNum(err.Error()))
			return err
		}
		select {
		case core.Event.Event_RoleEdit <- core.EventObject{
			Context: c,
			Content: r,
		}:
		case <- time.After(time.Second * 3):
		}

		return nil
	})

	return err
}


// Get all role
// 获取全部角色
type ApiRoleGetAll struct {
	ID        uint `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-"`
	Name string `json:"name" binding:"required" validate:"max=30"`
	Remark string `json:"remark"`
	Permissions []Permission `gorm:"many2many:role_permissions;" json:"permissions"`
}
func (this *Role) RoleGetAll(c *gin.Context) (interface{}, int64) {
	type Role ApiRoleGetAll
	var data []Role

	limit, offset := Common_GetPageByQuery(c)
	var count int64
	// Get the title of the search, if not get all the data
	// 获取搜索的标题，如果没有获取全部数据
	search := c.Query("search")

	database.DB.Model(&Role{}).Preload("Permissions").Where("name " + database.FUZZY_SEARCH + " ?", "%"+ search +"%").Count(&count).Limit(limit).Offset(offset).Find(&data)

	return data, count
}


// Delete roles (allow multiple deletions at the same time)
// 删除角色（允许同时删除多个）
type ApiRoleDelete struct {
	ID uint `json:"-"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	DeletedAt *time.Time `json:"-"`
	Roles []Role `json:"roles" binding:"required" validate:"min=1"`
}
func (this *Role) RoleDelete(c *gin.Context, ids []uint) (err error) {
	err = database.DB.Transaction(func(tx *gorm.DB) error {
		var r []Role
		tx.Where("id IN ?", ids).Find(&r)

		// 删除角色
		result := tx.Where("id IN ?", ids).Delete(&Role{})
		if result.Error != nil {
			core.LogError.Output(core.MessageWithLineNum(result.Error.Error()))
			return result.Error
		}

		// 删除关联
		err = tx.Model(&r).Association("Users").Clear()
		if err != nil {
			core.LogError.Output(core.MessageWithLineNum(err.Error()))
			return err
		}

		// 删除关联
		err = tx.Model(&r).Association("Permissions").Clear()
		if err != nil {
			core.LogError.Output(core.MessageWithLineNum(err.Error()))
			return err
		}

		// Role Delete Event for module
		select {
		case core.Event.Event_RoleDelete <- core.EventObject{
			Context: c,
			Content: r,
		}:
		case <- time.After(time.Second * 3):
		}

		return nil
	})

	return err
}


// interface conversion Role
// interface转换Role
func (this *Role) Role_InterfaceToRole(i interface{}) (Role, error) {
	var r Role
	us, err := json.Marshal(i)
	if err != nil {
		core.LogError.Output(core.MessageWithLineNum("Object to JSON failed!" + err.Error()))
		return Role{}, err
	}
	err = json.Unmarshal(us, &r)
	if err != nil {
		core.LogError.Output(core.MessageWithLineNum("JSON conversion object failed!" + err.Error()))
		return Role{}, err
	}
	return r, nil
}