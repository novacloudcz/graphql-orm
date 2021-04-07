package templates

// DummyModel ...
var DummyModel = `type User @entity {
	email: String @column
	firstName: String @column
	lastName: String @column

	tasks: [Task!]! @relationship(inverse:"assignee")
}

type Task @entity {
	title: String @column
	completed: Boolean @column
	dueDate: Time @column

	assignee: User @relationship(inverse:"tasks")
}

`
