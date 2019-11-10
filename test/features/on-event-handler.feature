Feature: OnEvent handler should be called for mutation

    Background: We have test data
        Given I send query:
            """
            mutation {
            deleteAllUsers
            deleteAllTasks
            createUser(input:{id:"johny",firstName:"John",lastName:"Doe",tasksIds:["test"]}) { id firstName lastName tasks { id } }
            }
            """

    Scenario: Changing user firstName should create new task for user
        When I send query:
            """
            mutation {
            updateUser(id:"johny",input:{firstName:"Johny"}) { id firstName }
            }
            """
        Then the error should be:
            """
            null
            """
        Then the response should be:
            """
            {
                "updateUser": {
                    "id": "johny",
                    "firstName": "Johny"
                }
            }
            """
        When I send query:
            """
            query {
            user(id:"johny") { tasks { title } }
            }
            """
        Then the response should be:
            """
            {
                "user": {
                    "tasks": [
                        {
                            "title": "Hello Johny!"
                        }
                    ]
                }
            }
            """
