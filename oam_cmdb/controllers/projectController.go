package controllers

import (
	"OAM/models"
	fn "OAM/util"
	"strings"

	"github.com/beego/beego/v2/core/logs"
	"github.com/emirpasic/gods/sets/hashset"
)

type ProjectController struct {
	AuthController
}

//查询项目明细,返回包括:项目基本信息,项目文档,项目账号,项目成员
func (c *ProjectController) Detail() {
	projId, err := c.GetInt("projId")
	if err != nil {
		c.JsonParamError("请选择项目")
	}
	proj := models.GetProjectById(projId)
	if proj == nil {
		c.JsonParamError("项目不存在或已删除")
	}

	curUser := c.getLoginUser()
	proj.Members = models.FindUserByProjId(proj.ProjId)
	if !proj.HasPermission(curUser) {
		c.JsonStatus(403, "无权查看该项目")
	}

	proj.Docs = models.FindDocsByProjId(projId)
	c.JsonOk(proj)
}

//查询项目列表
func (c *ProjectController) List() {
	curUser := c.getLoginUser()
	var projs []*models.Project
	if curUser.IsSupervisor() {
		projs = models.FindAllProjects()
	} else {
		projs = models.FindAuthorizedProjects(curUser.UserId)
	}
	c.Data["projs"] = projs
	//转为map给select用
	c.Data["projItems"] = fn.SliceToMap[*models.Project, int, string](projs, func(proj *models.Project) int {
		return proj.ProjId
	}, func(proj *models.Project) string {
		return proj.ProjName
	})
	c.TplName = "proj.html"
}

//保存项目
func (c *ProjectController) Save() {
	c.justPost()
	proj := models.Project{}
	if err := c.ParseForm(&proj); err != nil {
		c.JsonParamError(err.Error())
	}

	var isAdd = proj.ProjId == 0
	if errmsg := proj.Valid(isAdd); errmsg != "" {
		c.JsonParamError(errmsg)
	}
	curUser := c.getLoginUser()
	proj.UpdaterId = curUser.UserId
	var err error
	if isAdd {
		proj.CreatorId = curUser.UserId
		_, err = models.SaveProject(&proj)
	} else {
		_, err = models.UpdateProject(&proj)
	}
	if err != nil {
		c.JsonFailed("保存失败," + err.Error())
		return
	}
	logs.Info("用户%s保存项目:%s", curUser.UserName, proj.ProjName)
	c.JsonOk(proj)
}

//删除项目
func (c *ProjectController) DelProject() {
	projId, err := c.GetInt("projId")
	if err != nil || projId <= 0 {
		c.JsonParamError("请选择要删除的项目")
	}
	curUser := c.getLoginUser()
	//只有超管有删除项目权限
	if !curUser.IsSupervisor() {
		c.JsonForbidden()
	}

	ok := models.DeleteProject(projId)

	if !ok {
		c.JsonFailed("删除项目失败")
	}
	logs.Info("删除项目,删除人:%s,projId:%d", curUser.UserName, projId)
	c.JsonOk(ok)
}

//删除项目成员
func (c *ProjectController) DelMember() {
	projId, err := c.GetInt("projId")
	if err != nil || projId <= 0 {
		c.JsonParamError("请选择项目")
	}
	userId, err := c.GetInt("userId")
	if err != nil || userId <= 0 {
		c.JsonParamError("请选择要删除的成员")
	}
	ok := models.DeleteProjectMember(projId, userId)
	if ok {
		c.JsonOk(ok)
	} else {
		c.JsonFailed("删除成员失败")
	}
}

//增加项目成员
func (c *ProjectController) AddMember() {
	projId, err := c.GetInt("projId")
	if err != nil || projId <= 0 {
		c.JsonParamError("请选择项目")
	}
	userIdStr := c.GetString("userIds")
	if userIdStr == "" {
		c.JsonParamError("请选择成员")
		return
	}
	userIds := fn.ToIntSlice(strings.Split(userIdStr, ","))
	err = models.AddProjectMember(projId, userIds...)
	if err == nil {
		c.JsonOk(nil)
	} else {
		c.JsonFailed("新增成员失败")
	}
}

// 查询非项目成员,用于添加
func (c *ProjectController) SelectableMembers() {
	projId, err := c.GetInt("projId")
	if err != nil || projId <= 0 {
		c.JsonParamError("请选择项目")
	}
	members := models.FindNotProjMember(projId)
	c.JsonOk(members)
}

func (c *ProjectController) FindHostAndAppTree() {
	projId, err := c.GetInt("projId")
	if err != nil || projId <= 0 {
		c.JsonParamError("请选择项目")
	}
	//项目下主机
	hosts := models.FindHostsByProjId(projId)
	//项目下应用
	apps := models.FindAppInfoByProjId(projId)
	rootNode := models.HostTreeGridNode{}
	hostAppMapSize := 0
	var hostAppMap map[int]*hashset.Set
	if len(apps) > 0 {
		appIds := fn.MapToIntSlice(apps, func(src *models.AppInfo) int {
			return src.AppId
		})
		// 查出应用和主机的关联关系
		hostAppMap = models.FindRelHostIdByAppIds(appIds)
		hostAppMapSize = len(hostAppMap)
	}
	usedAppSet := hashset.New()
	if len(hosts) > 0 {
		for _, host := range hosts {
			node := rootNode.GetChild(host.HostId)
			if node == nil {
				node = rootNode.AddHostNode(*host)
			}
			if hostAppMapSize > 0 {
				appids, ok := hostAppMap[host.HostId]
				if ok {
					//给主机添加应用子节点
					for _, appid := range appids.Values() {
						for _, app := range apps {
							if app.AppId == appid {
								node.AddAppNode(*app)
								usedAppSet.Add(app.AppId)
								break
							}
						}
					}
				}
			}
		}
	}
	//未与主机关联的应用
	if len(apps) > 0 {
		for _, app := range apps {
			if usedAppSet.Size() > 0 {
				if !usedAppSet.Contains(app.AppId) {
					rootNode.AddAppNode(*app)
				}
			} else {
				rootNode.AddAppNode(*app)
			}
		}
	}
	c.JsonOk(rootNode.Children)

}
