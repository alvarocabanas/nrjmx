package gojmx

import (
	"context"
	"errors"
	"fmt"
	"github.com/apache/thrift/lib/go/thrift"

	"time"

	"github.com/newrelic/nrjmx/gojmx/nrprotocol"
)

const pingTimeout = 1000 * time.Millisecond

var (
	ErrAlreadyStarted = errors.New("nrjmx subprocess already started")
	ErrNotRunning     = errors.New("nrjmx subprocess is not running")
)

type JMXClient struct {
	jmxService nrprotocol.JMXService
	jmxProcess *jmxProcess
	ctx        context.Context
	socket     *thrift.TSocket
}

func NewJMXClient(ctx context.Context) *JMXClient {
	return &JMXClient{
		ctx: ctx,
	}
}

func (j *JMXClient) Init() (*JMXClient, error) {

	var err error
	j.jmxProcess, err = NewJMXProcess(j.ctx).Start()
	if err != nil {
		j.jmxProcess.stop() // TODO: Handle err
		return j, err
	}

	transport := thrift.NewStreamTransport(j.jmxProcess.Stdout, j.jmxProcess.Stdin)
	j.jmxService, err = j.configureJMXServiceClient(transport)

	if err != nil {
		j.jmxProcess.stop() // TODO: Handle err
		return j, err
	}

	err = j.ping(pingTimeout)
	if err != nil {
		j.jmxProcess.stop() // TODO: Handle err
		return j, err
	}

	return j, nil
}

func (j *JMXClient) configureJMXServiceClient(transport thrift.TTransport) (*nrprotocol.JMXServiceClient, error) {
	var protocolFactory thrift.TProtocolFactory
	protocolFactory = thrift.NewTJSONProtocolFactory()

	var transportFactory thrift.TTransportFactory
	transportFactory = thrift.NewTBufferedTransportFactory(8192)
	transportFactory = thrift.NewTFramedTransportFactory(transportFactory)

	transport, err := transportFactory.GetTransport(transport)
	if err != nil {
		return nil, err
	}

	iprot := protocolFactory.GetProtocol(transport)
	oprot := protocolFactory.GetProtocol(transport)
	jmxServiceClient := nrprotocol.NewJMXServiceClient(thrift.NewTStandardClient(iprot, oprot))
	return jmxServiceClient, err
}

//Connect(ctx context.Context, config *JMXConfig) (err error)
//Disconnect(ctx context.Context) (err error)
// Parameters:
//  - BeanName
//QueryMbean(ctx context.Context, beanName string) (r []*JMXAttribute, err error)
//GetLogs(ctx context.Context) (r []*LogMessage, err error)

func (j *JMXClient) ping(timeout time.Duration) error {
	ctx, cancel := context.WithCancel(j.ctx)
	defer cancel()
	done := make(chan struct{}, 1)
	go func() {
		for ctx.Err() == nil {
			err := j.jmxService.Ping(ctx)
			if err != nil {
				continue
			}
			done <- struct{}{}
			break
		}
	}()
	select {
	case <-time.After(timeout):
		return fmt.Errorf("ping timeout")
	case err := <-j.jmxProcess.errCh:
		return err
	case <-done:
		return nil
	}
}

func (j *JMXClient) checkState() error {
	if !j.jmxProcess.IsRunning() {
		return ErrNotRunning
	}

	err := j.jmxProcess.Error()
	if err != nil {
		return err
	}

	return nil
}

func (j *JMXClient) checkForTransportError(err error) error {
	if _, ok := err.(thrift.TTransportException); ok {
		return j.jmxProcess.WaitExitError(5 * time.Second)
	}
	return err
}

func (j *JMXClient) Connect(config *nrprotocol.JMXConfig, timeout int64) error {
	if err := j.checkState(); err != nil {
		return err
	}
	err := j.jmxService.Connect(j.ctx, config, timeout)
	return j.checkForTransportError(err)
}

func (j *JMXClient) GetMBeanNames(mbean string, timeout int64) ([]string, error) {
	if err := j.checkState(); err != nil {
		return nil, err
	}
	result, err := j.jmxService.GetMBeanNames(j.ctx, mbean, timeout)

	return result, j.checkForTransportError(err)
}

func (j *JMXClient) GetMBeanAttrNames(mbean string, timeout int64) ([]string, error) {
	if err := j.checkState(); err != nil {
		return nil, err
	}
	result, err := j.jmxService.GetMBeanAttrNames(j.ctx, mbean, timeout)
	return result, j.checkForTransportError(err)
}

func (j *JMXClient) GetMBeanAttr(mBeanName, mBeanAttrName string, timeout int64) (*nrprotocol.JMXAttribute, error) {
	if err := j.checkState(); err != nil {
		return nil, err
	}
	result, err := j.jmxService.GetMBeanAttr(j.ctx, mBeanName, mBeanAttrName, timeout)
	return result, j.checkForTransportError(err)
}

func (j *JMXClient) Disconnect() error {
	if err := j.checkState(); err != nil {
		return err
	}
	defer func() {
		//j.jmxProcess.stop()
		//j = NewJMXClient(j.ctx)
	}()
	err := j.jmxService.Disconnect(j.ctx)
	return j.checkForTransportError(err)
}

func (j *JMXClient) WriteJunk() {
	//_, err := j.socket.Write([]byte("some junk\n"))
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Fprintln(j.socket)
	//err = j.socket.Flush(context.Background())
	//if err != nil {
	//	panic(err)
	//}
	fmt.Fprintf(j.jmxProcess.Stdin, "aa")
	//fmt.Fprintf(j.jmxProcess.Stdin, "junk")
	//fmt.Fprintln(j.jmxProcess.Stdin)
	//transport := thrift.NewStreamTransport(j.jmxProcess.Stdout, j.jmxProcess.Stdin)
	//var err error
	//j.jmxService, err = j.configureJMXServiceClient(transport)
	//if err != nil {
	//	panic(err)
	//}
}

func (j *JMXClient) Close() error {
	j.jmxService.Disconnect(j.ctx)
	return j.jmxProcess.stop()
}
