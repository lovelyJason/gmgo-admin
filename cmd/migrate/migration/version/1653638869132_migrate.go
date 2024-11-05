package version

import (
	"go-admin/cmd/migrate/migration/models"
	common "go-admin/common/models"
	"gorm.io/gorm"
	"runtime"
	"strconv"

	"go-admin/cmd/migrate/migration"
)

func init() {
	_, fileName, _, _ := runtime.Caller(0)
	migration.Migrate.SetVersion(migration.GetFilename(fileName), _1653638869132Test)
}

func _1653638869132Test(db *gorm.DB, version string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		var list []models.SysMenu
		//查询 SysMenu 表中的所有菜单，并按 parent_id 和 menu_id 排序，结果存储在 list 中
		err := tx.Model(&models.SysMenu{}).Order("parent_id,menu_id").Find(&list).Error
		if err != nil {
			return err
		}
		//遍历查询到的菜单列表，根据每个菜单的 ParentId 更新其 Paths 字段：
		//如果 ParentId 为 0，说明它是顶级菜单，Paths 被设置为 "/0/{MenuId}"。
		//如果 ParentId 不为 0，则查找其父菜单，获取父菜单的 Paths，并将当前菜单的 Paths 更新为 父菜单的路径 + "/" + 当前菜单ID。
		//更新后的 Paths 字段通过 UpdateByUserId 方法保存到数据库。
		for _, v := range list {
			if v.ParentId == 0 {
				v.Paths = "/0/" + strconv.Itoa(v.MenuId)
			} else {
				var e models.SysMenu
				err = tx.Model(&models.SysMenu{}).Where("menu_id=?", v.ParentId).First(&e).Error
				if err != nil {
					if err == gorm.ErrRecordNotFound {
						continue
					}
					return err
				}
				v.Paths = e.Paths + "/" + strconv.Itoa(v.MenuId)
			}
			err = tx.Model(&v).Update("paths", v.Paths).Error
			if err != nil {
				return err
			}
		}
		// 记录迁移版本
		return tx.Create(&common.Migration{
			Version: version,
		}).Error
	})
}
