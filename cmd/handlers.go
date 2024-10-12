package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	models "rest_api_server/internal/Models"
	connctordb "rest_api_server/internal/connctorDB"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func startUI(c echo.Context) error {
	r := new(models.ResponseUI)

	r.Ui = "OKOK"

	return c.JSON(http.StatusOK, r)
}

func startV1(c echo.Context) error {

	r := new(models.ResposeV1)

	r.V1 = "OK"

	return c.JSON(http.StatusOK, r)

}

func get_items(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid ID format"})
	}
	items := connctordb.Get(id)
	return c.JSON(http.StatusOK, items)
}

func new_item(c echo.Context) error {
	item := new(models.Item)

	err := json.NewDecoder(c.Request().Body).Decode(&item)
	if err != nil {
		return err
	}

	result := connctordb.Save(*item)
	fmt.Print(result)
	return c.JSON(http.StatusOK, item)
}

func Get_user(c echo.Context) error {
	// fmt.Println("TEST")
	Cashe(c)

	return nil
}

// ////////////////////////////////

func Cashe(c echo.Context) {
	n, _ := io.ReadAll(c.Request().Body)

	i, err := GetID(string(n))
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to getid"})
		return
	}

	db, err := Mongodbcon("retual", "users")
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to connect to MongoDB"})
		return
	}

	st, err := SelectRequest(db, i, "retual", "users")
	if err != nil {
		HandlerFindUser(i, db, c)
		return
	} else {
		c.JSON(http.StatusOK, st)
	}

	// if u, ok := st.(user); ok {
	// fmt.Printf("User ID: %d, Name: %s, Age: %d\n", u.ID, u.NAME, u.AGE)
	// c.JSON(http.StatusOK, u)
	// } else {
	// HandlerFindUser(i, db, c)
	// c.JSON(http.StatusBadRequest, map[string]string{"error": "Data is not a user instance"})
	// }
}

func Mongodbcon(dbName string, table string) (*mongo.Client, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI("mongodb://127.0.0.1:27017")

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	// Используем параметр dbName вместо переменной db
	database := client.Database(dbName)

	// Получаем список коллекций
	collections, err := database.ListCollectionNames(ctx, bson.D{})
	if err != nil {
		return nil, fmt.Errorf("failed to list collections: %v", err)
	}

	collectionExists := false
	for _, coll := range collections {
		if coll == table {
			collectionExists = true
			break
		}
	}

	if !collectionExists {
		fmt.Printf("Collection '%s' does not exist. Creating collection...\n", table)

		// Вставляем документ для создания коллекции
		_, err := database.Collection(table).InsertOne(ctx, bson.D{{"init", "created"}})
		if err != nil {
			return nil, fmt.Errorf("failed to create collection '%s': %v", table, err)
		}
		fmt.Printf("Collection '%s' created successfully!\n", table)
	} else {
		fmt.Printf("Collection '%s' exists!\n", table)
	}

	return client, nil
}

type User struct {
	Id   int    `bson:"id"`
	Name string `bson:"name"`
	Age  int    `bson:"age"`
}

func SelectRequest(client *mongo.Client, id int, dbName, collectionName string) (*User, error) {
	db := client.Database(dbName)
	collection := db.Collection(collectionName)
	filter := bson.D{{Key: "id", Value: id}}

	var person User
	err := collection.FindOne(context.Background(), filter).Decode(&person)
	if err != nil {
		return nil, errors.New("person not found")
	}
	return &person, nil
}

func HandlerFindUser(id int, con *mongo.Client, c echo.Context) {
	// db2, err := (&DataBase{}).Init()
	// if err != nil {
	//  c.JSON(http.StatusBadRequest)
	// }
	resp := User{
		Id:   4,
		Name: "Jhon",
		Age:  44,
	}
	// err, resp := ReadByID(db2, id)
	// if err != nil {
	//  c.JSON(http.StatusGatewayTimeout)

	coll := con.Database("retual").Collection("users")
	_, err := coll.InsertOne(context.TODO(), resp)
	if err != nil {
		c.JSON(http.StatusNotFound, map[string]string{"error": "No select"})
	} else {
		c.JSON(http.StatusOK, resp)
	}
}

func GetID(jsonStr string) (int, error) {
	dec := json.NewDecoder(strings.NewReader(jsonStr))
	_, err := dec.Token()
	if err != nil {
		return 0, nil
	}
	var idStr string
	for dec.More() {
		var key json.Token
		key, err = dec.Token()
		if err != nil {
			return 0, nil
		}
		keyStr, ok := key.(string)
		if !ok {
			return 0, nil
		}

		if keyStr == "id" {
			var value json.Token
			value, err = dec.Token()
			if err != nil {
				return 0, nil
			}
			idStr, ok = value.(string)
			if !ok {
				return 0, nil
			}

			break
		} else {
			_, err = dec.Token()
			if err != nil {
				return 0, nil
			}
		}
	}
	_, err = dec.Token()
	if err != nil {
		return 0, nil
	}
	num, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, nil
	}
	return num, nil
}
