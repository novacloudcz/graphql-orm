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
                    "sdl": "scalar Time\n\ntype Query {\n  company(id: ID, q: String, filter: CompanyFilterType): Company\n  companies(offset: Int, limit: Int = 30, q: String, sort: [CompanySortType!], filter: CompanyFilterType): CompanyResultType\n  user(id: ID, q: String, filter: UserFilterType): User\n  users(offset: Int, limit: Int = 30, q: String, sort: [UserSortType!], filter: UserFilterType): UserResultType\n  task(id: ID, q: String, filter: TaskFilterType): Task\n  tasks(offset: Int, limit: Int = 30, q: String, sort: [TaskSortType!], filter: TaskFilterType): TaskResultType\n  plainEntity(id: ID, q: String, filter: PlainEntityFilterType): PlainEntity\n  plainEntities(offset: Int, limit: Int = 30, q: String, sort: [PlainEntitySortType!], filter: PlainEntityFilterType): PlainEntityResultType\n}\n\ntype Mutation {\n  createCompany(input: CompanyCreateInput!): Company!\n  updateCompany(id: ID!, input: CompanyUpdateInput!): Company!\n  deleteCompany(id: ID!): Company!\n  deleteAllCompanies: Boolean!\n  createUser(input: UserCreateInput!): User!\n  updateUser(id: ID!, input: UserUpdateInput!): User!\n  deleteUser(id: ID!): User!\n  deleteAllUsers: Boolean!\n  createTask(input: TaskCreateInput!): Task!\n  updateTask(id: ID!, input: TaskUpdateInput!): Task!\n  deleteTask(id: ID!): Task!\n  deleteAllTasks: Boolean!\n  createPlainEntity(input: PlainEntityCreateInput!): PlainEntity!\n  updatePlainEntity(id: ID!, input: PlainEntityUpdateInput!): PlainEntity!\n  deletePlainEntity(id: ID!): PlainEntity!\n  deleteAllPlainEntities: Boolean!\n}\n\nextend type Query {\n  hello: String!\n  topCompanies: [Company!]!\n}\n\ntype Company @key(fields: \"id\") {\n  id: ID!\n  name: String\n  countryId: ID\n  country: Country\n  employees: [User!]!\n  reviews: [Review!]!\n  updatedAt: Time\n  createdAt: Time!\n  updatedBy: ID\n  createdBy: ID\n  employeesIds: [ID!]!\n}\n\nextend type Company {\n  uppercaseName: String!\n}\n\ntype User {\n  id: ID!\n  email: String\n  firstName: String\n  lastName: String\n  employers: [Company!]!\n  tasks: [Task!]!\n  createdTasks: [Task!]!\n  updatedAt: Time\n  createdAt: Time!\n  updatedBy: ID\n  createdBy: ID\n  employersIds: [ID!]!\n  tasksIds: [ID!]!\n  createdTasksIds: [ID!]!\n}\n\nenum TaskState {\n  CREATED\n  IN_PROGRESS\n  RESOLVED\n}\n\ntype Task {\n  id: ID!\n  title: String\n  completed: Boolean\n  state: TaskState\n  dueDate: Time\n  assignee: User\n  owner: User!\n  assigneeId: ID\n  ownerId: ID\n  updatedAt: Time\n  createdAt: Time!\n  updatedBy: ID\n  createdBy: ID\n}\n\nextend type Review @key(fields: \"id\") {\n  id: ID! @external\n  referenceID: ID! @external\n  company: Company @requires(fields: \"referenceID\")\n}\n\nextend type Country @key(fields: \"id\") {\n  id: ID! @external\n}\n\ntype PlainEntity {\n  id: ID!\n  date: Time\n  text: String\n  updatedAt: Time\n  createdAt: Time!\n  updatedBy: ID\n  createdBy: ID\n}\n\nextend type PlainEntity {\n  shortText: String!\n}\n\nunion _Entity = Company\n\ninput CompanyCreateInput {\n  id: ID\n  name: String\n  countryId: ID\n  employeesIds: [ID!]\n}\n\ninput CompanyUpdateInput {\n  name: String\n  countryId: ID\n  employeesIds: [ID!]\n}\n\nenum CompanySortType {\n  ID_ASC\n  ID_DESC\n  NAME_ASC\n  NAME_DESC\n  COUNTRY_ID_ASC\n  COUNTRY_ID_DESC\n  UPDATED_AT_ASC\n  UPDATED_AT_DESC\n  CREATED_AT_ASC\n  CREATED_AT_DESC\n  UPDATED_BY_ASC\n  UPDATED_BY_DESC\n  CREATED_BY_ASC\n  CREATED_BY_DESC\n  EMPLOYEES_IDS_ASC\n  EMPLOYEES_IDS_DESC\n}\n\ninput CompanyFilterType {\n  AND: [CompanyFilterType!]\n  OR: [CompanyFilterType!]\n  id: ID\n  id_ne: ID\n  id_gt: ID\n  id_lt: ID\n  id_gte: ID\n  id_lte: ID\n  id_in: [ID!]\n  name: String\n  name_ne: String\n  name_gt: String\n  name_lt: String\n  name_gte: String\n  name_lte: String\n  name_in: [String!]\n  name_like: String\n  name_prefix: String\n  name_suffix: String\n  countryId: ID\n  countryId_ne: ID\n  countryId_gt: ID\n  countryId_lt: ID\n  countryId_gte: ID\n  countryId_lte: ID\n  countryId_in: [ID!]\n  updatedAt: Time\n  updatedAt_ne: Time\n  updatedAt_gt: Time\n  updatedAt_lt: Time\n  updatedAt_gte: Time\n  updatedAt_lte: Time\n  updatedAt_in: [Time!]\n  createdAt: Time\n  createdAt_ne: Time\n  createdAt_gt: Time\n  createdAt_lt: Time\n  createdAt_gte: Time\n  createdAt_lte: Time\n  createdAt_in: [Time!]\n  updatedBy: ID\n  updatedBy_ne: ID\n  updatedBy_gt: ID\n  updatedBy_lt: ID\n  updatedBy_gte: ID\n  updatedBy_lte: ID\n  updatedBy_in: [ID!]\n  createdBy: ID\n  createdBy_ne: ID\n  createdBy_gt: ID\n  createdBy_lt: ID\n  createdBy_gte: ID\n  createdBy_lte: ID\n  createdBy_in: [ID!]\n  employees: UserFilterType\n}\n\ntype CompanyResultType {\n  items: [Company!]!\n  count: Int!\n}\n\ninput UserCreateInput {\n  id: ID\n  email: String\n  firstName: String\n  lastName: String\n  employersIds: [ID!]\n  tasksIds: [ID!]\n  createdTasksIds: [ID!]\n}\n\ninput UserUpdateInput {\n  email: String\n  firstName: String\n  lastName: String\n  employersIds: [ID!]\n  tasksIds: [ID!]\n  createdTasksIds: [ID!]\n}\n\nenum UserSortType {\n  ID_ASC\n  ID_DESC\n  EMAIL_ASC\n  EMAIL_DESC\n  FIRST_NAME_ASC\n  FIRST_NAME_DESC\n  LAST_NAME_ASC\n  LAST_NAME_DESC\n  UPDATED_AT_ASC\n  UPDATED_AT_DESC\n  CREATED_AT_ASC\n  CREATED_AT_DESC\n  UPDATED_BY_ASC\n  UPDATED_BY_DESC\n  CREATED_BY_ASC\n  CREATED_BY_DESC\n  EMPLOYERS_IDS_ASC\n  EMPLOYERS_IDS_DESC\n  TASKS_IDS_ASC\n  TASKS_IDS_DESC\n  CREATED_TASKS_IDS_ASC\n  CREATED_TASKS_IDS_DESC\n}\n\ninput UserFilterType {\n  AND: [UserFilterType!]\n  OR: [UserFilterType!]\n  id: ID\n  id_ne: ID\n  id_gt: ID\n  id_lt: ID\n  id_gte: ID\n  id_lte: ID\n  id_in: [ID!]\n  email: String\n  email_ne: String\n  email_gt: String\n  email_lt: String\n  email_gte: String\n  email_lte: String\n  email_in: [String!]\n  email_like: String\n  email_prefix: String\n  email_suffix: String\n  firstName: String\n  firstName_ne: String\n  firstName_gt: String\n  firstName_lt: String\n  firstName_gte: String\n  firstName_lte: String\n  firstName_in: [String!]\n  firstName_like: String\n  firstName_prefix: String\n  firstName_suffix: String\n  lastName: String\n  lastName_ne: String\n  lastName_gt: String\n  lastName_lt: String\n  lastName_gte: String\n  lastName_lte: String\n  lastName_in: [String!]\n  lastName_like: String\n  lastName_prefix: String\n  lastName_suffix: String\n  updatedAt: Time\n  updatedAt_ne: Time\n  updatedAt_gt: Time\n  updatedAt_lt: Time\n  updatedAt_gte: Time\n  updatedAt_lte: Time\n  updatedAt_in: [Time!]\n  createdAt: Time\n  createdAt_ne: Time\n  createdAt_gt: Time\n  createdAt_lt: Time\n  createdAt_gte: Time\n  createdAt_lte: Time\n  createdAt_in: [Time!]\n  updatedBy: ID\n  updatedBy_ne: ID\n  updatedBy_gt: ID\n  updatedBy_lt: ID\n  updatedBy_gte: ID\n  updatedBy_lte: ID\n  updatedBy_in: [ID!]\n  createdBy: ID\n  createdBy_ne: ID\n  createdBy_gt: ID\n  createdBy_lt: ID\n  createdBy_gte: ID\n  createdBy_lte: ID\n  createdBy_in: [ID!]\n  employers: CompanyFilterType\n  tasks: TaskFilterType\n  createdTasks: TaskFilterType\n}\n\ntype UserResultType {\n  items: [User!]!\n  count: Int!\n}\n\ninput TaskCreateInput {\n  id: ID\n  title: String\n  completed: Boolean\n  state: TaskState\n  dueDate: Time\n  assigneeId: ID\n  ownerId: ID\n}\n\ninput TaskUpdateInput {\n  title: String\n  completed: Boolean\n  state: TaskState\n  dueDate: Time\n  assigneeId: ID\n  ownerId: ID\n}\n\nenum TaskSortType {\n  ID_ASC\n  ID_DESC\n  TITLE_ASC\n  TITLE_DESC\n  COMPLETED_ASC\n  COMPLETED_DESC\n  STATE_ASC\n  STATE_DESC\n  DUE_DATE_ASC\n  DUE_DATE_DESC\n  ASSIGNEE_ID_ASC\n  ASSIGNEE_ID_DESC\n  OWNER_ID_ASC\n  OWNER_ID_DESC\n  UPDATED_AT_ASC\n  UPDATED_AT_DESC\n  CREATED_AT_ASC\n  CREATED_AT_DESC\n  UPDATED_BY_ASC\n  UPDATED_BY_DESC\n  CREATED_BY_ASC\n  CREATED_BY_DESC\n}\n\ninput TaskFilterType {\n  AND: [TaskFilterType!]\n  OR: [TaskFilterType!]\n  id: ID\n  id_ne: ID\n  id_gt: ID\n  id_lt: ID\n  id_gte: ID\n  id_lte: ID\n  id_in: [ID!]\n  title: String\n  title_ne: String\n  title_gt: String\n  title_lt: String\n  title_gte: String\n  title_lte: String\n  title_in: [String!]\n  title_like: String\n  title_prefix: String\n  title_suffix: String\n  completed: Boolean\n  completed_ne: Boolean\n  completed_gt: Boolean\n  completed_lt: Boolean\n  completed_gte: Boolean\n  completed_lte: Boolean\n  completed_in: [Boolean!]\n  state: TaskState\n  state_ne: TaskState\n  state_gt: TaskState\n  state_lt: TaskState\n  state_gte: TaskState\n  state_lte: TaskState\n  state_in: [TaskState!]\n  dueDate: Time\n  dueDate_ne: Time\n  dueDate_gt: Time\n  dueDate_lt: Time\n  dueDate_gte: Time\n  dueDate_lte: Time\n  dueDate_in: [Time!]\n  assigneeId: ID\n  assigneeId_ne: ID\n  assigneeId_gt: ID\n  assigneeId_lt: ID\n  assigneeId_gte: ID\n  assigneeId_lte: ID\n  assigneeId_in: [ID!]\n  ownerId: ID\n  ownerId_ne: ID\n  ownerId_gt: ID\n  ownerId_lt: ID\n  ownerId_gte: ID\n  ownerId_lte: ID\n  ownerId_in: [ID!]\n  updatedAt: Time\n  updatedAt_ne: Time\n  updatedAt_gt: Time\n  updatedAt_lt: Time\n  updatedAt_gte: Time\n  updatedAt_lte: Time\n  updatedAt_in: [Time!]\n  createdAt: Time\n  createdAt_ne: Time\n  createdAt_gt: Time\n  createdAt_lt: Time\n  createdAt_gte: Time\n  createdAt_lte: Time\n  createdAt_in: [Time!]\n  updatedBy: ID\n  updatedBy_ne: ID\n  updatedBy_gt: ID\n  updatedBy_lt: ID\n  updatedBy_gte: ID\n  updatedBy_lte: ID\n  updatedBy_in: [ID!]\n  createdBy: ID\n  createdBy_ne: ID\n  createdBy_gt: ID\n  createdBy_lt: ID\n  createdBy_gte: ID\n  createdBy_lte: ID\n  createdBy_in: [ID!]\n  assignee: UserFilterType\n  owner: UserFilterType\n}\n\ntype TaskResultType {\n  items: [Task!]!\n  count: Int!\n}\n\ninput PlainEntityCreateInput {\n  id: ID\n  date: Time\n  text: String\n}\n\ninput PlainEntityUpdateInput {\n  date: Time\n  text: String\n}\n\nenum PlainEntitySortType {\n  ID_ASC\n  ID_DESC\n  DATE_ASC\n  DATE_DESC\n  TEXT_ASC\n  TEXT_DESC\n  UPDATED_AT_ASC\n  UPDATED_AT_DESC\n  CREATED_AT_ASC\n  CREATED_AT_DESC\n  UPDATED_BY_ASC\n  UPDATED_BY_DESC\n  CREATED_BY_ASC\n  CREATED_BY_DESC\n}\n\ninput PlainEntityFilterType {\n  AND: [PlainEntityFilterType!]\n  OR: [PlainEntityFilterType!]\n  id: ID\n  id_ne: ID\n  id_gt: ID\n  id_lt: ID\n  id_gte: ID\n  id_lte: ID\n  id_in: [ID!]\n  date: Time\n  date_ne: Time\n  date_gt: Time\n  date_lt: Time\n  date_gte: Time\n  date_lte: Time\n  date_in: [Time!]\n  text: String\n  text_ne: String\n  text_gt: String\n  text_lt: String\n  text_gte: String\n  text_lte: String\n  text_in: [String!]\n  text_like: String\n  text_prefix: String\n  text_suffix: String\n  updatedAt: Time\n  updatedAt_ne: Time\n  updatedAt_gt: Time\n  updatedAt_lt: Time\n  updatedAt_gte: Time\n  updatedAt_lte: Time\n  updatedAt_in: [Time!]\n  createdAt: Time\n  createdAt_ne: Time\n  createdAt_gt: Time\n  createdAt_lt: Time\n  createdAt_gte: Time\n  createdAt_lte: Time\n  createdAt_in: [Time!]\n  updatedBy: ID\n  updatedBy_ne: ID\n  updatedBy_gt: ID\n  updatedBy_lt: ID\n  updatedBy_gte: ID\n  updatedBy_lte: ID\n  updatedBy_in: [ID!]\n  createdBy: ID\n  createdBy_ne: ID\n  createdBy_gt: ID\n  createdBy_lt: ID\n  createdBy_gte: ID\n  createdBy_lte: ID\n  createdBy_in: [ID!]\n}\n\ntype PlainEntityResultType {\n  items: [PlainEntity!]!\n  count: Int!\n}"
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
