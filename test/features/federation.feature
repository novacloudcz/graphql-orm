Feature: It should be possible fetch fields from apollo federation specs
    Background: We have test company
        Given I send query:
            """
            mutation {
            deleteAllCompanies
            test:createCompany(input:{id:"test",name:"Test company"}) { id }
            test2:createCompany(input:{id:"test2",name:"Test2 company"}) { id }
            }
            """

    Scenario: Fetching _entities with empty representations should return null
        When I send query:
            """
            query { _entities(representations:[{__typename:"Company"}]) { __typename } }
            """
        Then the response should be:
            """
            {
                "_entities": [
                    null
                ]
            }
            """
    Scenario: Fetching _entities with resolving reference
        When I send query:
            """
            query { _entities(representations:[{__typename:"Company",id:"test"},{__typename:"Company",id:"test2"}]) {
            __typename
            ... on Company { id name }
            } }
            """
        Then the response should be:
            """
            {
                "_entities": [
                    {
                        "__typename": "Company",
                        "id": "test",
                        "name": "Test company"
                    },
                    {
                        "__typename": "Company",
                        "id": "test2",
                        "name": "Test2 company"
                    }
                ]
            }
            """
    Scenario: Fetching _entities by non ID field with resolving reference
        When I send query:
            """
            query { _entities(representations:[{__typename:"Company",name:"Test company"},{__typename:"Company",name:"Test2 company"}]) {
            __typename
            ... on Company { id name }
            } }
            """
        Then the response should be:
            """
            {
                "_entities": [
                    {
                        "__typename": "Company",
                        "id": "test",
                        "name": "Test company"
                    },
                    {
                        "__typename": "Company",
                        "id": "test2",
                        "name": "Test2 company"
                    }
                ]
            }
            """
    Scenario: Fetching _entities by non existing fields with resolving reference
        When I send query:
            """
            query { _entities(representations:[{__typename:"Company",blah:"xx"},{__typename:"Company",foo:"xx"}]) {
            __typename
            ... on Company { id name }
            } }
            """
        Then the response should be:
            """
            {
                "_entities": [
                    null,
                    null
                ]
            }
            """
    Scenario: Fetching _entities by nonexisting ID field with resolving reference
        When I send query:
            """
            query { _entities(representations:[{__typename:"Company",id:"aaa"}]) {
            __typename
            ... on Company { id name }
            } }
            """
        Then the response should be:
            """
            {
                "_entities": [
                    null
                ]
            }
            """
        And the error should be:
            """
            null
            """
