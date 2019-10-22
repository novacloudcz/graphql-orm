Feature: It should be fetch by number column

    Background: We have Johny and Jane
        Given I send query:
            """
            mutation {
            deleteAllUsers
            null:createUser(input:{id:"null",firstName:null}) { id }
            jane:createUser(input:{id:"jane",firstName:"Jane"}) { id }
            }
            """

    Scenario: Fetching single user by null column
        When I send query:
            """
            query { users(filter:{firstName_null:true}) { items { id firstName } count } }
            """
        Then the response should be:
            """
            {
                "users": {
                    "items": [
                        {
                            "id": "null",
                            "firstName": null
                        }
                    ],
                    "count": 1
                }
            }
            """

    Scenario: Fetching single user by non null column
        When I send query:
            """
            query { users(filter:{firstName_null:false}) { items { id firstName } count } }
            """
        Then the response should be:
            """
            {
                "users": {
                    "items": [
                        {
                            "id": "jane",
                            "firstName": "Jane"
                        }
                    ],
                    "count": 1
                }
            }
            """
