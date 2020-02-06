Feature: It should not fail to fetch entity with invalid ID (this situation probably will not be possible to create once we support foreign keys)

    Background: We have task with invalid assigneeID
        Given I send query:
            """
            mutation {
                deleteAllUsers
                deleteAllCompanies
                deleteAllTasks
                john: createUser(input: { id: "john" }) {
                    id
                }
                jane: createUser(input: { id: "jane" }) {
                    id
                }
                createTask(
                    input: {
                    id: "test"
                    title: "do something"
                    completed: true
                    assigneeId: "jane"
                    ownerId: "john"
                    }
                ) {
                    id
                }
                deleteJohn:deleteUser(id: "john") {
                    id
                }
                deleteJane:deleteUser(id: "jane") {
                    id
                }
            }
            """

    Scenario: Fetching invalid many2many should be ok for optional fields
        When I send query:
            """
            query {
                task(id:"test"){ id title assignee { id firstName } }
            }
            """
        Then the response should be:
            """
            {
                "task": {
                    "id": "test",
                    "title": "do something",
                    "assignee": null
                }
            }
            """
        And the error should be:
            """
            null
            """

    Scenario: Fetching invalid many2many should throw error non-optional fields
        When I send query:
            """
            query {
                task(id:"test"){ id title owner { id firstName } }
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
            graphql: must not be null
            """
