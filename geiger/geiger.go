package geiger

import (
	"context"
	"fmt"
	"github.com/peter-mount/go-kernel/cron"
	"github.com/peter-mount/home-automation/util/graphite"
	"github.com/peter-mount/home-automation/util/mq"
	"go.bug.st/serial"
	"time"
)

type Geiger struct {
	cron       *cron.CronService `kernel:"inject"`
	mq         *mq.MQ            `kernel:"inject"`
	graphite   graphite.Graphite `kernel:"inject"`
	serialPort *string           `kernel:"config,geigerPort"`
	prefix     *string           `kernel:"config,geigerPrefix"`
	port       serial.Port
}

func (m *Geiger) Start() error {
	// This sets the baud rate of the connection.
	// This is the default for the GMC-300 Plus Geiger counter.
	// You might have to change this for other models.
	// See GQ-RFC1201.txt for details.
	mode := &serial.Mode{
		BaudRate: 115200,
		DataBits: 8,
		Parity:   serial.NoParity,
		StopBits: 0,
	}

	port, err := serial.Open(*m.serialPort, mode)
	if err != nil {
		return err
	}
	m.port = port

	// Poll the Geiger counter every 10 seconds
	_, err = m.cron.AddTask("0/10 * * * * *", m.getStats)
	if err != nil {
		return err
	}

	return m.getStats(context.Background())
}

// sendCmd sends a command to the Geiger Counter.
// cmd is the text content of the command without the prefix or suffix characters.
// len is the number of bytes expected from the command.
// len=0 if no response is expected.
func (m *Geiger) sendCmd(cmd string, len int) ([]byte, error) {
	c := fmt.Sprintf("<%s>>%c", cmd, 13)
	//log.Printf("sending %q", c)
	n, err := m.port.Write([]byte(c))
	if err != nil {
		return nil, err
	}
	//log.Printf("Sent %d", n)

	time.Sleep(time.Millisecond * 250)

	if len <= 0 {
		return nil, nil
	}

	//log.Printf("Reading %d", len)
	buf := make([]byte, len)

	n, err = m.port.Read(buf)
	if err != nil {
		return nil, err
	}
	//log.Printf("Read %d/%d", n, len)

	if n < len {
		return nil, fmt.Errorf("expected %d got %d (%v)", len, n, buf)
	}

	return buf, nil
}

// getStats is a Task that sends the relevant commands to the Geiger counter.
func (m *Geiger) getStats(ctx context.Context) error {
	//log.Println("Getting data")

	now := time.Now().UTC()

	b, err := m.sendCmd("GETCPM", 2)
	if err != nil {
		return err
	}
	err = m.graphite.Publish(now, *m.prefix+".cpm", toInt16(b))
	if err != nil {
		return err
	}

	b, err = m.sendCmd("GETTEMP", 4)
	if err != nil {
		return err
	}
	err = m.graphite.Publish(now, *m.prefix+".temperature", toFloat64(b))
	if err != nil {
		return err
	}

	b, err = m.sendCmd("GETVOLT", 1)
	if err != nil {
		return err
	}
	err = m.graphite.Publish(now, *m.prefix+".volt", float64(b[0])/10)
	if err != nil {
		return err
	}

	b, err = m.sendCmd("GETGYRO", 7)
	if err != nil {
		return err
	}
	err = m.graphite.Publish(now, *m.prefix+".gyro_x", toInt16(b[0:2]))
	if err != nil {
		return err
	}
	err = m.graphite.Publish(now, *m.prefix+".gyro_y", toInt16(b[2:4]))
	if err != nil {
		return err
	}
	err = m.graphite.Publish(now, *m.prefix+".gyro_z", toInt16(b[4:6]))
	if err != nil {
		return err
	}

	return nil
}

// toInt16 converts a 2 byte value to an unsigned integer
func toInt16(b []byte) int {
	return (int(b[0]) * 256) + int(b[1])
}

// toFloat64 converts a 4 byte value to a signed float64.
// Byte 0 contains the integer portion.
// Byte 1 contains the fraction
// Byte 2 if not 0 indicates the value is negative.
// Byte 3 is not referenced here but is not used in the protocol
func toFloat64(b []byte) float64 {
	temp := float64(b[0]) + (float64(b[1]) / 100)
	if b[2] != 0 {
		temp = -temp
	}
	return temp
}
