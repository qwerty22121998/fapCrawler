package main

import (
	"net/http"
	"github.com/PuerkitoBio/goquery"

	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const URL_CLASS_DEFAULT = "http://fap.fpt.edu.vn/Course/Groups.aspx?group="
const TOKEN = "ASP.NET_SessionId=t1nccdovtrtf52nks1doucnv"
const DB_NAME = "fuHL"
const C_NAME = "student"



func getSession() (session *mgo.Session, err error) {
	session, err = mgo.Dial("localhost")
	if err != nil {
		fmt.Println("Can't connect to db")
		return nil, err
	}
	return session, nil
}
func getCollections(session *mgo.Session) *mgo.Collection {
	return session.DB(DB_NAME).C(C_NAME)
}

type Student struct {
	ObID bson.ObjectId `bson:"_id"`
	Name interface{} `bson:"name,omitempty"`
	Id   string	`bson:"id"`
}

func getClassInfomartion(id string, collection *mgo.Collection) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", URL_CLASS_DEFAULT+id, nil)

	if err != nil {
		return
	}

	req.Header.Add("Cookie", TOKEN)
	resp, err := client.Do(req)

	if err != nil {
		return
	}
	defer resp.Body.Close()
	stu := make([]Student, 0)


	doc, _ := goquery.NewDocumentFromReader(resp.Body)
	stus := doc.Find(".sortable").Find("tbody").Find("tr").Find("td").Nodes
	class := doc.Find("#ctl00_mainContent_divGroup").Find("table").Find("tbody").Find("tr").Find("td").Find("b").Text()


	fmt.Println("CLASS", class)

	for i := 0; i < len(stus); i += 7 {
		id := stus[i+3].FirstChild.Data
		name := []string{stus[i+4].FirstChild.Data, stus[i+5].FirstChild.Data, stus[i+6].FirstChild.Data}
		stu = append(stu, Student{bson.NewObjectId(), name, id})
	}

	for _, v := range stu {
		fmt.Println(v)

		err := collection.Insert(v)
		if err != nil {
			fmt.Println(err)
			return

		}
	}

}

func main() {



	s, _ := getSession()
	c :=getCollections(s)
	getClassInfomartion("3113", c)
}
