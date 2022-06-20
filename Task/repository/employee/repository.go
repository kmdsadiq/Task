package employeerepository

import (
	"Task/common"
	"context"
	"errors"
	"time"

	"github.com/go-chassis/openlog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type EmployeeRepo struct {
	DbClient     *mongo.Client
	DatabaseName string
}

var employeecollection = "employee"

//EmployeeInsert function will insert the data into database
func (ec *EmployeeRepo) EmployeeInsert(meta map[string]interface{}) (map[string]interface{}, int, error) {
	collection := ec.DbClient.Database(ec.DatabaseName).Collection(employeecollection)

	meta["createdon"] = time.Now().Unix()
	res, err := collection.InsertOne(context.Background(), meta)
	if err != nil {
		openlog.Error(err.Error())
		return make(map[string]interface{}), 500, errors.New("Internal Server Error")
	}
	id := res.InsertedID.(primitive.ObjectID).Hex()
	result, _, err := ec.FindEmployee(id)
	if err != nil {
		openlog.Error(err.Error())
		return result, 500, errors.New("Internal Server Error")
	}
	return result, 0, nil
}

//Find function will find the document by id returns that document
func (ec *EmployeeRepo) FindEmployee(id string) (map[string]interface{}, int, error) {
	collection := ec.DbClient.Database(ec.DatabaseName).Collection(employeecollection)
	result := make(map[string]interface{})
	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		openlog.Error(err.Error())
		return result, 400, errors.New("Invalid Id")
	}
	err = collection.FindOne(context.TODO(), bson.M{"_id": docID}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return result, 404, errors.New("No Documents Found")
		}
		openlog.Error(err.Error())
		return result, 500, errors.New("Internal Server Error")
	}
	return result, 0, nil
}

//FindAndUpdate function will find the documentation by id and update the data for that document
func (ec *EmployeeRepo) FindAndUpdateEmployee(id string, document map[string]interface{}) (map[string]interface{}, int, error) {
	collection := ec.DbClient.Database(ec.DatabaseName).Collection(employeecollection)
	update := make(map[string]interface{})
	update["$set"] = document
	result := make(map[string]interface{})
	docID, convErr := primitive.ObjectIDFromHex(id)
	if convErr != nil {
		openlog.Error(convErr.Error())
		return result, 400, errors.New("Invalid ID")
	}
	document["updatedon"] = time.Now().Unix()
	err := collection.FindOneAndUpdate(context.TODO(), bson.M{"_id": docID}, update).Decode(&result)
	if err != nil {
		openlog.Error(err.Error())
		return result, 500, errors.New("Internal Server Error")
	}
	result, _, err = ec.FindEmployee(id)
	if err != nil {
		openlog.Error(err.Error())
		return result, 500, errors.New("Internal Server Error")
	}
	return result, 0, nil
}

//FindByFiltersAndPagenation will find the document by filters and pagenation and returns the data
func (ec *EmployeeRepo) FindByFiltersAndPagenation(page int64, size int64, filters map[string]interface{}, sort map[string]interface{}) ([]map[string]interface{}, int, error) {

	options := *options.Find()

	collection := ec.DbClient.Database(ec.DatabaseName).Collection(employeecollection)
	var result []map[string]interface{}

	cursor, err := collection.Find(context.TODO(), filters, options.SetLimit(size), options.SetSkip((page-1)*size), options.SetSort(sort))
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return result, 0, nil
		}
		openlog.Error(err.Error())
		return result, 500, errors.New("Internal Server Error")
	}
	for cursor.Next(context.TODO()) {
		var doc map[string]interface{}
		err := cursor.Decode(&doc)
		if err != nil {
			openlog.Error(err.Error())
			return result, 500, errors.New("Internal Server Error")
		}
		result = append(result, doc)
	}
	return result, 0, nil
}

//FindByFilters will find the document by filters and pagenation and returns the data
func (ec *EmployeeRepo) FindByFilters(filters map[string]interface{}, sort map[string]interface{}) ([]map[string]interface{}, int, error) {

	options := *options.Find()

	collection := ec.DbClient.Database(ec.DatabaseName).Collection(employeecollection)
	var result []map[string]interface{}

	cursor, err := collection.Find(context.TODO(), filters, options.SetSort(sort))
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return result, 200, nil
		}
		return result, 500, errors.New("Internal server error")
	}
	for cursor.Next(context.TODO()) {
		var doc map[string]interface{}
		err := cursor.Decode(&doc)
		if err != nil {
			openlog.Error(err.Error())
			return result, 500, errors.New("Internal Server Error")
		}
		result = append(result, doc)
	}
	return result, 200, nil
}

//IsBoilerPlateNameExists checks for any Boilerplate conflict
func (ec *EmployeeRepo) IsEmployeeExists(email string, id float64) (int, error) {
	collection := ec.DbClient.Database(ec.DatabaseName).Collection(employeecollection)
	result := make(map[string]interface{})
	flag := false
	var err error
	if email != "" {
		err = collection.FindOne(context.TODO(), bson.M{"email": email}).Decode(&result)
	} else if id > 0 {
		flag = true
		err = collection.FindOne(context.TODO(), bson.M{"id": id}).Decode(&result)
	}
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return 404, nil
		}
		openlog.Error(err.Error())
		return 500, errors.New("Internal Server Error")
	}
	if flag {
		return 409, errors.New("Email already exists")
	} else {
		return 409, errors.New("Employee id already exists")
	}
}

//Delete will allows to delete the  document
func (ec *EmployeeRepo) DeleteEmployee(id string) (map[string]interface{}, int, error) {
	collection := ec.DbClient.Database(ec.DatabaseName).Collection(employeecollection)
	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		openlog.Error(err.Error())
		return nil, 400, errors.New("Invalid id")
	}
	result := make(map[string]interface{})
	err = collection.FindOneAndDelete(context.TODO(), bson.M{"_id": docID}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, 404, errors.New("No Documents Found for given id")
		}
		openlog.Error(err.Error())
		return nil, 500, errors.New("Internal Server Error")
	}
	return result, 0, nil
}

func (ec *EmployeeRepo) FindDepartment(depid float64) (map[string]interface{}, int, error) {
	collection := ec.DbClient.Database(ec.DatabaseName).Collection(employeecollection)
	result := make(map[string]interface{})
	err := collection.FindOne(context.TODO(), bson.M{common.Departmentid: depid}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return result, 404, nil
		}
		openlog.Error(err.Error())
		return result, 500, errors.New("Internal Server Error")
	}
	return result, 0, nil
}
