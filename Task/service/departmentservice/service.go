package service

import (
	"Task/common"
	departmentrepository "Task/repository/department"
	employeerepository "Task/repository/employee"
	"encoding/json"
	"strconv"

	"github.com/go-chassis/openlog"
)

type Service struct {
	DepartmentRepo departmentrepository.DepartmentRepo
	Repo           employeerepository.EmployeeRepo
}

func (h *Service) CreateDepartment(input common.CreateInput) common.Response {
	departmentname := input.Create["department name"].(string)
	id := input.Create["id"].(float64)
	code, err := h.DepartmentRepo.IsDepartmentExists(departmentname, 0)
	if err != nil {
		openlog.Error(err.Error())
		return common.Response{Msg: err.Error(), Data: nil, Status: code}
	}
	code1, err1 := h.DepartmentRepo.IsDepartmentExists("", id)
	if err1 != nil {
		openlog.Error(err.Error())
		return common.Response{Msg: err.Error(), Data: nil, Status: code1}
	}

	res, code, err := h.DepartmentRepo.DepartmentInsert(input.Create)
	if err != nil {
		openlog.Error(err.Error())
		return common.Response{Msg: err.Error(), Data: nil, Status: code}
	}
	return common.Response{Msg: "Department inserted successfully", Data: res, Status: 201}
}

func (h *Service) GetDepartmentbyid(input common.FetchDepartmentInput) common.Response {
	res, code, err := h.DepartmentRepo.DepartmentFind(input.ID)
	if err != nil {
		openlog.Error(err.Error())
		return common.Response{Msg: err.Error(), Data: nil, Status: code}
	}
	return common.Response{Msg: "Department Fetched successfully", Data: res, Status: 200}
}

//CheckAndInitializePageDetails will initilize page,and it's size
func CheckAndInitializePageDetails(page, size string) (int64, int64) {
	pageNo, err := strconv.ParseInt(page, 10, 64)
	if err != nil || pageNo < 1 {
		pageNo = -1
	}

	var pageSize int64
	pageSize, err = strconv.ParseInt(size, 10, 64)
	if err != nil || pageSize < 1 {
		pageSize = 5
	}
	return pageNo, pageSize
}

// FetchAllDatamodelsByPagenation function will helps to get the field by considering page number, size and filters.
func (h *Service) FetchAllDepartment(input common.FetchAllInput) common.Response {
	var filter = make(map[string]interface{})
	if input.Filters != "" {
		bytes := []byte(input.Filters)
		json.Unmarshal(bytes, &filter)
	}
	var sortorder = make(map[string]interface{})
	if input.Sort != "" {
		bytes := []byte(input.Sort)
		_ = json.Unmarshal(bytes, &sortorder)
	}
	var res = make([]map[string]interface{}, 0)
	var err error
	var code int
	pageNo, pageSize := CheckAndInitializePageDetails(input.Page, input.Size)
	if pageNo < 0 {
		res, code, err = h.DepartmentRepo.FindByFilters(filter, sortorder)
		if err != nil {
			openlog.Error(err.Error())
			return common.Response{Msg: err.Error(), Data: nil, Status: code}
		}
	} else {
		res, code, err = h.DepartmentRepo.FindByFiltersAndPagenation(pageNo, pageSize, filter, sortorder)
		if err != nil {
			openlog.Error(err.Error())
			return common.Response{Msg: err.Error(), Data: nil, Status: code}

		}
	}
	return common.Response{Msg: "All Users Fetched successfully", Data: res, Status: 200}
}

func (h *Service) DeleteDepartment(input common.DeleteInput) common.Response {
	openlog.Info("Got a request to Delete Department")

	Employeefetch, code, err := h.DepartmentRepo.DepartmentFind(input.ID)
	if err != nil {
		openlog.Error(err.Error())
		return common.Response{Msg: err.Error(), Data: nil, Status: code}
	}
	result, _, err := h.Repo.FindDepartment(Employeefetch["id"].(float64))
	if len(result) > 0 {
		return common.Response{Msg: "Employee exist with department id", Data: result, Status: 400}

	}
	if err != nil {
		openlog.Error(err.Error())
		return common.Response{Msg: err.Error(), Data: nil, Status: code}

	}
	res, code, err := h.DepartmentRepo.DeleteDepartment(input.ID)
	if err != nil {
		openlog.Error(err.Error())
		return common.Response{Msg: err.Error(), Data: nil, Status: code}

	}
	return common.Response{Msg: "Department deleted successfully", Data: res, Status: 200}
}

func (h *Service) UpdateDepartment(input common.UpdateInput) common.Response {
	fetchresult, code, err := h.DepartmentRepo.DepartmentFind(input.ID)
	var flag bool
	if err != nil {
		openlog.Error(err.Error())
		return common.Response{Msg: err.Error(), Data: nil, Status: code}
	}
	name, ok := input.Update["department name"].(string)
	id, idok := input.Update["id"].(float64)
	if idok {
		if fetchresult["id"] == id {
		} else {
			code, err = h.DepartmentRepo.IsDepartmentExists("", id)
			flag = true
		}
	}
	if ok {
		if fetchresult["department name"].(string) == name {
		} else {
			code, err = h.DepartmentRepo.IsDepartmentExists(name, 0)
			flag = true
		}
		if flag {
			if err != nil {
				openlog.Error(err.Error())
				return common.Response{Msg: err.Error(), Data: nil, Status: code}
			}
		}
	}
	res, code, err := h.DepartmentRepo.FindAndUpdateDepartment(input.ID, input.Update)
	if err != nil {
		openlog.Error(err.Error())
		return common.Response{Msg: err.Error(), Data: nil, Status: code}
	}
	return common.Response{Msg: "Department Updated successfully", Data: res, Status: 201}
}
