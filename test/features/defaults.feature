Feature: It should possible to specify default value

    Background: We have test company
        Given I send query:
            """
            mutation {
                deleteAllTasks
                test:createTask(input:{id:"test",title:"Test"}) { id }
            }
            """

    Scenario: Fetching single item
        When I send query:
            """
            query {
                task(id: "test") {
                    id
                    completed
                }
            }
            """
        Then the response should be:
            """
            {
                "task": {
                    "id": "test",
                    "completed": false
                }
            }
            """
