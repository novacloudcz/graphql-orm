Feature: It should be possible fetch extended query field

    Background: We have test company
        Given I send query:
            """
            mutation {
            deleteAllCompanies
            test:createCompany(input:{id:"test",name:"Test company"}) { id }
            test2:createCompany(input:{id:"test2",name:"Test2 company"}) { id }
            }
            """

    Scenario: Fetching hello world
        When I send query:
            """
            query { hello }
            """
        Then the response should be:
            """
            {
                "hello": "world"
            }
            """

    Scenario: Fetching entities using extended query
        When I send query:
            """
            query { topCompanies { id name } }
            """
        Then the response should be:
            """
            {
                "topCompanies": [
                    {
                        "id": "test",
                        "name": "Test company"
                    },
                    {
                        "id": "test2",
                        "name": "Test2 company"
                    }
                ]
            }
            """

    Scenario: Fetching entity extended fields
        When I send query:
            """
            query { company(id:"test") { id name uppercaseName } }
            """
        Then the response should be:
            """
            {
                "company": {
                    "id": "test",
                    "name": "Test company",
                    "uppercaseName": "TEST COMPANY"
                }
            }
            """
