package s3utils

import (
	"sync"

	"github.com/aws/aws-sdk-go/aws/session"
)

var s3Session *session.Session
var s3SessionMutex = &sync.Mutex{}

// UseSession sets the AWS session to be used by all items that need an AWS
// Session in this package. This way we don't have to keep passing a pointer
// to the session around everywhere.
func UseSession(s *session.Session) {
	s3SessionMutex.Lock()
	defer s3SessionMutex.Unlock()

	s3Session = s
}

// getSession - A package internal function for getting the session to use
func getSession() *session.Session {
	if s3Session == nil {
		s3SessionMutex.Lock()
		defer s3SessionMutex.Unlock()

		if s3Session == nil {
			var err error
			s3Session, err = session.NewSession()

			if err != nil {
				panic(err)
			}
		}
	}
	return s3Session
}
