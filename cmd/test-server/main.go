// Package main runs a local test server that replicates the Litmos REST API.
// Used by CI when no real-tenant credentials are available.
// Start with: go run ./cmd/test-server/
// Point the connector at it with LITMOS_BASE_URL=http://localhost:8080.
package main

import (
	"encoding/xml"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	testAPIKey = "test-api-key"
	serverPort = "8080"
)

// Response wrapper types — element names match what the Litmos API returns.
// The client's XML decoders expect these exact outer element names.

type usersResponse struct {
	XMLName xml.Name `xml:"Users"`
	Users   []*User  `xml:"User"`
}

type teamsResponse struct {
	XMLName xml.Name `xml:"Teams"`
	Teams   []*Team  `xml:"Team"`
}

type coursesResponse struct {
	XMLName xml.Name  `xml:"Courses"`
	Courses []*Course `xml:"Course"`
}

type courseUsersResponse struct {
	XMLName xml.Name      `xml:"Users"`
	Users   []*CourseUser `xml:"User"`
}

type modulesResponse struct {
	XMLName xml.Name  `xml:"Modules"`
	Modules []*Module `xml:"Module"`
}

// Server holds the shared state for all handlers.
type Server struct {
	state *State
}

func main() {
	state := NewState()
	s := &Server{state: state}

	mux := http.NewServeMux()

	// Doc URL: https://support.litmos.com/hc/en-us/articles/227738287-REST-API-Developer-Guide#users-list
	mux.HandleFunc("GET /v1.svc/users", s.requireAuth(s.handleListUsers))

	// Doc URL: https://support.litmos.com/hc/en-us/articles/227738287-REST-API-Developer-Guide#teams-list
	mux.HandleFunc("GET /v1.svc/teams", s.requireAuth(s.handleListTeams))

	// Doc URL: https://support.litmos.com/hc/en-us/articles/227738287-REST-API-Developer-Guide#teams-users
	mux.HandleFunc("GET /v1.svc/teams/{id}/users", s.requireAuth(s.handleListTeamUsers))

	// Doc URL: https://support.litmos.com/hc/en-us/articles/227738287-REST-API-Developer-Guide#courses-list
	mux.HandleFunc("GET /v1.svc/courses", s.requireAuth(s.handleListCourses))

	// Doc URL: https://support.litmos.com/hc/en-us/articles/227738287-REST-API-Developer-Guide#courses-get
	mux.HandleFunc("GET /v1.svc/courses/{id}", s.requireAuth(s.handleGetCourse))

	// Doc URL: https://support.litmos.com/hc/en-us/articles/227738287-REST-API-Developer-Guide#courses-users
	mux.HandleFunc("GET /v1.svc/courses/{id}/users", s.requireAuth(s.handleListCourseUsers))

	// Doc URL: https://support.litmos.com/hc/en-us/articles/227738287-REST-API-Developer-Guide#courses-modules
	mux.HandleFunc("GET /v1.svc/courses/{id}/modules", s.requireAuth(s.handleListModules))

	// Doc URL: https://support.litmos.com/hc/en-us/articles/227738287-REST-API-Developer-Guide#users-courses-assign
	mux.HandleFunc("POST /v1.svc/users/{id}/courses", s.requireAuth(s.handleAssignCourseToUser))

	// Doc URL: https://support.litmos.com/hc/en-us/articles/227738287-REST-API-Developer-Guide#users-courses-remove
	mux.HandleFunc("DELETE /v1.svc/users/{id}/courses/{courseId}", s.requireAuth(s.handleRemoveCourseFromUser))

	log.Printf("litmos test server listening on :%s (apikey=%s)", serverPort, testAPIKey)
	srv := &http.Server{
		Addr:              ":" + serverPort,
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

// requireAuth validates the apikey header on every request.
// The Litmos API uses a custom apikey header (not Authorization: Bearer).
func (s *Server) requireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("apikey") != testAPIKey {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}

// writeXML serializes v as XML with application/xml Content-Type.
func writeXML(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/xml")
	if err := xml.NewEncoder(w).Encode(v); err != nil {
		log.Printf("writeXML: %v", err)
	}
}

// atoiOr parses s as an integer, returning def on failure or empty string.
func atoiOr(s string, def int) int {
	if s == "" {
		return def
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return n
}

// pageParams returns start offset and limit from the request query.
// Litmos pagination uses ?start=N&limit=M (default limit 500).
func pageParams(r *http.Request) (int, int) {
	return atoiOr(r.URL.Query().Get("start"), 0),
		atoiOr(r.URL.Query().Get("limit"), 500)
}

// Doc URL: https://support.litmos.com/hc/en-us/articles/227738287-REST-API-Developer-Guide#users-list
func (s *Server) handleListUsers(w http.ResponseWriter, r *http.Request) {
	start, limit := pageParams(r)
	users := s.state.ListUsers(start, limit)
	writeXML(w, &usersResponse{Users: users})
}

// Doc URL: https://support.litmos.com/hc/en-us/articles/227738287-REST-API-Developer-Guide#teams-list
func (s *Server) handleListTeams(w http.ResponseWriter, r *http.Request) {
	start, limit := pageParams(r)
	teams := s.state.ListTeams(start, limit)
	writeXML(w, &teamsResponse{Teams: teams})
}

// Doc URL: https://support.litmos.com/hc/en-us/articles/227738287-REST-API-Developer-Guide#teams-users
func (s *Server) handleListTeamUsers(w http.ResponseWriter, r *http.Request) {
	teamID := r.PathValue("id")
	start, limit := pageParams(r)
	users, ok := s.state.ListTeamUsers(teamID, start, limit)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	writeXML(w, &usersResponse{Users: users})
}

// Doc URL: https://support.litmos.com/hc/en-us/articles/227738287-REST-API-Developer-Guide#courses-list
func (s *Server) handleListCourses(w http.ResponseWriter, r *http.Request) {
	start, limit := pageParams(r)
	courses := s.state.ListCourses(start, limit)
	writeXML(w, &coursesResponse{Courses: courses})
}

// Doc URL: https://support.litmos.com/hc/en-us/articles/227738287-REST-API-Developer-Guide#courses-get
func (s *Server) handleGetCourse(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	course, ok := s.state.GetCourse(id)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	writeXML(w, course)
}

// Doc URL: https://support.litmos.com/hc/en-us/articles/227738287-REST-API-Developer-Guide#courses-users
func (s *Server) handleListCourseUsers(w http.ResponseWriter, r *http.Request) {
	courseID := r.PathValue("id")
	start, limit := pageParams(r)
	users, ok := s.state.ListCourseUsers(courseID, start, limit)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	writeXML(w, &courseUsersResponse{Users: users})
}

// Doc URL: https://support.litmos.com/hc/en-us/articles/227738287-REST-API-Developer-Guide#courses-modules
func (s *Server) handleListModules(w http.ResponseWriter, r *http.Request) {
	courseID := r.PathValue("id")
	start, limit := pageParams(r)
	modules, ok := s.state.ListModules(courseID, start, limit)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	writeXML(w, &modulesResponse{Modules: modules})
}

// assignCoursesBody is the XML body sent by POST /v1.svc/users/{id}/courses.
// Doc URL: https://support.litmos.com/hc/en-us/articles/227738287-REST-API-Developer-Guide#users-courses-assign
type assignCoursesBody struct {
	XMLName xml.Name `xml:"Courses"`
	Courses []struct {
		Id string `xml:"Id"`
	} `xml:"Course"`
}

// Doc URL: https://support.litmos.com/hc/en-us/articles/227738287-REST-API-Developer-Guide#users-courses-assign
func (s *Server) handleAssignCourseToUser(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("id")

	var body assignCoursesBody
	if err := xml.NewDecoder(r.Body).Decode(&body); err != nil || len(body.Courses) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	for _, c := range body.Courses {
		if !s.state.AssignCourseToUser(userID, c.Id) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
	}

	// The real Litmos API returns an empty 200 with no body.
	w.WriteHeader(http.StatusOK)
}

// Doc URL: https://support.litmos.com/hc/en-us/articles/227738287-REST-API-Developer-Guide#users-courses-remove
func (s *Server) handleRemoveCourseFromUser(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("id")
	courseID := r.PathValue("courseId")

	enrolled, courseExists := s.state.RemoveCourseFromUser(userID, courseID)
	if !courseExists {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if !enrolled {
		// User was not enrolled — return 404 so the connector's not-found path fires.
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// The real Litmos API returns an empty 200 with no body.
	w.WriteHeader(http.StatusOK)
}
