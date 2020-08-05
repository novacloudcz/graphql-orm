Feature: Pagination should return correct data set

    Background: We have multiple users and tasks
        Given I send query:
            """
            mutation {
            deleteAllUsers
            deleteAllTasks
            u1:createUser(input: { id:"u1" }) { id }
            u2:createUser(input: { id:"u2" }) { id }
            t1:createTask(input: { id:"t1",title:"test", assigneeId:"u1"}) { id }
            t2:createTask(input: { id:"t2",title:"test", assigneeId:"u1"}) { id }
            t3:createTask(input: { id:"t3",title:"test", assigneeId:"u2"}) { id }
            }
            """

    Scenario: Fetching users joined with task should return proper item count
        When I send query:
            """
            query {
                users(filter: { tasks: { title: "test" } }) {
                    items {
                        id
                    }
                    count
                }
            }
            """
        Then the response should be:
            """
            {
                "users": {
                    "items": [
                        {
                            "id": "u1"
                        },
                        {
                            "id": "u2"
                        }
                    ],
                    "count": 2
                }
            }
            """