Feature: It should be possible to fetch with q parameter

    Background: We have Johny and Jane
        Given I send query:
            """
            mutation {
            deleteAllUsers
            johny:createUser(input:{id:"johny",firstName:"John",lastName:"Doe"}) { id }
            jane:createUser(input:{id:"jane",firstName:"Jane",lastName:"Siri"}) { id }
            }
            """

    Scenario: Fetching list by q parameter
        When I send query:
            """
            query { users(q:"John") { items {firstName lastName createdBy updatedBy} count } }
            """
        Then the response should be:
            """
            {
                "users": {
                    "items": [
                        {
                            "firstName": "John",
                            "lastName": "Doe",
                            "createdBy": null,
                            "updatedBy": null
                        }
                    ],
                    "count": 1
                }
            }
            """

    Scenario: Fetching single item by q parameter
        When I send query:
            """
            query { user(q:"John") { firstName lastName createdBy updatedBy} }
            """
        Then the response should be:
            """
            {
                "user": {
                    "firstName": "John",
                    "lastName": "Doe",
                    "createdBy": null,
                    "updatedBy": null
                }
            }
            """
