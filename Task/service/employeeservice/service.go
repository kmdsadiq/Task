package employeeservice

import (
	"Task/common"
	departmentrepository "Task/repository/department"
	repository "Task/repository/employee"
	"encoding/json"
	"strconv"

	"github.com/go-chassis/openlog"
)

type EmployeeService struct {
	Repo           repository.EmployeeRepo
	DepartmentRepo departmentrepository.DepartmentRepo
}

func (e *EmployeeService) CreateEmployee(input common.CreateInput) common.Response {
	mail := input.Create["email"].(string)
	id := input.Create["id"].(int)
	departmentid := input.Create[common.Departmentid].(float64)
	fetchresult, code, err := e.DepartmentRepo.DepartmentExists(departmentid)
	if err != nil {
		openlog.Error(err.Error())
		return common.Response{Msg: err.Error(), Data: nil, Status: code}
	}
	if len(fetchresult) == 0 {
		return common.Response{Msg: "No department found with given input"}
	}
	code, err = e.Repo.IsEmployeeExists(mail, 0)
	if err != nil {
		openlog.Error(err.Error())
		return common.Response{Msg: err.Error(), Data: nil, Status: code}
	}
	code1, err1 := e.Repo.IsEmployeeExists("", id)
	if err1 != nil {
		openlog.Error(err.Error())
		return common.Response{Msg: err.Error(), Data: nil, Status: code1}
	}
	res, code, err := e.Repo.EmployeeInsert(input.Create)
	if err != nil {
		openlog.Error(err.Error())
		return common.Response{Msg: err.Error(), Data: nil, Status: code}
	}
	return common.Response{Msg: "Employee inserted successfully", Data: res, Status: 201}
}

func (e *EmployeeService) GetEmployeebyid(input common.FetchDepartmentInput) common.Response {
	res, code, err := e.Repo.FindEmployee(input.ID)
	if err != nil {
		openlog.Error(err.Error())
		return common.Response{Msg: err.Error(), Data: nil, Status: code}
	}
	return common.Response{Msg: "Employee Fetched successfully", Data: res, Status: 200}
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
func (e *EmployeeService) FetchAllEmployees(input common.FetchAllInput) common.Response {
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
		res, code, err = e.Repo.FindByFilters(filter, sortorder)
		if err != nil {
			openlog.Error(err.Error())
			return common.Response{Msg: err.Error(), Data: nil, Status: code}
		}
	} else {
		res, code, err = e.Repo.FindByFiltersAndPagenation(pageNo, pageSize, filter, sortorder)
		if err != nil {
			openlog.Error(err.Error())
			return common.Response{Msg: err.Error(), Data: nil, Status: code}

		}
	}
	return common.Response{Msg: "All Employees Fetched successfully", Data: res, Status: 200}
}

func (e *EmployeeService) DeleteEmployee(input common.DeleteInput) common.Response {
	openlog.Info("Got a request to Delete User")
	res, code, err := e.Repo.DeleteEmployee(input.ID)
	if err != nil {
		openlog.Error(err.Error())
		return common.Response{Msg: err.Error(), Data: nil, Status: code}

	}
	return common.Response{Msg: "Employee record deleted successfully", Data: res, Status: 200}
}

func (e *EmployeeService) UpdateEmployee(input common.UpdateInput) common.Response {
	fetchresult, code, err := e.Repo.FindEmployee(input.ID)
	var flag bool
	if err != nil {
		openlog.Error(err.Error())
		return common.Response{Msg: err.Error(), Data: nil, Status: code}
	}
	mail, ok := input.Update["email"].(string)
	id, idok := input.Update["id"].(int)
	if idok {
		if fetchresult["id"] == id {
		} else {
			code, err = e.Repo.IsEmployeeExists("", id)
			flag = true
		}
	}
	if ok {
		if fetchresult["email"].(string) == mail {
		} else {
			code, err = e.Repo.IsEmployeeExists(mail, 0)
			flag = true
		}
		if flag {
			if err != nil {
				openlog.Error(err.Error())
				return common.Response{Msg: err.Error(), Data: nil, Status: code}
			}
		}
	}
	res, code, err := e.Repo.FindAndUpdateEmployee(input.ID, input.Update)
	if err != nil {
		openlog.Error(err.Error())
		return common.Response{Msg: err.Error(), Data: nil, Status: code}
	}
	return common.Response{Msg: "Employee details updated successfully", Data: res, Status: 201}
}
