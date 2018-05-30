package models

import (
  "fmt"
  "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
  "github.com/aws/aws-sdk-go/service/dynamodb"
  "github.com/aws/aws-sdk-go/aws"
  //"os"
  "github.com/google/uuid"
  "time"

)

type Rank struct {
  Question
  Options       []string      `dynamodbav:"options"`
}

func CreateRank (title string, description string, options []string,  hd string, owner string) (rid string){

      t := time.Now()

  		ruuid := uuid.New()
  		item := Rank{
          Question: Question{
            Id: ruuid.String(),
            OwnerID: owner,
            Hd: hd,
            Title: title,
            Description: description,
            Records: []string{t.Format("2006-01-02") + ": Created.", },
            URL: "rank/" + ruuid.String(),
          },
  				Options: options,
  		}

  		PutItem(item, "questions-tf")

  		fmt.Println("Successfully added.")

      return ruuid.String()

}

func (r Rank) GetURL() (url string) {

  return "/view/rank/" + r.Id

}

func UpdateRank (rid string, title string, description string, options []string, user string) {

  //Key for the table
  key := map[string]*dynamodb.AttributeValue {
    "id": {
      S: aws.String(rid),
    },
  }

  //Changing this into a list of attributes.
  o, _ := dynamodbattribute.MarshalList(options)

  //Make our list of "expressions"
  expressionattrvalues:= map[string]*dynamodb.AttributeValue {
    ":t": {
      S: aws.String(title),
    },
    ":d": {
      S: aws.String(description),
    },
    ":o": {
      L: o,
    },
    ":user": {
      S: aws.String(user),
    },
  }

  //Case issue. Options has mixed case in other tables, should fix on production launch. See #24
  updateexpression := "SET title = :t, description = :d, options = :o"

  //Enforce moderator
  conditionexpression := "ownerid = :user"

  UpdateItem(key, updateexpression, expressionattrvalues, "questions-tf", conditionexpression)

  fmt.Println("Updated rank.")


}

func GetRank (rid string) (r Rank) {

  result := GetPrimaryItem(rid, "id", "questions-tf")

  r = Rank{}

  err := dynamodbattribute.UnmarshalMap(result.Item, &r)

  if err != nil {
    panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
  }

  if r.Question.Id == "" {
      fmt.Println("Could not find that scenario.")
      return
  }

  return r

}

func ListRanks(user string) (r []Rank) {

  result := GetPrimaryIndexItem(user, "ownerid", "ownerid-index", "questions-tf")

  r = []Rank{}

  err := dynamodbattribute.UnmarshalListOfMaps(result.Items, &r)

  if err != nil {
    panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
  }

  return r

}

func DeleteRank(rid string, owner string) {

  DeletePrimaryItem(rid, "id", "questions-tf", "ownerid", owner)

  fmt.Println("Deleted rank.", rid)

//  DeleteRankRanges(rid)

}

func ViewRankResults (rid string) (s []Sort) {
  //Need to do a HD check here to prevent IDOR.

    result := GetPrimaryIndexItem(rid, "id", "id-index", "answers-tf")

    s = []Sort{}

    err := dynamodbattribute.UnmarshalListOfMaps(result.Items, &s)

    if err != nil {
      panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
    }

    return s
}
