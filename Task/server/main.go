package main

import (
	"log"
	"net/http"

	"Task/handlers"

	"github.com/go-chassis/openlog"
	"github.com/gorilla/mux"

	"Task/database"
	departmentrepository "Task/repository/department"
	repository "Task/repository/employee"

	service "Task/service/departmentservice"
	employeeservice "Task/service/employeeservice"
)

func GetService(dbname string) (service.Service, employeeservice.EmployeeService) {
	employeerepo := repository.EmployeeRepo{DbClient: database.GetClient(), DatabaseName: dbname}
	departmentrepo := departmentrepository.DepartmentRepo{DbClient: database.GetClient(), DatabaseName: dbname}

	return service.Service{DepartmentRepo: departmentrepo, Repo: employeerepo}, employeeservice.EmployeeService{Repo: employeerepo, DepartmentRepo: departmentrepo}
}
func main() {
	r := mux.NewRouter()
	err := database.Connect()
	if err != nil {
		openlog.Error(err.Error())
		return
	}
	service, employee := GetService("Users")
	h := handlers.Handler{Service: service, Employeeservice: employee}
	r.HandleFunc("/departments/{id}", h.GetDepartment).Methods("GET")
	r.HandleFunc("/departments", h.CreateDepartment).Methods("POST")
	r.HandleFunc("/departments/{id}", h.UpdateDepartment).Methods("PUT")
	r.HandleFunc("/departments/{id}", h.DeleteDepartment).Methods("DELETE")
	r.HandleFunc("/departments", h.FetchAllDepartments).Methods("GET")
	//
	r.HandleFunc("/employees/{id}", h.GetEmployeebyid).Methods("GET")
	r.HandleFunc("/employees", h.CreateEmployee).Methods("POST")
	r.HandleFunc("/employees/{id}", h.UpdateEmployee).Methods("PUT")
	r.HandleFunc("/employees/{id}", h.DeleteEmployee).Methods("DELETE")
	r.HandleFunc("/employees", h.FetchAllEmployees).Methods("GET")

	openlog.Info("Started listening at http://localhost:8070")
	log.Fatal(http.ListenAndServe(":8070", r))

}
