Feature: It should be possible to mutate with relationships

    Background: We have Johny and Jane
        Given I send query:
            """
            mutation {
            deleteAllUsers
            johny:createUser(input:{id:"johny",firstName:"John",lastName:"Doe"}) { id }
            jane:createUser(input:{id:"jane",firstName:"Jane",lastName:"Siri"}) { id }
            }
            """

    Scenario: Fetching multiple items with count
        When I send query:
            """
            query { users(filter:{id:"johny"}) { items {firstName lastName createdBy updatedBy} count } }
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

    Scenario: Fetching single item
        When I send query:
            """
            query { user(id:"johny") { firstName lastName } }
            """
        Then the response should be:
            """
            {
                "user": {
                    "firstName": "John",
                    "lastName": "Doe"
                }
            }
            """

    Scenario: Fetching single item using filters
        When I send query:
            """
            query { user(filter:{id:"johny"}) { firstName lastName } }
            """
        Then the response should be:
            """
            {
                "user": {
                    "firstName": "John",
                    "lastName": "Doe"
                }
            }
            """

    Scenario: Fetching single item using IN filters
        When I send query:
            """
            query { user(filter:{id_in:["jane"]}) { firstName lastName } }
            """
        Then the response should be:
            """
            {
                "user": {
                    "firstName": "Jane",
                    "lastName": "Siri"
                }
            }
            """
    Scenario: Fetching single item using not equal filter
        When I send query:
            """
            query { user(filter:{id_ne:"johny"}) { firstName lastName } }
            """
        Then the response should be:
            """
            {
                "user": {
                    "firstName": "Jane",
                    "lastName": "Siri"
                }
            }
            """
    Scenario: Fetching with or filter should correctly compose and filter for single object
        When I send query:
            """
            query { users(filter:{OR:[{id:"johny",firstName:"Jane"}]}) { items {id} count } }
            """
        Then the response should be:
            """
            {
                "users": {
                    "items": [],
                    "count": 0
                }
            }
            """
