Feature: It should be possible to create/update entities with relations

    Background: We have test task
        Given I send query:
            """
            mutation {
            deleteAllUsers
            createTask(input:{id:"test",title:"do something",completed:true}) { id }
            }
            """

    Scenario: Creating user with already existing task
        When I send query:
            """
            mutation {
            createUser(input:{id:"johny",firstName:"John",lastName:"Doe",tasksIds:["test"]}) { id firstName lastName tasks { id } }
            }
            """
        Then the error should be:
            """
            null
            """
        Then the response should be:
            """
            {
                "createUser": {
                    "id": "johny",
                    "firstName": "John",
                    "lastName": "Doe",
                    "tasks": [
                        {
                            "id": "test"
                        }
                    ]
                }
            }
            """
