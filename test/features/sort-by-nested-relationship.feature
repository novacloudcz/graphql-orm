Feature: It should be possible to fetch nested relationships with preload enabled
    Background: We have users Johny and Jane with task and company
        Given I send query:
            """
            mutation {
            deleteAllUsers
            deleteAllCompanies
            deleteAllTasks
            johny:createUser(input:{id:"johny",firstName:"John",lastName:"Doe"}) { id }
            jane:createUser(input:{id:"jane",firstName:"Jane",lastName:"Siri"}) { id }
            dummy:createUser(input:{id:"dummy",firstName:"dummy",lastName:"User"}) { id }
            c1:createCompany(input:{id:"test",name:"test company",employeesIds:["johny"]}) { id }
            c2:createCompany(input:{id:"test2",name:"AAA company",employeesIds:["jane","dummy"]}) { id }
            task1:createTask(input:{id:"test",title:"do something",completed:true,assigneeId:"jane"}) { id }
            task2:createTask(input:{id:"test2",title:"do another thing",completed:false,assigneeId:"johny"}) { id }
            task3:createTask(input:{id:"test3",title:"do another second thing",completed:false,assigneeId:"johny"}) { id }
            }
            """

    Scenario: Fetching sorted and nested preloaded relationships
        When I send query:
            """
            query {
            tasks(sort:[{assignee:{employers:{name:ASC}}},{assignee:{firstName:DESC}}]) { items { id assignee { id firstName employers { id name } } } }
            }
            """
        Then the response should be:
            """
            {
                "tasks": {
                    "items": [
                        {
                            "id": "test",
                            "assignee": {
                                "id": "jane",
                                "firstName": "Jane",
                                "employers": [
                                    {
                                        "id": "test2",
                                        "name": "AAA company"
                                    }
                                ]
                            }
                        },
                        {
                            "id": "test2",
                            "assignee": {
                                "id": "johny",
                                "firstName": "John",
                                "employers": [
                                    {
                                        "id": "test",
                                        "name": "test company"
                                    }
                                ]
                            }
                        },
                        {
                            "id": "test3",
                            "assignee": {
                                "id": "johny",
                                "firstName": "John",
                                "employers": [
                                    {
                                        "id": "test",
                                        "name": "test company"
                                    }
                                ]
                            }
                        }
                    ]
                }
            }
            """
        And the error should be empty