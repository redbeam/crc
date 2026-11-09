package fs9p

import (
	"errors"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWaitForErrorAfterStop(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)

	server, err := New9pServer(listener, t.TempDir())
	require.NoError(t, err)
	require.NoError(t, server.Start())
	require.NoError(t, server.Stop())

	assert.NoError(t, server.WaitForError())
}

func TestWaitForErrorReturnsServerError(t *testing.T) {
	acceptErr := errors.New("accept failed")
	listener := &failingListener{err: acceptErr, fail: make(chan struct{})}

	server, err := New9pServer(listener, t.TempDir())
	require.NoError(t, err)
	require.NoError(t, server.Start())
	close(listener.fail)

	assert.ErrorIs(t, server.WaitForError(), acceptErr)
}

// failingListener blocks in Accept() until fail is closed, then returns err
type failingListener struct {
	err  error
	fail chan struct{}
}

func (l *failingListener) Accept() (net.Conn, error) {
	<-l.fail
	return nil, l.err
}

func (l *failingListener) Close() error {
	return nil
}

func (l *failingListener) Addr() net.Addr {
	return &net.TCPAddr{}
}
