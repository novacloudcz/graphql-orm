Feature: It should be possible to mutate with relationships

    Background: We have users Johny and Jane
        Given I send query:
            """
            mutation {
            deleteAllUsers
            deleteAllCompanies
            deleteAllTasks
            johny:createUser(input:{id:"johny",firstName:"John",lastName:"Doe"}) { id }
            jane:createUser(input:{id:"jane",firstName:"Jane",lastName:"Siri"}) { id }
            createCompany(input:{id:"test",name:"test company",employeesIds:["johny"]}) { id }
            createTask(input:{id:"test",title:"do something",completed:true,assigneeId:"jane"}) { id }
            }
            """

    Scenario: Fetching many2many
        When I send query:
            """
            query {
            users(sort:[FIRST_NAME_DESC]) { items { firstName lastName createdBy updatedBy employers { name } } count }
            companies { items { name employees { firstName } } }
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
                                }
                            ]
                        },
                        {
                            "firstName": "Jane",
                            "lastName": "Siri",
                            "createdBy": null,
                            "updatedBy": null,
                            "employers": []
                        }
                    ],
                    "count": 2
                },
                "companies": {
                    "items": [
                        {
                            "name": "test company",
                            "employees": [
                                {
                                    "firstName": "John"
                                }
                            ]
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
