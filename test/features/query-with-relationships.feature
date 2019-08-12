Feature: It should be possible to fetch items by relationships

    Background: We have users Johny and Jane with company
        Given I send query:
            """
            mutation {
            deleteAllUsers
            deleteAllCompanies
            johny:createUser(input:{id:"johny",firstName:"John",lastName:"Doe"}) { id }
            jane:createUser(input:{id:"jane",firstName:"Jane",lastName:"Siri"}) { id }
            createCompany(input:{id:"test",name:"test company",employeesIds:["johny"]}) { id }
            }
            """

    Scenario: Fetching users by company name
        When I send query:
            """
            query {
            users(q:"test") { items { firstName lastName employers { name } } count }
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

    Scenario: Fetching users by company name with multiple words
        When I send query:
            """
            query {
            users(q:"test company") { items { firstName lastName employers { name } } count }
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
