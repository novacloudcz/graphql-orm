package templates

var DummyModel = `type User {
	email: String
	firstName: String
	lastName: String

	tasks: [Task!]! @relationship(inverse:"assignee")
}

type Task {
	title: String
	completed: Boolean
	dueDate: Time

	assignee: User @relationship(inverse:"tasks")
}
`
