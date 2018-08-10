// Copyright 2018 Amazon.com, Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may not
// use this file except in compliance with the License. A copy of the
// License is located at
//
// http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
// either express or implied. See the License for the specific language governing
// permissions and limitations under the License.

// Package shell implements session shell plugin.
package shell

import (
	"os"
	"testing"
	"time"

	"io/ioutil"

	"github.com/aws/amazon-ssm-agent/agent/appconfig"
	"github.com/aws/amazon-ssm-agent/agent/context"
	"github.com/aws/amazon-ssm-agent/agent/contracts"
	"github.com/aws/amazon-ssm-agent/agent/framework/processor/executer/iohandler/mock"
	"github.com/aws/amazon-ssm-agent/agent/framework/runpluginutil"
	"github.com/aws/amazon-ssm-agent/agent/log"
	mgsContracts "github.com/aws/amazon-ssm-agent/agent/session/contracts"
	dataChannelMock "github.com/aws/amazon-ssm-agent/agent/session/datachannel/mocks"
	"github.com/aws/amazon-ssm-agent/agent/task"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/twinj/uuid"
)

var (
	payload       = []byte("testPayload")
	messageId     = "dd01e56b-ff48-483e-a508-b5f073f31b16"
	schemaVersion = uint32(1)
	createdDate   = uint64(1503434274948)
	mockLog       = log.NewMockLog()
)

type ShellTestSuite struct {
	suite.Suite
	mockContext     *context.Mock
	mockLog         log.T
	mockCancelFlag  *task.MockCancelFlag
	mockDataChannel *dataChannelMock.IDataChannel
	mockIohandler   *iohandlermocks.MockIOHandler
	stdin           *os.File
	stdout          *os.File
	plugin          runpluginutil.SessionPlugin
}

func (suite *ShellTestSuite) SetupTest() {
	mockContext := context.NewMockDefault()
	mockCancelFlag := &task.MockCancelFlag{}
	mockDataChannel := &dataChannelMock.IDataChannel{}

	suite.mockContext = mockContext
	suite.mockCancelFlag = mockCancelFlag
	suite.mockLog = mockLog
	suite.mockDataChannel = mockDataChannel
	suite.mockIohandler = new(iohandlermocks.MockIOHandler)
	stdout, stdin, _ := os.Pipe()
	suite.stdin = stdin
	suite.stdout = stdout
	suite.plugin = &ShellPlugin{
		stdin:  stdin,
		stdout: stdout,
	}
}

func (suite *ShellTestSuite) TearDownTest() {
	suite.stdin.Close()
	suite.stdout.Close()
}

// Testing Name
func (suite *ShellTestSuite) TestName() {
	rst := suite.plugin.Name()
	assert.Equal(suite.T(), rst, appconfig.PluginNameStandardStream)
}

// Testing GetOnMessageHandler
func (suite *ShellTestSuite) TestGetOnMessageHandler() {
	rst := suite.plugin.GetOnMessageHandler(suite.mockLog, suite.mockCancelFlag)

	assert.NotNil(suite.T(), rst)
}

// Testing Execute
func (suite *ShellTestSuite) TestExecuteWhenCancelFlagIsShutDown() {
	suite.mockCancelFlag.On("ShutDown").Return(true)
	suite.mockIohandler.On("MarkAsShutdown").Return(nil)

	suite.plugin.Execute(suite.mockContext,
		contracts.Configuration{},
		suite.mockCancelFlag,
		suite.mockIohandler,
		suite.mockDataChannel)

	suite.mockCancelFlag.AssertExpectations(suite.T())
	suite.mockIohandler.AssertExpectations(suite.T())
}

// Testing Execute
func (suite *ShellTestSuite) TestExecuteWhenCancelFlagIsCancelled() {
	suite.mockCancelFlag.On("Canceled").Return(true)
	suite.mockCancelFlag.On("ShutDown").Return(false)
	suite.mockIohandler.On("MarkAsCancelled").Return(nil)

	suite.plugin.Execute(suite.mockContext,
		contracts.Configuration{},
		suite.mockCancelFlag,
		suite.mockIohandler,
		suite.mockDataChannel)

	suite.mockCancelFlag.AssertExpectations(suite.T())
	suite.mockIohandler.AssertExpectations(suite.T())
}

// Testing Execute
func (suite *ShellTestSuite) TestExecute() {
	suite.mockCancelFlag.On("Canceled").Return(false)
	suite.mockCancelFlag.On("ShutDown").Return(false)
	suite.mockCancelFlag.On("Wait").Return(task.Completed)
	suite.mockIohandler.On("SetExitCode", 0).Return(nil)
	suite.mockIohandler.On("SetStatus", contracts.ResultStatusSuccess).Return()

	suite.plugin.Execute(suite.mockContext,
		contracts.Configuration{},
		suite.mockCancelFlag,
		suite.mockIohandler,
		suite.mockDataChannel)

	suite.mockCancelFlag.AssertExpectations(suite.T())
	suite.mockIohandler.AssertExpectations(suite.T())
}

// Testing writepump separately
func (suite *ShellTestSuite) TestWritePump() {
	stdout, stdin, _ := os.Pipe()
	stdin.Write(payload)

	//suite.mockDataChannel := &dataChannelMock.IDataChannel{}
	suite.mockDataChannel.On("SendStreamDataMessage", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	plugin := &ShellPlugin{
		stdout:      stdout,
		ipcFilePath: "test.log",
		dataChannel: suite.mockDataChannel,
	}

	// Spawning a separate go routine to close read and write pipes after a few seconds.
	// This is required as plugin.writePump() has a for loop which will continuosly read data from pipe until it is closed.
	go func() {
		time.Sleep(1800 * time.Millisecond)
		stdin.Close()
		stdout.Close()
	}()
	plugin.writePump(suite.mockLog)

	// Assert if SendStreamDataMessage function was called with same data from stdout
	suite.mockDataChannel.AssertExpectations(suite.T())
}

//Execute the test suite
func TestShellTestSuite(t *testing.T) {
	suite.Run(t, new(ShellTestSuite))
}

func TestProcessStreamMessage(t *testing.T) {
	stdinFile, _ := ioutil.TempFile("/tmp", "stdin")
	stdoutFile, _ := ioutil.TempFile("/tmp", "stdout")
	defer os.Remove(stdinFile.Name())
	defer os.Remove(stdoutFile.Name())
	plugin := &ShellPlugin{
		stdin:  stdinFile,
		stdout: stdoutFile,
	}
	agentMessage := getAgentMessage(uint32(mgsContracts.Output), payload)
	plugin.processStreamMessage(mockLog, *agentMessage)

	stdinFileContent, _ := ioutil.ReadFile(stdinFile.Name())
	assert.Equal(t, "testPayload", string(stdinFileContent))
}

// getAgentMessage constructs and returns AgentMessage with given sequenceNumber, messageType & payload
func getAgentMessage(payloadType uint32, payload []byte) *mgsContracts.AgentMessage {
	messageUUID, _ := uuid.Parse(messageId)
	agentMessage := mgsContracts.AgentMessage{
		MessageType:    mgsContracts.InputStreamDataMessage,
		SchemaVersion:  schemaVersion,
		CreatedDate:    createdDate,
		SequenceNumber: 1,
		Flags:          2,
		MessageId:      messageUUID,
		PayloadType:    payloadType,
		Payload:        payload,
	}
	return &agentMessage
}