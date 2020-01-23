Feature: It should be possible to create/update/fetch embedded types with column persistance

    Background: We have test task
        Given I send query:
            """
            mutation {
            deleteAllTasks
            createTask(input:{id:"test",meta:{key:"hello",value:"world"},metas:[{key:"john",value:"doe"}]}) { id }
            }
            """

    Scenario: Fetching task with embedded types
        When I send query:
            """
            query {
            task(id:"test") { id meta { key value } metas { key value } }
            }
            """
        Then the response should be:
            """
            {
                "task": {
                    "id": "test",
                    "meta": {
                        "key": "hello",
                        "value": "world"
                    },
                    "metas": [
                        {
                            "key": "john",
                            "value": "doe"
                        }
                    ]
                }
            }
            """

    Scenario: Creating task with empty embedded columns
        When I send query:
            """
            mutation {
            createTask(input:{id:"test2",meta:null,metas:[]}) { id meta { key value } metas { key value } }
            }
            """
        Then the response should be:
            """
            {
                "createTask": {
                    "id": "test2",
                    "meta": null,
                    "metas": []
                }
            }
            """
    Scenario: Updating task with embedded columns
        When I send query:
            """
            mutation {
            updateTask(id:"test",input:{meta:{key:"hello2",value:"world2"},metas:[]}) { id meta { key value } metas { key value } }
            }
            """
        Then the response should be:
            """
            {
                "updateTask": {
                    "id": "test",
                    "meta": {
                        "key": "hello2",
                        "value": "world2"
                    },
                    "metas": []
                }
            }
            """