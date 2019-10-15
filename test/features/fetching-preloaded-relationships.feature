Feature: It should be possible to mutate with relationships

    Background: We have users Johny and Jane with task and company
        Given I send query:
            """
            mutation {
            deleteAllUsers
            deleteAllCompanies
            deleteAllTasks
            johny:createUser(input:{id:"johny",firstName:"John",lastName:"Doe"}) { id }
            jane:createUser(input:{id:"jane",firstName:"Jane",lastName:"Siri"}) { id }
            createCompany(input:{id:"test",name:"test company",employeesIds:["johny"]}) { id }
            task1:createTask(input:{id:"test",title:"do something",completed:true,assigneeId:"jane"}) { id }
            task2:createTask(input:{id:"test2",title:"do another thing",completed:true,assigneeId:"johny"}) { id }
            }
            """

    Scenario: Fetching many2many
        When I send query:
            """
            query {
            users { items { id lastName firstName employers { id name } } }
            }
            """
        Then the response should be:
            """
            {
                "users": {
                    "items": [
                        {
                            "id": "johny",
                            "lastName": "Doe",
                            "firstName": "John",
                            "employers": [
                                {
                                    "id": "test",
                                    "name": "test company"
                                }
                            ]
                        },
                        {
                            "id": "jane",
                            "lastName": "Siri",
                            "firstName": "Jane",
                            "employers": []
                        }
                    ]
                }
            }
            """

    Scenario: Fetching with nested preloaded relationships
        When I send query:
            """
            query {
            tasks { items {id assignee { id firstName employers { id name } } } }
            }
            """
        Then the response should be:
            """
            {
                "tasks": {
                    "items": [
                        {
                            "assignee": {
                                "employers": [],
                                "firstName": "Jane",
                                "id": "jane"
                            },
                            "id": "test"
                        },
                        {
                            "assignee": {
                                "employers": [
                                    {
                                        "id": "test",
                                        "name": "test company"
                                    }
                                ],
                                "firstName": "John",
                                "id": "johny"
                            },
                            "id": "test2"
                        }
                    ]
                }
            }
            """
