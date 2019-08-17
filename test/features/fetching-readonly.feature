Feature: It should be possible to fetch readonly field

    Background: We have test company
        Given I send query:
            """
            mutation {
            deleteAllCompanies
            test:createCompany(input:{id:"test",name:"Test company"}) { id }
            }
            """

    Scenario: Fetching single item
        When I send query:
            """
            query { company { id name reviews { id } } }
            """
        Then the response should be:
            """
            {
                "company": {
                    "id": "test",
                    "name": "Test company",
                    "reviews": [
                        {
                            "id": "1"
                        },
                        {
                            "id": "2"
                        }
                    ]
                }
            }
            """
