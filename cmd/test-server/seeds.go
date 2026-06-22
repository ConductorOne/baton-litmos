package main

func seed(s *State) {
	// Courses first — enrollments reference them.
	courses := []*Course{
		{
			Id:          "course-intro",
			Code:        "INT-101",
			Name:        "Introduction to Security",
			Active:      true,
			Description: "Security fundamentals course",
			CreatedBy:   "admin@example.com",
			SeqId:       "1",
		},
		{
			Id:          "course-compliance",
			Code:        "COMP-201",
			Name:        "Compliance Training",
			Active:      true,
			Description: "Annual compliance certification",
			CreatedBy:   "admin@example.com",
			SeqId:       "2",
		},
		{
			Id:          "course-deprecated",
			Code:        "DEP-001",
			Name:        "Legacy Systems Overview",
			Active:      false, // inactive course — tests Active=false path
			Description: "No longer maintained",
			CreatedBy:   "admin@example.com",
			SeqId:       "3",
		},
	}
	for _, c := range courses {
		cp := *c
		s.courses[c.Id] = &cp
		s.courseList = append(s.courseList, &cp)
	}

	// Modules per course (exercised when enableModules=true).
	s.modules["course-intro"] = []*Module{
		{Id: "mod-intro-1", Code: "INT-101-A", Name: "What is Security?", Description: "Overview of security concepts"},
		{Id: "mod-intro-2", Code: "INT-101-B", Name: "Threat Landscape", Description: "Common threats and attack vectors"},
		{Id: "mod-intro-3", Code: "INT-101-C", Name: "Best Practices", Description: "Recommended security practices"},
	}
	s.modules["course-compliance"] = []*Module{
		{Id: "mod-comp-1", Code: "COMP-201-A", Name: "Policy Overview", Description: "Company policy summary"},
		{Id: "mod-comp-2", Code: "COMP-201-B", Name: "Data Handling", Description: "How to handle sensitive data"},
		{Id: "mod-comp-3", Code: "COMP-201-C", Name: "Reporting Incidents", Description: "Incident reporting procedures"},
	}
	s.modules["course-deprecated"] = []*Module{
		{Id: "mod-dep-1", Code: "DEP-001-A", Name: "Legacy Introduction", Description: "Introduction to legacy systems"},
	}

	// Teams.
	teams := []*Team{
		{Id: "team-engineering", Name: "Engineering", TeamCodeForBulkImport: "ENG"},
		{Id: "team-product", Name: "Product", TeamCodeForBulkImport: "PROD"},
		{Id: "team-operations", Name: "Operations", TeamCodeForBulkImport: "OPS"},
	}
	for _, t := range teams {
		cp := *t
		s.teams[t.Id] = &cp
		s.teamList = append(s.teamList, &cp)
	}

	// Users — diversity matters more than count.
	users := []*User{
		{
			Id: "user-1", UserName: "alice@example.com",
			FirstName: "Alice", LastName: "Admin",
			Active: true, Email: "alice@example.com",
			AccessLevel: "Admin", Brand: "default",
		},
		{
			Id: "user-2", UserName: "bob@example.com",
			FirstName: "Bob", LastName: "Builder",
			Active: true, Email: "bob@example.com",
			AccessLevel: "Learner", Brand: "default",
		},
		{
			Id: "user-3", UserName: "carol@example.com",
			FirstName: "Carol", LastName: "Collins",
			Active: false, Email: "carol@example.com", // disabled — tests STATUS_DISABLED
			AccessLevel: "Learner", Brand: "default",
		},
		{
			Id: "user-4", UserName: "dave@example.com",
			FirstName: "Dave", LastName: "Doe",
			Active: true, Email: "dave@example.com", // no teams, no courses — tests empty-grants path
			AccessLevel: "Learner", Brand: "default",
		},
		{
			Id: "user-5", UserName: "eve@example.com",
			FirstName: "Eve", LastName: "Evans",
			Active: true, Email: "eve@example.com",
			AccessLevel: "Learner", Brand: "default",
		},
	}
	for _, u := range users {
		cp := *u
		s.users[u.Id] = &cp
		s.userList = append(s.userList, &cp)
	}

	// Team memberships — overlapping to exercise double-counting checks.
	// user-2 is in two teams; user-3 (disabled) is in a team.
	s.teamUsers["team-engineering"] = []string{"user-1", "user-2"}
	s.teamUsers["team-product"] = []string{"user-2", "user-3"}
	s.teamUsers["team-operations"] = []string{"user-5"}

	// Course enrollments — cover completed, in-progress, and disabled-user paths.
	// course-intro: alice completed, bob in-progress, carol (disabled) in-progress
	s.courseEnrollments["course-intro"] = []*CourseEnrollment{
		{UserID: "user-1", Completed: true, PercentageComplete: 100},
		{UserID: "user-2", Completed: false, PercentageComplete: 60},
		{UserID: "user-3", Completed: false, PercentageComplete: 10}, // disabled user with enrollment
	}
	// course-compliance: alice and eve completed (tests same course, multiple completions)
	s.courseEnrollments["course-compliance"] = []*CourseEnrollment{
		{UserID: "user-1", Completed: true, PercentageComplete: 100},
		{UserID: "user-5", Completed: true, PercentageComplete: 100},
	}
	// course-deprecated: bob in-progress (inactive course, still has enrollments)
	s.courseEnrollments["course-deprecated"] = []*CourseEnrollment{
		{UserID: "user-2", Completed: false, PercentageComplete: 30},
	}
	// user-4 (dave) has no enrollments — exercises the empty-grants path
}
