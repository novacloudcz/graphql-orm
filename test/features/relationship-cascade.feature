Feature: It should delete all objects with cascade relationship
    Background: We have nested tasks and delete them afterwards
        Given I send query:
            """
            mutation {
                deleteAllTasks
                t1:createTask(input:{id:"t1"}) { id }
                t2:createTask(input:{id:"t2",parentTaskId:"t1"}){id}
                t3:createTask(input:{id:"t3",parentTaskId:"t2"}){id}
                deleteTask(id:"t1") {id}
            }
            """

    Scenario: Fetching tasks should return empty list
        When I send query:
            """
            query {
                tasks {
                    items {
                    id
                    subtasksIds
                    }
                }
            }
            """
        Then the response should be:
            """
            {
                "tasks": {
                    "items": []
                }
            }
            """
