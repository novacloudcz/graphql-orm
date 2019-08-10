Feature: It should not fail to fetch entity with invalid ID (this situation probably will not be possible to create once we support foreign keys)

    Background: We have task with invalid assigneeID
        Given I send query:
            """
            mutation {
            deleteAllUsers
            deleteAllCompanies
            deleteAllTasks
            createTask(input:{id:"test",title:"do something",completed:true,assigneeId:"jane"}) { id }
            }
            """

    Scenario: Fetching many2many
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
