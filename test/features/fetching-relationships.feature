Feature: It should be possible to fetch with relationships

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
            company2:createCompany(input:{id:"test2",name:"test2 company",employeesIds:["johny","jane"]}) { id }
            createTask(input:{id:"test",title:"do something",completed:true,assigneeId:"jane"}) { id }
            }
            """

    Scenario: Fetching many2many
        When I send query:
            """
            query {
            users(sort:[{firstName:DESC}]) { items { firstName lastName createdBy updatedBy employers { name } } count }
            companies(filter:{id:"test"}) { items { name employees { firstName } } }
            }
            """
        Then the response should be:
            """
            {
                "companies": {
                    "items": [
                        {
                            "employees": [
                                {
                                    "firstName": "John"
                                }
                            ],
                            "name": "test company"
                        }
                    ]
                },
                "users": {
                    "count": 2,
                    "items": [
                        {
                            "createdBy": null,
                            "employers": [
                                {
                                    "name": "test company"
                                },
                                {
                                    "name": "test2 company"
                                }
                            ],
                            "firstName": "John",
                            "lastName": "Doe",
                            "updatedBy": null
                        },
                        {
                            "createdBy": null,
                            "employers": [
                                {
                                    "name": "test2 company"
                                }
                            ],
                            "firstName": "Jane",
                            "lastName": "Siri",
                            "updatedBy": null
                        }
                    ]
                }
            }
            """

    Scenario: Fetching many2many with filter
        When I send query:
            """
            query {
            users(filter:{employers:{id:"test"}}) { items { firstName lastName createdBy updatedBy employers { name } } count }
            }
            """
        Then the response should be:
            """
            {
                "users": {
                    "items": [
                        {
                            "firstName": "John",
                            "lastName": "Doe",
                            "createdBy": null,
                            "updatedBy": null,
                            "employers": [
                                {
                                    "name": "test company"
                                },
                                {
                                    "name": "test2 company"
                                }
                            ]
                        }
                    ],
                    "count": 1
                }
            }
            """

    Scenario: Fetching many2many with filter by id
        When I send query:
            """
            query {
            users(filter:{employersIds_in:["test"]}) { items { firstName lastName createdBy updatedBy employers { name } } count }
            }
            """
        Then the response should be:
            """
            {
                "users": {
                    "items": [
                        {
                            "firstName": "John",
                            "lastName": "Doe",
                            "createdBy": null,
                            "updatedBy": null,
                            "employers": [
                                {
                                    "name": "test company"
                                },
                                {
                                    "name": "test2 company"
                                }
                            ]
                        }
                    ],
                    "count": 1
                }
            }
            """

    Scenario: Fetching toOne
        When I send query:
            """
            query {
            user(id:"jane") { firstName lastName tasks { id title } }
            }
            """
        Then the response should be:
            """
            {
                "user": {
                    "firstName": "Jane",
                    "lastName": "Siri",
                    "tasks": [
                        {
                            "id": "test",
                            "title": "do something"
                        }
                    ]
                }
            }
            """


    Scenario: Fetching toOne with filter
        When I send query:
            """
            query {
            user(filter:{tasks:{id:"test"}}) { firstName lastName tasks { id title } }
            }
            """
        Then the response should be:
            """
            {
                "user": {
                    "firstName": "Jane",
                    "lastName": "Siri",
                    "tasks": [
                        {
                            "id": "test",
                            "title": "do something"
                        }
                    ]
                }
            }
            """


    Scenario: Fetching toMany with filter
        When I send query:
            """
            query {
            task(filter:{assigneeId:"jane"}) { title completed assignee { id firstName } }
            }
            """
        Then the response should be:
            """
            {
                "task": {
                    "title": "do something",
                    "completed": true,
                    "assignee": {
                        "id": "jane",
                        "firstName": "Jane"
                    }
                }
            }
            """
