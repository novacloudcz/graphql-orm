Feature: It should not be possible to create multiple users with same email

    Background: We have Johny
        Given I send query:
            """
            mutation {
            deleteAllUsers
            createUser(input:{email:"john.doe@example.com"}) { id }
            }
            """

    Scenario: Creating user with existing email should fail
        When I send query:
            """
            mutation {
            createUser(input:{email:"john.doe@example.com"}) { id }
            }
            """
        Then the error should not be empty