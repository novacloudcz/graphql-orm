Feature: It should be possible to fetch relationship using connection

    Background: We have Johny and some tasks
        Given I send query:
            """
            mutation {
                deleteAllUsers
                deleteAllTasks
                createUser(input: { id: "john" }) { id }
                t1: createTask(input: { id:"t1",assigneeId:"john" }) { id }
                t2: createTask(input: { id:"t2",assigneeId:"john" }) { id }
            }
            """

    Scenario: Fetching user with tasks should be possible using connection
        When I send query:
            """
            query {
                users {
                    items {
                        tasksConnection(limit: 1) {
                            items {
                                id
                            }
                            count
                        }
                    }
                    count
                }
            }
            """
        Then the error should be empty
        And the response should be:
            """
            {
                "users": {
                    "items": [
                        {
                            "tasksConnection": {
                                "items": [
                                    {
                                        "id": "t1"
                                    }
                                ],
                                "count": 2
                            }
                        }
                    ],
                    "count": 1
                }
            }
            """