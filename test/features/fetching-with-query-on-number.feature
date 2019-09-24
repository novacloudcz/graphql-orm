Feature: It should be fetch by number column

    Background: We have Johny and Jane
        Given I send query:
            """
            mutation {
            deleteAllUsers
            johny:createUser(input:{id:"johny",firstName:"John",lastName:"Doe",code:5554321}) { id }
            jane:createUser(input:{id:"jane",firstName:"Jane",lastName:"Siri",code:1234}) { id }
            }
            """

    Scenario: Fetching single user by code column prefix
        When I send query:
            """
            query { users(q:"555") { items { code firstName lastName createdBy updatedBy} count } }
            """
        Then the response should be:
            """
            {
                "users": {
                    "items": [
                        {
                            "code": 5554321,
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
