package main

import (
	"encoding/xml"
	"slices"
	"sync"
)

// Data types mirror the client's XML structs exactly — the test server
// must serialize the same fields with the same element names.

type User struct {
	XMLName     xml.Name `xml:"User"`
	Id          string   `xml:"Id"`
	UserName    string   `xml:"UserName"`
	FirstName   string   `xml:"FirstName"`
	LastName    string   `xml:"LastName"`
	Active      bool     `xml:"Active"`
	Email       string   `xml:"Email"`
	AccessLevel string   `xml:"AccessLevel"`
	Brand       string   `xml:"Brand"`
}

type Team struct {
	XMLName               xml.Name `xml:"Team"`
	Id                    string   `xml:"Id"`
	Name                  string   `xml:"Name"`
	TeamCodeForBulkImport string   `xml:"TeamCodeForBulkImport"`
	ParentTeamId          string   `xml:"ParentTeamId"`
}

type Course struct {
	XMLName                   xml.Name `xml:"Course"`
	Id                        string   `xml:"Id"`
	Code                      string   `xml:"Code"`
	Name                      string   `xml:"Name"`
	Active                    bool     `xml:"Active"`
	ForSale                   bool     `xml:"ForSale"`
	OriginalId                string   `xml:"OriginalId"`
	Description               string   `xml:"Description"`
	EcommerceShortDescription string   `xml:"EcommerceShortDescription"`
	EcommerceLongDescription  string   `xml:"EcommerceLongDescription"`
	CourseCodeForBulkImport   string   `xml:"CourseCodeForBulkImport"`
	Price                     string   `xml:"Price"`
	AccessTillDate            string   `xml:"AccessTillDate"`
	AccessTillDays            string   `xml:"AccessTillDays"`
	CourseTeamLibrary         bool     `xml:"CourseTeamLibrary"`
	CreatedBy                 string   `xml:"CreatedBy"`
	SeqId                     string   `xml:"SeqId"`
}

type CourseUser struct {
	XMLName            xml.Name `xml:"User"`
	Id                 string   `xml:"Id"`
	UserName           string   `xml:"UserName"`
	FirstName          string   `xml:"FirstName"`
	LastName           string   `xml:"LastName"`
	Completed          bool     `xml:"Completed"`
	PercentageComplete float64  `xml:"PercentageComplete"`
	CompliantTill      string   `xml:"CompliantTill"`
	DueDate            string   `xml:"DueDate"`
	AccessTillDate     string   `xml:"AccessTillDate"`
}

type Module struct {
	XMLName     xml.Name `xml:"Module"`
	Id          string   `xml:"Id"`
	Code        string   `xml:"Code"`
	Name        string   `xml:"Name"`
	Description string   `xml:"Description"`
}

// CourseEnrollment tracks a user's enrollment status in a course.
type CourseEnrollment struct {
	UserID             string
	Completed          bool
	PercentageComplete float64
}

// State holds all in-memory data. Every read and write goes through a
// typed method that locks mu — CI tests hammer the server in parallel.
type State struct {
	mu sync.Mutex

	users    map[string]*User
	userList []*User

	teams     map[string]*Team
	teamList  []*Team
	teamUsers map[string][]string // teamId → []userId (ordered)

	courses    map[string]*Course
	courseList []*Course

	// modules per course, in insertion order for stable pagination.
	modules map[string][]*Module // courseId → []*Module

	// courseEnrollments per course, in insertion order for stable pagination.
	courseEnrollments map[string][]*CourseEnrollment // courseId → enrollments
}

func NewState() *State {
	s := &State{
		users:             make(map[string]*User),
		teams:             make(map[string]*Team),
		teamUsers:         make(map[string][]string),
		courses:           make(map[string]*Course),
		modules:           make(map[string][]*Module),
		courseEnrollments: make(map[string][]*CourseEnrollment),
	}
	seed(s)
	return s
}

// ListUsers returns a page of users starting at offset start.
func (s *State) ListUsers(start, limit int) []*User {
	s.mu.Lock()
	defer s.mu.Unlock()
	all := s.userList
	if start >= len(all) {
		return nil
	}
	end := start + limit
	if end > len(all) {
		end = len(all)
	}
	page := make([]*User, end-start)
	for i, u := range all[start:end] {
		cp := *u
		page[i] = &cp
	}
	return page
}

// ListTeams returns a page of teams starting at offset start.
func (s *State) ListTeams(start, limit int) []*Team {
	s.mu.Lock()
	defer s.mu.Unlock()
	all := s.teamList
	if start >= len(all) {
		return nil
	}
	end := start + limit
	if end > len(all) {
		end = len(all)
	}
	page := make([]*Team, end-start)
	for i, t := range all[start:end] {
		cp := *t
		page[i] = &cp
	}
	return page
}

// ListTeamUsers returns a page of users belonging to teamID.
func (s *State) ListTeamUsers(teamID string, start, limit int) ([]*User, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.teams[teamID]; !ok {
		return nil, false
	}
	memberIDs := s.teamUsers[teamID]
	if start >= len(memberIDs) {
		return nil, true
	}
	end := start + limit
	if end > len(memberIDs) {
		end = len(memberIDs)
	}
	page := make([]*User, 0, end-start)
	for _, uid := range memberIDs[start:end] {
		if u, ok := s.users[uid]; ok {
			cp := *u
			page = append(page, &cp)
		}
	}
	return page, true
}

// ListCourses returns a page of courses starting at offset start.
func (s *State) ListCourses(start, limit int) []*Course {
	s.mu.Lock()
	defer s.mu.Unlock()
	all := s.courseList
	if start >= len(all) {
		return nil
	}
	end := start + limit
	if end > len(all) {
		end = len(all)
	}
	page := make([]*Course, end-start)
	for i, c := range all[start:end] {
		cp := *c
		page[i] = &cp
	}
	return page
}

// GetCourse returns the course with the given ID, or false if not found.
func (s *State) GetCourse(id string) (*Course, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	c, ok := s.courses[id]
	if !ok {
		return nil, false
	}
	cp := *c
	return &cp, true
}

// ListCourseUsers returns a page of course users for courseID.
// Each CourseUser is built by joining enrollment data with the user record.
func (s *State) ListCourseUsers(courseID string, start, limit int) ([]*CourseUser, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.courses[courseID]; !ok {
		return nil, false
	}
	enrollments := s.courseEnrollments[courseID]
	if start >= len(enrollments) {
		return nil, true
	}
	end := start + limit
	if end > len(enrollments) {
		end = len(enrollments)
	}
	page := make([]*CourseUser, 0, end-start)
	for _, e := range enrollments[start:end] {
		u, ok := s.users[e.UserID]
		if !ok {
			continue
		}
		page = append(page, &CourseUser{
			Id:                 u.Id,
			UserName:           u.UserName,
			FirstName:          u.FirstName,
			LastName:           u.LastName,
			Completed:          e.Completed,
			PercentageComplete: e.PercentageComplete,
		})
	}
	return page, true
}

// ListModules returns a page of modules for courseID.
func (s *State) ListModules(courseID string, start, limit int) ([]*Module, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.courses[courseID]; !ok {
		return nil, false
	}
	all := s.modules[courseID]
	if start >= len(all) {
		return nil, true
	}
	end := start + limit
	if end > len(all) {
		end = len(all)
	}
	page := make([]*Module, end-start)
	for i, m := range all[start:end] {
		cp := *m
		page[i] = &cp
	}
	return page, true
}

// AssignCourseToUser adds userID to courseID's enrollments. Returns false if
// course or user does not exist. Duplicate enrollments are silently ignored.
func (s *State) AssignCourseToUser(userID, courseID string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.courses[courseID]; !ok {
		return false
	}
	if _, ok := s.users[userID]; !ok {
		return false
	}
	for _, e := range s.courseEnrollments[courseID] {
		if e.UserID == userID {
			return true // idempotent: already enrolled
		}
	}
	s.courseEnrollments[courseID] = append(s.courseEnrollments[courseID], &CourseEnrollment{
		UserID:             userID,
		Completed:          false,
		PercentageComplete: 0,
	})
	return true
}

// RemoveCourseFromUser removes userID from courseID's enrollments.
// Returns (found, courseExists).
func (s *State) RemoveCourseFromUser(userID, courseID string) (bool, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.courses[courseID]; !ok {
		return false, false
	}
	before := len(s.courseEnrollments[courseID])
	s.courseEnrollments[courseID] = slices.DeleteFunc(
		s.courseEnrollments[courseID],
		func(e *CourseEnrollment) bool { return e.UserID == userID },
	)
	return len(s.courseEnrollments[courseID]) < before, true
}
