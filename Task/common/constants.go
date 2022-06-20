package common

type Response struct {
	Msg    string      `json:"_msg"`
	Status int         `json:"_status"`
	Data   interface{} `json:"data"`
}

type CreateInput struct {
	Create map[string]interface{}
}
type FetchDepartmentInput struct {
	ID string
}

type FetchAllInput struct {
	Page, Size, Filters, Sort string
}
type DeleteInput struct {
	ID string
}
type UpdateInput struct {
	ID     string
	Update map[string]interface{}
}

var Departmentid = "department_id"
