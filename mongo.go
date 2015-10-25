package main

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"os"
	"gopkg.in/mgo.v2/bson"
)

func MongoCreate(Information *Info) {

	sess, err := mgo.Dial("mongodb://cmpe273:cmpe273@ds037234.mongolab.com:37234/cmpe273")
	if err != nil {
		fmt.Printf("Can't connect to mongo, go error %v\n", err)
		panic(err)
	}
	defer sess.Close()
	//fmt.Printf("open DB successfully\n")

	sess.SetSafe(&mgo.Safe{})
	db := sess.DB("cmpe273").C("Location")

	Information.Id = bson.NewObjectId()
	//fmt.Println("Information.Id:", Information.Id);

	err = db.Insert(&Information)
	if err != nil {
		fmt.Printf("Can't insert DB: %v\n", err)
		os.Exit(1)
	}

	var results []Info
	err = db.Find(bson.M{"_id": Information.Id}).Sort("-timestamp").All(&results)

	if err != nil {
		panic(err)
	}
	//fmt.Println("Results All: ", results)

	err = db.Find(bson.M{}).Sort("-timestamp").All(&results)

	if err != nil {
		panic(err)
	}
	//fmt.Println("Results All: ", results)
}

func MongoQuery(Id bson.ObjectId ) (Info, error) {

    var result []Info
	sess, err := mgo.Dial("mongodb://cmpe273:cmpe273@ds037234.mongolab.com:37234/cmpe273")
	if err != nil {
		fmt.Printf("Can't connect to mongo, go error %v\n", err)
		panic(err)
		return result[0], err
	}
	defer sess.Close()
	//fmt.Printf("MongoQuery: open DB successfully")
	sess.SetSafe(&mgo.Safe{})
	collection := sess.DB("cmpe273").C("Location")


	err = collection.Find(bson.M{"_id": Id}).All(&result)

	if err != nil {
		panic(err)
		return result[0], err
	}
	//fmt.Println("Results : ", result)

	return result[0],nil
}

func MongoUpdate(Information Info ) (Info, error) {

	var LocalInfo Info

	sess, err := mgo.Dial("mongodb://cmpe273:cmpe273@ds037234.mongolab.com:37234/cmpe273")
	if err != nil {
		fmt.Printf("Can't connect to mongo, go error %v\n", err)
		panic(err)
		return LocalInfo,err;
	}
	defer sess.Close()
	//fmt.Printf("MongoQuery: open DB successfully")
	sess.SetSafe(&mgo.Safe{})
	collection := sess.DB("cmpe273").C("Location")

	colQuerier := bson.M{"_id": Information.Id}

	change := bson.M{"$set": bson.M{"address": Information.Address,
									"city": Information.City,
									"state": Information.State,
									"zip": Information.Zip,
									"coordinate": bson.M{"lat":Information.Coordinate.Lat,
											"lng":Information.Coordinate.Lng}}}
	err = collection.Update(colQuerier, change)
	if err != nil {
		panic(err)
		return LocalInfo,err
	}
	Info,error := MongoQuery(Information.Id)
	return Info,error
}

func MongoRemove(Id bson.ObjectId ) (error) {

	sess, err := mgo.Dial("mongodb://cmpe273:cmpe273@ds037234.mongolab.com:37234/cmpe273")
	if err != nil {
		fmt.Printf("Can't connect to mongo, go error ", err)
		panic(err)
		return err;
	}
	defer sess.Close()
	//fmt.Printf("MongoRemove: open DB successfully")
	sess.SetSafe(&mgo.Safe{})
	collection := sess.DB("cmpe273").C("Location")

	err = collection.Remove(bson.M{"_id": Id})
	if err != nil {
		panic(err)
		return err
	}

	return nil
}
