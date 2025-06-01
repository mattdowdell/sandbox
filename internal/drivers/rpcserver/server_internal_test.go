package rpcserver

// URL returns the URL the server is listening upon.
//
// This is a test-only method for when port 0 is used to dynamically select a free port for use in
// unit tests.
func (s *Server) URL(path string) string {
	<-s.ch
	return "http://" + s.url + path
}
