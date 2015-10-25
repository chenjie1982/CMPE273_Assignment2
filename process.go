package main
import (
	"net/http"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"io"
	"encoding/json"
	"github.com/gorilla/mux"
)

type Info struct {
	Id          bson.ObjectId `json:"id"      bson:"_id,omitempty"`
	Name        string `json:"name"    bson:"name"`
	Address     string `json:"address" bson:"address"`
	City        string `json:"city"    bson:"city"`
	State       string `json:"state"   bson:"state"`
	Zip         string `json:"zip"     bson:"zip"`
	Coordinate struct {
		Lat float64 `json:"lat" bson:"lat"`
		Lng float64 `json:"lng" bson:"lng"`
	} `json:"coordinate"        bson:"coordinate"`
}

type AddressInfo struct {
	Name        string `json:"name"    bson:"name"`
	Address     string `json:"address" bson:"address"`
	City        string `json:"city"    bson:"city"`
	State       string `json:"state"   bson:"state"`
	Zip         string `json:"zip"     bson:"zip"`
}

type Service struct{}


func Create(w http.ResponseWriter, r *http.Request) {

	var args AddressInfo
	var reply Info

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	if err := json.Unmarshal(body, &args); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
		fmt.Println("Error! Can't unmarshal Json from create emulator request.", body)
		return
	}
	addr := args.Address+","+args.City+","+args.State+","+args.Zip

	Information,err := QueryInfo(addr);

	reply.Address = args.Address;
	reply.City = args.City;
	reply.State = args.State;
	reply.Zip = args.Zip;
	reply.Name = args.Name;
	reply.Coordinate.Lat = Information.Coordinate.Lat
	reply.Coordinate.Lng = Information.Coordinate.Lng
	MongoCreate(&reply)
	//fmt.Println("Information.Id:"+reply.Id)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(reply); err != nil {
		panic(err)
	}
	return
}

func Update(w http.ResponseWriter, r *http.Request) {

	var args AddressInfo
	var reply Info

	vars := mux.Vars(r)
	reply.Id = bson.ObjectIdHex(vars["location_id"])

	//fmt.Println(reply.Id)

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	if err := json.Unmarshal(body, &args); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
		fmt.Println("Error! Can't unmarshal Json from create emulator request.", body)
		return
	}

	addr := args.Address+","+args.City+","+args.State+","+args.Zip

	Information,err := QueryInfo(addr);

	//fmt.Println(Information)

	reply.Address = args.Address;
	reply.City = args.City;
	reply.State = args.State;
	reply.Zip = args.Zip;
	reply.Coordinate.Lat = Information.Coordinate.Lat
	reply.Coordinate.Lng = Information.Coordinate.Lng

	Information, err = MongoUpdate(reply)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}
	reply.Name = Information.Name

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(reply); err != nil {
		panic(err)
	}

}

func Query(w http.ResponseWriter, r *http.Request) {

	var reply Info
	var err error
	vars := mux.Vars(r)
	reply.Id = bson.ObjectIdHex(vars["location_id"])


	reply, err = MongoQuery(reply.Id )
	if err != nil {

		fmt.Printf(err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(reply); err != nil {
		panic(err)
	}
}

func Remove(w http.ResponseWriter, r *http.Request) {

	var reply Info

	vars := mux.Vars(r)
	reply.Id = bson.ObjectIdHex(vars["location_id"])

	err := MongoRemove(reply.Id )
	if err != nil {

		fmt.Printf(err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	return

}