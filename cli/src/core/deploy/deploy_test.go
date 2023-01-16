package deploy

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"net"
	"strings"
	"testing"

	"github.com/docker/docker/api/types"
	_container "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/luno/jettison/jtest"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_attachUntilRemoved(t *testing.T) {
	mockApiClient := new(MockApiClient)

	type cases struct {
		expectedError string
		hookFunc      func()
	}

	testCases := []cases{
		// Case: receive error from cli.Container attach
		{
			expectedError: "test error",
			hookFunc: func() {
				mockApiClient.On("ContainerAttach").Return("test error").Once()
			},
		},
		// Case: receive success message to successChannel
		{
			hookFunc: func() {
				mockApiClient.On("ContainerAttach").Return(nil)

				mockApiClient.On("ContainerWait").Return(_container.ContainerWaitOKBody{
					StatusCode: 0,
					Error:      nil,
				}, nil).Once()
			},
		},
		// Case: receive expected "No such container" message
		{
			expectedError: "No such container",
			hookFunc: func() {
				mockApiClient.On("ContainerWait").Return(nil, "No such container").Once()
			},
		},
		// Case: receive error to errorChannel
		{
			expectedError: "test error",
			hookFunc: func() {
				mockApiClient.On("ContainerWait").Return(nil, "test error").Once()
			},
		},
	}

	for _, testCase := range testCases {
		testCase.hookFunc()

		err := attachUntilRemoved(mockApiClient, context.Background(), "")
		if err != nil {
			require.Equal(t, strings.Contains(err.Error(), testCase.expectedError), true)
		} else {
			jtest.RequireNil(t, err)
		}
	}
}

type MockApiClient struct {
	mock.Mock
	client.ContainerAPIClient
}

func (mock *MockApiClient) ContainerAttach(ctx context.Context, container string, options types.ContainerAttachOptions) (types.HijackedResponse, error) {
	args := mock.Called()

	response := types.HijackedResponse{
		Conn:   &net.IPConn{},
		Reader: bufio.NewReader(bytes.NewReader([]byte("test"))),
	}

	var err error
	if args.Get(0) != nil {
		err = errors.New(args.Get(0).(string))
	}

	return response, err
}

func (mock *MockApiClient) ContainerWait(ctx context.Context, container string, condition _container.WaitCondition) (<-chan _container.ContainerWaitOKBody, <-chan error) {
	args := mock.Called()

	newChan := make(chan _container.ContainerWaitOKBody)

	var containerWaitBody _container.ContainerWaitOKBody
	b := args.Get(0)
	if b != nil {
		containerWaitBody = args.Get(0).(_container.ContainerWaitOKBody)

		go func() {
			newChan <- containerWaitBody
		}()
	}

	errChan := make(chan error)

	var err error
	c := args.Get(1)
	if c != nil {
		err = errors.New(args.Get(1).(string))

		go func() {
			errChan <- err
		}()
	}

	return newChan, errChan
}
