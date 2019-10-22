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

    Scenario: Fetching _service
        When I send query:
            """
            query { _service { sdl } }
            """
        Then the response should be:
            """
            {
                "_service": {
                    "sdl": "scalar Time\n\ntype Query {\n  company(id: ID, q: String, filter: CompanyFilterType): Company\n  companies(offset: Int, limit: Int = 30, q: String, sort: [CompanySortType!], filter: CompanyFilterType): CompanyResultType\n  user(id: ID, q: String, filter: UserFilterType): User\n  users(offset: Int, limit: Int = 30, q: String, sort: [UserSortType!], filter: UserFilterType): UserResultType\n  task(id: ID, q: String, filter: TaskFilterType): Task\n  tasks(offset: Int, limit: Int = 30, q: String, sort: [TaskSortType!], filter: TaskFilterType): TaskResultType\n  plainEntity(id: ID, q: String, filter: PlainEntityFilterType): PlainEntity\n  plainEntities(offset: Int, limit: Int = 30, q: String, sort: [PlainEntitySortType!], filter: PlainEntityFilterType): PlainEntityResultType\n}\n\ntype Mutation {\n  createCompany(input: CompanyCreateInput!): Company!\n  updateCompany(id: ID!, input: CompanyUpdateInput!): Company!\n  deleteCompany(id: ID!): Company!\n  deleteAllCompanies: Boolean!\n  createUser(input: UserCreateInput!): User!\n  updateUser(id: ID!, input: UserUpdateInput!): User!\n  deleteUser(id: ID!): User!\n  deleteAllUsers: Boolean!\n  createTask(input: TaskCreateInput!): Task!\n  updateTask(id: ID!, input: TaskUpdateInput!): Task!\n  deleteTask(id: ID!): Task!\n  deleteAllTasks: Boolean!\n  createPlainEntity(input: PlainEntityCreateInput!): PlainEntity!\n  updatePlainEntity(id: ID!, input: PlainEntityUpdateInput!): PlainEntity!\n  deletePlainEntity(id: ID!): PlainEntity!\n  deleteAllPlainEntities: Boolean!\n}\n\nenum ObjectSortType {\n  ASC\n  DESC\n}\n\nextend type Query {\n  hello: String!\n  topCompanies: [Company!]!\n}\n\ntype Company @key(fields: \"id\") {\n  id: ID!\n  name: String\n  countryId: ID\n  country: Country\n  employees: [User!]!\n  reviews: [Review!]!\n  updatedAt: Time\n  createdAt: Time!\n  updatedBy: ID\n  createdBy: ID\n  employeesIds: [ID!]!\n}\n\nextend type Company {\n  uppercaseName: String!\n}\n\ntype Address {\n  street: String\n  city: String\n  zip: String\n}\n\ntype User {\n  id: ID!\n  code: Int\n  email: String\n  firstName: String\n  lastName: String\n  addressRaw: String\n  address: Address\n  employers: [Company!]!\n  tasks: [Task!]!\n  createdTasks: [Task!]!\n  updatedAt: Time\n  createdAt: Time!\n  updatedBy: ID\n  createdBy: ID\n  employersIds: [ID!]!\n  tasksIds: [ID!]!\n  createdTasksIds: [ID!]!\n}\n\nenum TaskState {\n  CREATED\n  IN_PROGRESS\n  RESOLVED\n}\n\ntype Task {\n  id: ID!\n  title: String\n  completed: Boolean\n  state: TaskState\n  dueDate: Time\n  assignee: User\n  owner: User!\n  assigneeId: ID\n  ownerId: ID\n  updatedAt: Time\n  createdAt: Time!\n  updatedBy: ID\n  createdBy: ID\n}\n\nextend type Review @entity @key(fields: \"id\") {\n  id: ID! @external\n  referenceID: ID! @external\n  company: Company @requires(fields: \"referenceID\")\n}\n\nextend type Country @entity @key(fields: \"id\") {\n  id: ID! @external\n}\n\ntype PlainEntity {\n  id: ID!\n  date: Time\n  text: String\n  shortText: String!\n  updatedAt: Time\n  createdAt: Time!\n  updatedBy: ID\n  createdBy: ID\n}\n\nunion _Entity = Company\n\ninput CompanyCreateInput {\n  id: ID\n  name: String\n  countryId: ID\n  employeesIds: [ID!]\n}\n\ninput CompanyUpdateInput {\n  name: String\n  countryId: ID\n  employeesIds: [ID!]\n}\n\ninput CompanySortType {\n  id: ObjectSortType\n  name: ObjectSortType\n  countryId: ObjectSortType\n  updatedAt: ObjectSortType\n  createdAt: ObjectSortType\n  updatedBy: ObjectSortType\n  createdBy: ObjectSortType\n  employeesIds: ObjectSortType\n  employees: UserSortType\n}\n\ninput CompanyFilterType {\n  AND: [CompanyFilterType!]\n  OR: [CompanyFilterType!]\n  id: ID\n  id_ne: ID\n  id_gt: ID\n  id_lt: ID\n  id_gte: ID\n  id_lte: ID\n  id_in: [ID!]\n  id_null: Boolean\n  name: String\n  name_ne: String\n  name_gt: String\n  name_lt: String\n  name_gte: String\n  name_lte: String\n  name_in: [String!]\n  name_like: String\n  name_prefix: String\n  name_suffix: String\n  name_null: Boolean\n  countryId: ID\n  countryId_ne: ID\n  countryId_gt: ID\n  countryId_lt: ID\n  countryId_gte: ID\n  countryId_lte: ID\n  countryId_in: [ID!]\n  countryId_null: Boolean\n  updatedAt: Time\n  updatedAt_ne: Time\n  updatedAt_gt: Time\n  updatedAt_lt: Time\n  updatedAt_gte: Time\n  updatedAt_lte: Time\n  updatedAt_in: [Time!]\n  updatedAt_null: Boolean\n  createdAt: Time\n  createdAt_ne: Time\n  createdAt_gt: Time\n  createdAt_lt: Time\n  createdAt_gte: Time\n  createdAt_lte: Time\n  createdAt_in: [Time!]\n  createdAt_null: Boolean\n  updatedBy: ID\n  updatedBy_ne: ID\n  updatedBy_gt: ID\n  updatedBy_lt: ID\n  updatedBy_gte: ID\n  updatedBy_lte: ID\n  updatedBy_in: [ID!]\n  updatedBy_null: Boolean\n  createdBy: ID\n  createdBy_ne: ID\n  createdBy_gt: ID\n  createdBy_lt: ID\n  createdBy_gte: ID\n  createdBy_lte: ID\n  createdBy_in: [ID!]\n  createdBy_null: Boolean\n  employees: UserFilterType\n}\n\ntype CompanyResultType {\n  items: [Company!]!\n  count: Int!\n}\n\ninput UserCreateInput {\n  id: ID\n  code: Int\n  email: String\n  firstName: String\n  lastName: String\n  addressRaw: String\n  employersIds: [ID!]\n  tasksIds: [ID!]\n  createdTasksIds: [ID!]\n}\n\ninput UserUpdateInput {\n  code: Int\n  email: String\n  firstName: String\n  lastName: String\n  addressRaw: String\n  employersIds: [ID!]\n  tasksIds: [ID!]\n  createdTasksIds: [ID!]\n}\n\ninput UserSortType {\n  id: ObjectSortType\n  code: ObjectSortType\n  email: ObjectSortType\n  firstName: ObjectSortType\n  lastName: ObjectSortType\n  addressRaw: ObjectSortType\n  updatedAt: ObjectSortType\n  createdAt: ObjectSortType\n  updatedBy: ObjectSortType\n  createdBy: ObjectSortType\n  employersIds: ObjectSortType\n  tasksIds: ObjectSortType\n  createdTasksIds: ObjectSortType\n  employers: CompanySortType\n  tasks: TaskSortType\n  createdTasks: TaskSortType\n}\n\ninput UserFilterType {\n  AND: [UserFilterType!]\n  OR: [UserFilterType!]\n  id: ID\n  id_ne: ID\n  id_gt: ID\n  id_lt: ID\n  id_gte: ID\n  id_lte: ID\n  id_in: [ID!]\n  id_null: Boolean\n  code: Int\n  code_ne: Int\n  code_gt: Int\n  code_lt: Int\n  code_gte: Int\n  code_lte: Int\n  code_in: [Int!]\n  code_null: Boolean\n  email: String\n  email_ne: String\n  email_gt: String\n  email_lt: String\n  email_gte: String\n  email_lte: String\n  email_in: [String!]\n  email_like: String\n  email_prefix: String\n  email_suffix: String\n  email_null: Boolean\n  firstName: String\n  firstName_ne: String\n  firstName_gt: String\n  firstName_lt: String\n  firstName_gte: String\n  firstName_lte: String\n  firstName_in: [String!]\n  firstName_like: String\n  firstName_prefix: String\n  firstName_suffix: String\n  firstName_null: Boolean\n  lastName: String\n  lastName_ne: String\n  lastName_gt: String\n  lastName_lt: String\n  lastName_gte: String\n  lastName_lte: String\n  lastName_in: [String!]\n  lastName_like: String\n  lastName_prefix: String\n  lastName_suffix: String\n  lastName_null: Boolean\n  addressRaw: String\n  addressRaw_ne: String\n  addressRaw_gt: String\n  addressRaw_lt: String\n  addressRaw_gte: String\n  addressRaw_lte: String\n  addressRaw_in: [String!]\n  addressRaw_like: String\n  addressRaw_prefix: String\n  addressRaw_suffix: String\n  addressRaw_null: Boolean\n  updatedAt: Time\n  updatedAt_ne: Time\n  updatedAt_gt: Time\n  updatedAt_lt: Time\n  updatedAt_gte: Time\n  updatedAt_lte: Time\n  updatedAt_in: [Time!]\n  updatedAt_null: Boolean\n  createdAt: Time\n  createdAt_ne: Time\n  createdAt_gt: Time\n  createdAt_lt: Time\n  createdAt_gte: Time\n  createdAt_lte: Time\n  createdAt_in: [Time!]\n  createdAt_null: Boolean\n  updatedBy: ID\n  updatedBy_ne: ID\n  updatedBy_gt: ID\n  updatedBy_lt: ID\n  updatedBy_gte: ID\n  updatedBy_lte: ID\n  updatedBy_in: [ID!]\n  updatedBy_null: Boolean\n  createdBy: ID\n  createdBy_ne: ID\n  createdBy_gt: ID\n  createdBy_lt: ID\n  createdBy_gte: ID\n  createdBy_lte: ID\n  createdBy_in: [ID!]\n  createdBy_null: Boolean\n  employers: CompanyFilterType\n  tasks: TaskFilterType\n  createdTasks: TaskFilterType\n}\n\ntype UserResultType {\n  items: [User!]!\n  count: Int!\n}\n\ninput TaskCreateInput {\n  id: ID\n  title: String\n  completed: Boolean\n  state: TaskState\n  dueDate: Time\n  assigneeId: ID\n  ownerId: ID\n}\n\ninput TaskUpdateInput {\n  title: String\n  completed: Boolean\n  state: TaskState\n  dueDate: Time\n  assigneeId: ID\n  ownerId: ID\n}\n\ninput TaskSortType {\n  id: ObjectSortType\n  title: ObjectSortType\n  completed: ObjectSortType\n  state: ObjectSortType\n  dueDate: ObjectSortType\n  assigneeId: ObjectSortType\n  ownerId: ObjectSortType\n  updatedAt: ObjectSortType\n  createdAt: ObjectSortType\n  updatedBy: ObjectSortType\n  createdBy: ObjectSortType\n  assignee: UserSortType\n  owner: UserSortType\n}\n\ninput TaskFilterType {\n  AND: [TaskFilterType!]\n  OR: [TaskFilterType!]\n  id: ID\n  id_ne: ID\n  id_gt: ID\n  id_lt: ID\n  id_gte: ID\n  id_lte: ID\n  id_in: [ID!]\n  id_null: Boolean\n  title: String\n  title_ne: String\n  title_gt: String\n  title_lt: String\n  title_gte: String\n  title_lte: String\n  title_in: [String!]\n  title_like: String\n  title_prefix: String\n  title_suffix: String\n  title_null: Boolean\n  completed: Boolean\n  completed_ne: Boolean\n  completed_gt: Boolean\n  completed_lt: Boolean\n  completed_gte: Boolean\n  completed_lte: Boolean\n  completed_in: [Boolean!]\n  completed_null: Boolean\n  state: TaskState\n  state_ne: TaskState\n  state_gt: TaskState\n  state_lt: TaskState\n  state_gte: TaskState\n  state_lte: TaskState\n  state_in: [TaskState!]\n  state_null: Boolean\n  dueDate: Time\n  dueDate_ne: Time\n  dueDate_gt: Time\n  dueDate_lt: Time\n  dueDate_gte: Time\n  dueDate_lte: Time\n  dueDate_in: [Time!]\n  dueDate_null: Boolean\n  assigneeId: ID\n  assigneeId_ne: ID\n  assigneeId_gt: ID\n  assigneeId_lt: ID\n  assigneeId_gte: ID\n  assigneeId_lte: ID\n  assigneeId_in: [ID!]\n  assigneeId_null: Boolean\n  ownerId: ID\n  ownerId_ne: ID\n  ownerId_gt: ID\n  ownerId_lt: ID\n  ownerId_gte: ID\n  ownerId_lte: ID\n  ownerId_in: [ID!]\n  ownerId_null: Boolean\n  updatedAt: Time\n  updatedAt_ne: Time\n  updatedAt_gt: Time\n  updatedAt_lt: Time\n  updatedAt_gte: Time\n  updatedAt_lte: Time\n  updatedAt_in: [Time!]\n  updatedAt_null: Boolean\n  createdAt: Time\n  createdAt_ne: Time\n  createdAt_gt: Time\n  createdAt_lt: Time\n  createdAt_gte: Time\n  createdAt_lte: Time\n  createdAt_in: [Time!]\n  createdAt_null: Boolean\n  updatedBy: ID\n  updatedBy_ne: ID\n  updatedBy_gt: ID\n  updatedBy_lt: ID\n  updatedBy_gte: ID\n  updatedBy_lte: ID\n  updatedBy_in: [ID!]\n  updatedBy_null: Boolean\n  createdBy: ID\n  createdBy_ne: ID\n  createdBy_gt: ID\n  createdBy_lt: ID\n  createdBy_gte: ID\n  createdBy_lte: ID\n  createdBy_in: [ID!]\n  createdBy_null: Boolean\n  assignee: UserFilterType\n  owner: UserFilterType\n}\n\ntype TaskResultType {\n  items: [Task!]!\n  count: Int!\n}\n\ninput PlainEntityCreateInput {\n  id: ID\n  date: Time\n  text: String\n}\n\ninput PlainEntityUpdateInput {\n  date: Time\n  text: String\n}\n\ninput PlainEntitySortType {\n  id: ObjectSortType\n  date: ObjectSortType\n  text: ObjectSortType\n  updatedAt: ObjectSortType\n  createdAt: ObjectSortType\n  updatedBy: ObjectSortType\n  createdBy: ObjectSortType\n}\n\ninput PlainEntityFilterType {\n  AND: [PlainEntityFilterType!]\n  OR: [PlainEntityFilterType!]\n  id: ID\n  id_ne: ID\n  id_gt: ID\n  id_lt: ID\n  id_gte: ID\n  id_lte: ID\n  id_in: [ID!]\n  id_null: Boolean\n  date: Time\n  date_ne: Time\n  date_gt: Time\n  date_lt: Time\n  date_gte: Time\n  date_lte: Time\n  date_in: [Time!]\n  date_null: Boolean\n  text: String\n  text_ne: String\n  text_gt: String\n  text_lt: String\n  text_gte: String\n  text_lte: String\n  text_in: [String!]\n  text_like: String\n  text_prefix: String\n  text_suffix: String\n  text_null: Boolean\n  updatedAt: Time\n  updatedAt_ne: Time\n  updatedAt_gt: Time\n  updatedAt_lt: Time\n  updatedAt_gte: Time\n  updatedAt_lte: Time\n  updatedAt_in: [Time!]\n  updatedAt_null: Boolean\n  createdAt: Time\n  createdAt_ne: Time\n  createdAt_gt: Time\n  createdAt_lt: Time\n  createdAt_gte: Time\n  createdAt_lte: Time\n  createdAt_in: [Time!]\n  createdAt_null: Boolean\n  updatedBy: ID\n  updatedBy_ne: ID\n  updatedBy_gt: ID\n  updatedBy_lt: ID\n  updatedBy_gte: ID\n  updatedBy_lte: ID\n  updatedBy_in: [ID!]\n  updatedBy_null: Boolean\n  createdBy: ID\n  createdBy_ne: ID\n  createdBy_gt: ID\n  createdBy_lt: ID\n  createdBy_gte: ID\n  createdBy_lte: ID\n  createdBy_in: [ID!]\n  createdBy_null: Boolean\n}\n\ntype PlainEntityResultType {\n  items: [PlainEntity!]!\n  count: Int!\n}"
                }
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
