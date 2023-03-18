package rabbitmq

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProducer_NewProducer(t *testing.T) {

}

func TestProducer_Close(t *testing.T) {
	t.Parallel()
	// Create the producer
	testProducer, err := NewProducer(rabbitmqURL,
		WithNewSignInSessionQueue("new_sign_in_session_queue"),
	)
	require.NoError(t, err)
	require.NotNil(t, testProducer)
	require.NoError(t, testProducer.Close())
	// Create a producer & close `msgChan` before `close`
	testProducer, err = NewProducer(rabbitmqURL,
		WithNewSignInSessionQueue("new_sign_in_session_queue"),
	)
	require.NoError(t, err)
	require.NotNil(t, testProducer)
	require.NoError(t, testProducer.msgChan.Close())
	require.Error(t, testProducer.Close())
	// Create a producer & close `connection` before close
	testProducer, err = NewProducer(rabbitmqURL,
		WithNewSignInSessionQueue("new_sign_in_session_queue"),
	)
	require.NoError(t, err)
	require.NotNil(t, testProducer)
	require.NoError(t, testProducer.q.Close())
	require.Error(t, testProducer.Close())

}
