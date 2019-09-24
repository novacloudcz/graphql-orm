Feature: It should be possible to fetch embedded types

    Background: We have test company
        Given I send query:
            """
            mutation {
            deleteAllUsers
            createUser(input:{id:"test",addressRaw:"{\"street\":\"some street\",\"city\":\"some city\",\"zip\":\"aa\"}"}) { id }
            }
            """

    Scenario: Fetching embedded field address
        When I send query:
            """
            query {
            user(id:"test") { id addressRaw address { street city zip } }
            }
            """
        Then the response should be:
            """
            {
                "user": {
                    "id": "test",
                    "addressRaw": "{\"street\":\"some street\",\"city\":\"some city\",\"zip\":\"aa\"}",
                    "address": {
                        "street": "some street",
                        "city": "some city",
                        "zip": "aa"
                    }
                }
            }
            """