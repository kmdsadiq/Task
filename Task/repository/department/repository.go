package repository

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

type DepartmentRepo struct {
	DbClient     *mongo.Client
	DatabaseName string
}

var departmentcollection = "department"

//Insert function will insert the data into database
func (dr *DepartmentRepo) DepartmentInsert(meta map[string]interface{}) (map[string]interface{}, int, error) {
	collection := dr.DbClient.Database(dr.DatabaseName).Collection(departmentcollection)

	meta["createdon"] = time.Now().Unix()
	res, err := collection.InsertOne(context.Background(), meta)
	if err != nil {
		openlog.Error(err.Error())
		return make(map[string]interface{}), 500, errors.New("Internal Server Error")
	}
	id := res.InsertedID.(primitive.ObjectID).Hex()
	result, _, err := dr.DepartmentFind(id)
	if err != nil {
		openlog.Error(err.Error())
		return result, 500, errors.New("Internal Server Error")
	}
	return result, 0, nil
}

//Find function will find the document by id returns that document
func (dr *DepartmentRepo) DepartmentFind(id string) (map[string]interface{}, int, error) {
	collection := dr.DbClient.Database(dr.DatabaseName).Collection(departmentcollection)
	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		openlog.Error(err.Error())
		return nil, 400, errors.New("Invalid Id")
	}
	matchStage := bson.D{{"$match", bson.D{{"_id", docID}}}}
	result := make([]map[string]interface{}, 0)
	lookUpStage_servicegroup := bson.D{{"$lookup", bson.D{{"from", "employee"}, {"localField", "_id"}, {"foreignField", common.Departmentid}, {"as", "employees"}}}}
	db_result, err := collection.Aggregate(context.TODO(), mongo.Pipeline{matchStage, lookUpStage_servicegroup})

	if err != nil {
		openlog.Error(err.Error())
		return nil, 500, errors.New("Internal Server Error")
	}
	if err = db_result.All(context.TODO(), &result); err != nil {
		openlog.Error(err.Error())
		return nil, 500, errors.New("Internal Server Error")
	}
	if len(result) == 0 {
		return nil, 404, errors.New("No Documents")
	}
	return result[0], 0, nil
}

//FindAndUpdate function will find the documentation by id and update the data for that document
func (dr *DepartmentRepo) FindAndUpdateDepartment(id string, document map[string]interface{}) (map[string]interface{}, int, error) {
	collection := dr.DbClient.Database(dr.DatabaseName).Collection(departmentcollection)
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
	result, _, err = dr.DepartmentFind(id)
	if err != nil {
		openlog.Error(err.Error())
		return result, 500, errors.New("Internal Server Error")
	}
	return result, 0, nil
}

//FindByFiltersAndPagenation will find the document by filters and pagenation and returns the data
func (dr *DepartmentRepo) FindByFiltersAndPagenation(page int64, size int64, filters map[string]interface{}, sort map[string]interface{}) ([]map[string]interface{}, int, error) {

	options := *options.Find()

	collection := dr.DbClient.Database(dr.DatabaseName).Collection(departmentcollection)
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
func (dr *DepartmentRepo) FindByFilters(filters map[string]interface{}, sort map[string]interface{}) ([]map[string]interface{}, int, error) {

	options := *options.Find()

	collection := dr.DbClient.Database(dr.DatabaseName).Collection(departmentcollection)
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

func (dr *DepartmentRepo) IsDepartmentExists(name string, id float64) (int, error) {
	collection := dr.DbClient.Database(dr.DatabaseName).Collection(departmentcollection)
	result := make(map[string]interface{})
	flag := false
	var err error
	if name != "" {
		err = collection.FindOne(context.TODO(), bson.M{"department name": name}).Decode(&result)
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
		return 409, errors.New("Department already exists")
	} else {
		return 409, errors.New("Department id already exists")
	}
}

//DeleteDepartment will allows to delete the  document
func (dr *DepartmentRepo) DeleteDepartment(id string) (map[string]interface{}, int, error) {
	collection := dr.DbClient.Database(dr.DatabaseName).Collection(departmentcollection)
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

func (dr *DepartmentRepo) DepartmentExists(id float64) (map[string]interface{}, int, error) {
	collection := dr.DbClient.Database(dr.DatabaseName).Collection(departmentcollection)
	result := make(map[string]interface{})
	err := collection.FindOne(context.TODO(), bson.M{"id": id}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return result, 404, errors.New("No Documents Found")
		}
		openlog.Error(err.Error())
		return result, 500, errors.New("Internal Server Error")
	}
	return result, 0, nil
}
