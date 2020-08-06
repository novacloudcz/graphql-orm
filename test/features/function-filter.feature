Feature: It should be possible to filter out based on aggregation functions

    Background: We have users with salary
        Given I send query:
            """
            mutation {
            deleteAllUsers
            deleteAllCompanies
            john: createUser(input: { id: "john.doe", salary: 1500 }) { id }
            john2: createUser(input: { id: "john.doe2", salary: 500 }) { id }
            jane: createUser(input: { id: "jane.doe", salary: 2500 }) { id }
            e1: createCompany(input:{id:"e1",employeesIds:["john.doe","john.doe2"]}){id}
            e2: createCompany(input:{id:"e2",employeesIds:["jane.doe"]}){id}
            }
            """

    Scenario: Creating user with existing email should fail
        When I send query:
            """
            query {
                companies(filter:{employees:{salaryMin_gt: 400,salaryMax_lt: 1300}}) {
                    items {
                        id
                        employees {
                            salary
                        }
                    }
                }
            }
            """
        Then the error should be empty
        And the response should be:
            """
            {
                "companies": {
                    "items": []
                }
            }
            """