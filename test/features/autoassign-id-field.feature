Feature: Field xxxId should be automatically used for filling object for xxx field

    Background: We have test company
        Given I send query:
            """
            mutation {
            deleteAllCompanies
            deleteAllTasks
            createCompany(input:{id:"test",name:"test company",countryId:"xxx"}) { id }
            createTask(input:{id:"test",assigneeId:"xxx",ownerId:"xxx"}) { id }
            }
            """

    Scenario: Fetching country should use the countryId field as id
        When I send query:
            """
            query {
            company(id:"test") { id countryId country { id } }
            }
            """
        Then the response should be:
            """
            {
                "company": {
                    "id": "test",
                    "countryId": "xxx",
                    "country": {
                        "id": "xxx"
                    }
                }
            }
            """
    Scenario: Fetching assignee should fail due to resolver not being implemented
        When I send query:
            """
            query {
            task(id:"test") { id assigneeId assignee { id } }
            }
            """
        Then the response should be:
            """
            {
                "task": {
                    "id": "test",
                    "assigneeId": "xxx",
                    "assignee": null
                }
            }
            """
        And the error should be:
            """
            null
            """
    Scenario: Fetching owner (required field) should fail with error due to resolver not being implemented
        When I send query:
            """
            query {
            task(id:"test") { id ownerId owner { id } }
            }
            """
        Then the response should be:
            """
            {
                "task": null
            }
            """
        And the error should be:
            """
            graphql: User with id 'xxx' not found
            """