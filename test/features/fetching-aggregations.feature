Feature: Collection aggregations should be fetchable from entity ResultType

    Background: We have init data set
        Given I send query:
            """
 mutation init {
  deleteAllUsers
  deleteAllCompanies
  c1:createCompany(input:{id:"blah",name:"blah"}){id}
  u1:createUser(input:{
    firstName:"John",
    lastName:"Doe"
    salary: 1200,
    employersIds:["blah"]
  }){id}
  u2:createUser(input:{
    firstName:"Jane",
    lastName:"Doe"
    salary:1502
  }){id}
  u3:createUser(input:{
    firstName:"John",
    lastName:"Snow"
    salary:1201
  }){id}
}
            """
        Then the error should be:
            """
            null
            """

    Scenario: Fetching all users should return aggregations for all
        When I send query:
            """
query queryUsers {
  users(limit:1,sort:[{salary:ASC}]) {
    items{
      firstName
      lastName
      salary
    }
    aggregations {
      salaryMin
      salaryMax
      salaryAvg
      salarySum
    }
  }
}
            """
        Then the response should be:
            """
{
  "users": {
    "items": [
      {
        "firstName": "John",
        "lastName": "Doe",
        "salary": 1200
      }
    ],
    "aggregations": {
      "salaryMin": 1200,
      "salaryMax": 1502,
      "salaryAvg": 1301,
      "salarySum": 3903
    }
  }
}
            """
    Scenario: Fetching users with basic filter
        When I send query:
            """
query queryUsersWithBasicFilter {
  users(sort:[{salary:ASC}],filter:{firstName:"John"}) {
    items{
      firstName
      lastName
      salary
    }
    aggregations {
      salaryMin
      salaryMax
      salaryAvg
    }
  }
}
            """
        Then the response should be:
            """
{
  "users": {
    "items": [
      {
        "firstName": "John",
        "lastName": "Doe",
        "salary": 1200
      },
      {
        "firstName": "John",
        "lastName": "Snow",
        "salary": 1201
      }
    ],
    "aggregations": {
      "salaryMin": 1200,
      "salaryMax": 1201,
      "salaryAvg": 1200.5
    }
  }
}
            """
        And the error should be:
            """
            null
            """

Scenario: Fetching users with relationship filter
        When I send query:
            """
query queryUsersWithRelationshipFilter {
  users(sort:[{salary:ASC}],filter:{firstName:"John",employers:{name:"blah"}}) {
    items {
      firstName
      lastName
    }
    aggregations {
      salaryMin
      salaryMax
      salaryAvg
      salarySum
    }
  }
}
            """
        Then the response should be:
            """
{
  "users": {
    "items": [
      {
        "firstName": "John",
        "lastName": "Doe"
      }
    ],
    "aggregations": {
      "salaryMin": 1200,
      "salaryMax": 1200,
      "salaryAvg": 1200,
      "salarySum": 1200
    }
  }
}
            """
        And the error should be:
            """
            null
            """