package models

import (
  "fmt"
  "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
  "github.com/aws/aws-sdk-go/service/dynamodb"
  "github.com/aws/aws-sdk-go/aws"
  //"os"
  "github.com/google/uuid"

)

type Estimate struct {
  Eid           string        `dynamodbav:"eid"`
  Title         string        `dynamodbav:"title"`
  Owner         string        `dynamodbav:"ownerid"`
  Hd            string        `dynamodbav:"hd"`
  Unit          string        `dynamodbav:"unitname"`
  Description   string        `dynamodbav:"description"`
  AvgMinimum    float64       `dynamodbav:"minimum"`
  AvgMaximum    float64       `dynamodbav:"maximum"`
  Actual        float64       `dynamodbav:"actual"`
  BrierScore    float64       `dynamodbav:"brierscore"`
  Concluded     bool          `dynamodbav:"concluded"`
  ConcludedTime string        `dynamodbav:"concludetime"`
}

func CreateEstimate (title string, description string, unit string, hd string, owner string) (eid string){

  		euuid := uuid.New()
  		item := Estimate{
  				Eid: euuid.String(),
          Owner: owner,
          Hd: hd,
  		    Title: title,
  		    Description: description,
          Unit: unit,
  		}

  		PutItem(item, "estimates-tf")
      fmt.Println(unit)

      return euuid.String()
}

func UpdateEstimate (eid string, title string, description string, unit string, user string) {

  //Primary key for update query
  key := map[string]*dynamodb.AttributeValue {
    "eid": {
      S: aws.String(eid),
    },
  }

  expressionattrvalues:= map[string]*dynamodb.AttributeValue {
    ":t": {
      S: aws.String(title),
    },
    ":d": {
      S: aws.String(description),
    },
    ":unit": {
      S: aws.String(unit),
    },
    ":user": {
      S: aws.String(user),
    },
  }

  updateexpression := "SET title = :t, description = :d, unitname = :unit"
  conditionexpression := "ownerid = :user"

  UpdateItem(key, updateexpression, expressionattrvalues, "estimates-tf", conditionexpression)

}

func GetEstimate (eid string) (e Estimate) {

  result := GetPrimaryItem(eid, "eid", "estimates-tf")

  e = Estimate{}

  err := dynamodbattribute.UnmarshalMap(result.Item, &e)

  if err != nil {
    panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
  }

  if e.Eid == "" {
      fmt.Println("Could not find that scenario.")
      return
  }

  return e

}

func ListEstimates(user string) (e []Estimate) {

  result := GetPrimaryIndexItem(user, "ownerid", "ownerid-index", "estimates-tf")

  e = []Estimate{}

  err := dynamodbattribute.UnmarshalListOfMaps(result.Items, &e)

  if err != nil {
    panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
  }

  return e

}

func DeleteEstimate(eid string, owner string) {

  DeletePrimaryItem(eid, "eid", "estimates-tf", "ownerid", owner)

  fmt.Println("Deleted estimate.", eid)

  DeleteEstimateRanges(eid)

}

func DeleteEstimateRanges(eid string) {

    es := ViewEstimateResults(eid)


    for _, v  := range es {
      fmt.Println("Deleting: ", v.Eid, v.User)
      DeleteCompositeIndexItem(v.Eid, v.User, "eid", "user", "ranges-tf")
    }


    fmt.Println("Deleted forecasts associated with scenario.")

}
