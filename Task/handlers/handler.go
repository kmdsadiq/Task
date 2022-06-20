package handlers

import (
	"Task/common"
	validator "Task/payloadvalidator"
	service "Task/service/departmentservice"
	employeeservice "Task/service/employeeservice"
	"encoding/json"
	"net/http"

	"github.com/go-chassis/openlog"
	"github.com/gorilla/mux"
)

type Handler struct {
	Service         service.Service
	Employeeservice employeeservice.EmployeeService
}

type Response struct {
	Msg    string      `json:"_msg"`
	Status int         `json:"_status"`
	Data   interface{} `json:"data"`
}

func (h *Handler) CreateDepartment(w http.ResponseWriter, r *http.Request) {
	openlog.Info("Got a request to create Department")
	w.Header().Set("Content-Type", "application/json")

	department := make(map[string]interface{})
	json.NewDecoder(r.Body).Decode(&department)
	valres, err := validator.ValidatePaylaod("./../payloadschemas/createdepartment.json", department)
	if err != nil {
		openlog.Error(err.Error())
		response := Response{Msg: err.Error(), Data: valres, Status: 400}
		json.NewEncoder(w).Encode(response)
		return
	}
	input := common.CreateInput{Create: department}
	res := h.Service.CreateDepartment(input)
	w.WriteHeader(res.Status)
	json.NewEncoder(w).Encode(res)
}

func (h *Handler) GetDepartment(w http.ResponseWriter, r *http.Request) {
	openlog.Info("Got a request to fetch Department")
	// set header.
	w.Header().Set("Content-Type", "application/json")
	var params = mux.Vars(r)
	id := params["id"]
	input := common.FetchDepartmentInput{ID: id}
	res := h.Service.GetDepartmentbyid(input)
	w.WriteHeader(res.Status)
	json.NewEncoder(w).Encode(res)
}

// FetchAllDatamodelsByPagenation function will helps to get the field by considering page number, size and filters.
func (h *Handler) FetchAllDepartments(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	size := r.URL.Query().Get("size")
	filters := r.URL.Query().Get("filters")
	sort := r.URL.Query().Get("sort")
	input := common.FetchAllInput{Page: page, Size: size, Filters: filters, Sort: sort}
	w.Header().Set("Content-Type", "application/json")
	res := h.Service.FetchAllDepartment(input)
	w.WriteHeader(res.Status)
	json.NewEncoder(w).Encode(res)
}

func (h *Handler) DeleteDepartment(w http.ResponseWriter, r *http.Request) {
	openlog.Info("Got a request to Delete Department")
	// set header.
	w.Header().Set("Content-Type", "application/json")
	var params = mux.Vars(r)
	id := params["id"]
	input := common.DeleteInput{ID: id}
	res := h.Service.DeleteDepartment(input)
	w.WriteHeader(res.Status)
	json.NewEncoder(w).Encode(res)
}

func (h *Handler) UpdateDepartment(w http.ResponseWriter, r *http.Request) {
	openlog.Info("Got a request to Update Department")
	w.Header().Set("Content-Type", "application/json")
	var params = mux.Vars(r)
	id := params["id"]
	department := make(map[string]interface{})
	json.NewDecoder(r.Body).Decode(&department)
	valres, err := validator.ValidatePaylaod("./../payloadschemas/updatedepartment.json", department)
	if err != nil {
		openlog.Error(err.Error())
		response := Response{Msg: err.Error(), Data: valres, Status: 400}
		json.NewEncoder(w).Encode(response)
		return
	}
	input := common.UpdateInput{ID: id, Update: department}
	res := h.Service.UpdateDepartment(input)
	w.WriteHeader(res.Status)
	json.NewEncoder(w).Encode(res)
}

//

func (h *Handler) CreateEmployee(w http.ResponseWriter, r *http.Request) {
	openlog.Info("Got a request to create employee")
	w.Header().Set("Content-Type", "application/json")

	department := make(map[string]interface{})
	json.NewDecoder(r.Body).Decode(&department)
	valres, err := validator.ValidatePaylaod("./../payloadschemas/createemployee.json", department)
	if err != nil {
		openlog.Error(err.Error())
		response := Response{Msg: err.Error(), Data: valres, Status: 400}
		json.NewEncoder(w).Encode(response)
		return
	}
	input := common.CreateInput{Create: department}
	res := h.Employeeservice.CreateEmployee(input)
	w.WriteHeader(res.Status)
	json.NewEncoder(w).Encode(res)
}

func (h *Handler) GetEmployeebyid(w http.ResponseWriter, r *http.Request) {
	openlog.Info("Got a request to fetch Department")
	// set header.
	w.Header().Set("Content-Type", "application/json")
	var params = mux.Vars(r)
	id := params["id"]
	input := common.FetchDepartmentInput{ID: id}
	res := h.Employeeservice.GetEmployeebyid(input)
	w.WriteHeader(res.Status)
	json.NewEncoder(w).Encode(res)
}

func (h *Handler) FetchAllEmployees(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	size := r.URL.Query().Get("size")
	filters := r.URL.Query().Get("filters")
	sort := r.URL.Query().Get("sort")
	input := common.FetchAllInput{Page: page, Size: size, Filters: filters, Sort: sort}
	w.Header().Set("Content-Type", "application/json")
	res := h.Employeeservice.FetchAllEmployees(input)
	w.WriteHeader(res.Status)
	json.NewEncoder(w).Encode(res)
}

func (h *Handler) DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	openlog.Info("Got a request to Delete Employee")
	// set header.
	w.Header().Set("Content-Type", "application/json")
	var params = mux.Vars(r)
	id := params["id"]
	input := common.DeleteInput{ID: id}
	res := h.Employeeservice.DeleteEmployee(input)
	w.WriteHeader(res.Status)
	json.NewEncoder(w).Encode(res)
}

func (h *Handler) UpdateEmployee(w http.ResponseWriter, r *http.Request) {
	openlog.Info("Got a request to Update Employee")
	w.Header().Set("Content-Type", "application/json")
	var params = mux.Vars(r)
	id := params["id"]
	department := make(map[string]interface{})
	json.NewDecoder(r.Body).Decode(&department)
	valres, err := validator.ValidatePaylaod("./../payloadschemas/updateemployee.json", department)
	if err != nil {
		openlog.Error(err.Error())
		response := Response{Msg: err.Error(), Data: valres, Status: 400}
		json.NewEncoder(w).Encode(response)
		return
	}
	input := common.UpdateInput{ID: id, Update: department}
	res := h.Employeeservice.UpdateEmployee(input)
	w.WriteHeader(res.Status)
	json.NewEncoder(w).Encode(res)
}
