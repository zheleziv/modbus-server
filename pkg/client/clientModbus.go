package client

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/things-go/go-modbus"
	myerr "zheleznovux.com/modbus-console/pkg"
	"zheleznovux.com/modbus-console/pkg/client/logger"
	tag "zheleznovux.com/modbus-console/pkg/client/tag"
)

type сlientModbus struct {
	name               string
	connectionType     string
	ip                 string
	port               int
	slaveId            uint8
	connectionAttempts int
	connectionTimeout  time.Duration
	log                logger.Logger
	state              string
	tags               []tag.TagInterface
	sender             modbus.Client
}

func (c *сlientModbus) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Name    string
		Ip      string
		SlaveId uint8
		State   string
		Tags    []tag.TagInterface
	}{
		Name:    c.name,
		Ip:      c.ip,
		SlaveId: c.slaveId,
		State:   c.state,
		Tags:    c.tags,
	})
}

// ======================инициализация========================{

// =================================Name
func (c *сlientModbus) Name() string {
	return c.name
}
func (c *сlientModbus) SetName(name string) error {
	if strings.TrimSpace(name) == "" {
		return myerr.New("invalid client name. {setter client name}")
	}
	c.name = name
	return nil
}

// =================================Type
func (c *сlientModbus) Type() string {
	return c.connectionType
}
func (c *сlientModbus) SetType() {
	c.connectionType = MODBUS_TCP
}

// ===================================IP
func (c *сlientModbus) Ip() string {
	return c.ip
}

// using net.parseIp
func (c *сlientModbus) SetIp(ip string) error {
	ipAddr := net.ParseIP(strings.TrimSpace(ip))
	if ipAddr == nil {
		return myerr.New("invalid client Ip. {setter client Ip}")
	} else {
		c.ip = ip
		return nil
	}
}

// ===================ConnectionAttempts
func (c *сlientModbus) ConnectionAttempts() int {
	return c.connectionAttempts
}
func (c *сlientModbus) SetConnectionAttempts(ca int) error {
	if ca <= 0 {
		return myerr.New("invalid client connection attempts. {setter client connection attempts}")
	}
	c.connectionAttempts = ca
	return nil
}

// =================================Port
func (c *сlientModbus) Port() int {
	return c.port
}
func (c *сlientModbus) SetPort(port int) error {
	if (port > 0xFFFF) || (port < 0) {
		c.port = 502
		return myerr.New("invalid client port")
	} else {
		c.port = port
		return nil
	}
}

// ==============================SlaveID
func (c *сlientModbus) SalveId() uint8 {
	return c.slaveId
}
func (c *сlientModbus) SetSalveId(sid uint8) error {
	if sid > 0xFF {
		return myerr.New("invalid client slaveID")
	}
	c.slaveId = sid
	return nil
}

// =================================Tags
func (c *сlientModbus) Tags() []tag.TagInterface {
	return c.tags
}

// func (c *сlientModbus) SetTags(tags []tag.TagInterface) error {
// 	for id := range tags {
// 		if _, err := c.TagByName(tags[id].Name()); err != nil {
// 			return err
// 		}
// 	}
// 	c.tags = tags
// 	return nil
// }

// ============================TagById
func (c *сlientModbus) TagById(id int) (tag.TagInterface, error) {
	if (id >= len(c.tags)) || (id < 0) {
		return nil, myerr.New("invalid id client tag")
	}
	return c.tags[id], nil
}
func (c *сlientModbus) TagByName(name string) (tag.TagInterface, error) {
	for id := range c.tags {
		if c.tags[id].Name() == name {
			return c.tags[id], nil
		}
	}
	return nil, myerr.New("invalid client tag name")
}
func (c *сlientModbus) SetTag(name string, address uint32, scanPeriod float64, dataType string) error {
	if _, err := c.TagByName(name); err != nil {
		adr16, err := checkModbusAddress(address, dataType)
		if err != nil {
			return myerr.New(err.Error())
		}
		t, err := tag.NewTag(name, adr16, scanPeriod, dataType)
		if err != nil {
			return myerr.New(err.Error())
		}
		c.tags = append(c.tags, t)
		return nil
	}
	return myerr.New("client tag name already exists")
}

func checkModbusAddress(address uint32, dataType string) (uint16, error) {
	if address >= tag.UINT16_MAX_VALUE {
		t := FUNCTION__TAG_TYPE[dataType]
		if t == nil {
			return 0, myerr.New("invalid tag data type")
		}

		tmpINT := int(address / 100000.0)
		isDefined := false
		for i := range t {
			if t[i] == tmpINT {
				isDefined = isDefined || true
			}
		}
		if !isDefined {
			return 0, myerr.New("invalid tag address")
		}

		tmpUINT16 := uint16(address - uint32(tmpINT*100000))
		if tmpUINT16 >= tag.UINT16_MAX_VALUE {
			return 0, myerr.New("invalid tag address")
		}
		return tmpUINT16 - 1, nil
	}
	return uint16(address - 1), nil
}

// ====================ConnectionTimeout
func (c *сlientModbus) ConnectionTimeout() time.Duration {
	return c.connectionTimeout
}
func (c *сlientModbus) SetConnectionTimeout(s float64) error {
	if s < 0 {
		return myerr.New("client connection timeout < 0")
	}
	c.connectionTimeout = time.Duration(s) * time.Second
	return nil
}

// // ==============================State
func (c *сlientModbus) setState(state string) (isChanged bool) {
	if c.state != state {
		c.state = state
		isChanged = true
	} else {
		isChanged = false
	}
	return
}

// ==========================Constructor
func NewClinetModbus(ip string, port int, slaveID uint8, name string, debug bool, ConnectionAttempts int, ConnectionTimeout float64) (*сlientModbus, error) {
	var c сlientModbus

	err := c.SetName(name)
	if err != nil {
		return nil, myerr.New(err.Error())
	}
	c.SetType()
	err = c.SetIp(ip)
	if err != nil {
		return nil, myerr.New(err.Error())
	}
	err = c.SetPort(port)
	if err != nil {
		return nil, myerr.New(err.Error())
	}
	err = c.SetSalveId(slaveID)
	if err != nil {
		return nil, myerr.New(err.Error())
	}
	err = c.SetConnectionAttempts(ConnectionAttempts)
	if err != nil {
		return nil, myerr.New(err.Error())
	}
	err = c.SetConnectionTimeout(ConnectionTimeout)
	if err != nil {
		return nil, myerr.New(err.Error())
	}
	c.log = logger.Logger{
		ParentNodeName: c.name,
		ParentNodeIp:   c.ip,
		IsLogOutput:    debug,
	}

	return &c, nil
}

//======================инициализация========================}

func (c *сlientModbus) Connect() error {
	provider := modbus.NewTCPClientProvider(
		c.ip + ":" + fmt.Sprint(c.port))
	c.sender = modbus.NewClient(provider)
	defer c.sender.Close()

	//устанавливаем соединение
	err := c.sender.Connect()
	if err != nil {
		return myerr.New("Ошибка соединения! " + err.Error())
	} else {
		return nil
	}
}

func (c *сlientModbus) Send(id int) error {
	if c.sender == nil {
		c.tags[id].SetState(false)
		return myerr.New("sender nil")
	}
	defer c.sender.Close()

	//новые типы должны быть указаны здесь
	switch c.tags[id].DataType() {
	case tag.COIL_TYPE:
		{
			resp, err := c.sender.ReadDiscreteInputs(c.slaveId, c.tags[id].Address(), 1)

			if err != nil {
				c.tags[id].SetState(false)
				return err
			}

			if len(resp) > 0 {
				c.tags[id].(*tag.CoilTag).SetValue(resp[0])
				// c.log.WriteWithTag(logger.INFO, c.state, c.tags[id].Name(), "Значение: "+strconv.Itoa(int(resp[0]))+".")
				return nil
			} else {
				c.tags[id].SetState(false)
				// c.log.WriteWithTag(logger.WARNING, c.state, c.tags[id].Name(), "Значение не было считано!")
				return nil
			}
		}
	case tag.WORD_TYPE:
		{
			resp, err := c.sender.ReadHoldingRegisters(c.slaveId, c.tags[id].Address(), 1)

			if err != nil {
				c.tags[id].SetState(false)
				return err
			}

			if len(resp) > 0 {
				c.tags[id].(*tag.WordTag).SetValue(resp[0])
				// c.log.WriteWithTag(logger.INFO, c.state, c.tags[id].Name(), "Значение: "+strconv.Itoa(int(resp[0]))+".")
				return nil
			} else {
				c.tags[id].SetState(false)
				// c.log.WriteWithTag(logger.WARNING, c.state, c.tags[id].Name(), "Значение не было считано!")
				return nil
			}
		}
	case tag.DWORD_TYPE:
		{
			resp, err := c.sender.ReadHoldingRegisters(c.slaveId, c.tags[id].Address(), 2)

			if err != nil {
				c.tags[id].SetState(false)
				return err
			}

			if len(resp) > 1 {
				var tmp uint32 = (uint32(resp[0]) << 16) + uint32(resp[1])
				c.tags[id].(*tag.DWordTag).SetValue(tmp)
				// c.log.WriteWithTag(logger.INFO, c.state, c.tags[id].Name(), "Значение: "+strconv.Itoa(int(tmp))+".")
				return nil
			} else {
				c.tags[id].SetState(false)
				// c.log.WriteWithTag(logger.WARNING, c.state, c.tags[id].Name(), "Значение не было считано!")
				return nil
			}
		}
	default:
		return myerr.New("resp nil")
	}
}

func (c *сlientModbus) Start(stop chan struct{}, wg *sync.WaitGroup) {
	connection := make(chan bool)
	quit := make(chan struct{})
	defer wg.Done()

	var wgi sync.WaitGroup
	wgi.Add(1)
	go c.TryConnect(stop, connection, &wgi)
	for {
		select {
		case <-stop: //канал сверху. Завершение сессии
			{
				close(quit)
				wgi.Wait()
				return
			}
		case cb := <-connection: // канал снизу. Плохое подключение => реконект
			{
				if cb {
					close(quit)
					wgi.Wait()
					quit = make(chan struct{})
					for tagId := range c.tags {
						wgi.Add(1)
						go c.startSender(tagId, quit, &wgi, connection)
					}
				} else {
					close(quit)
					wgi.Wait()
					quit = make(chan struct{})
					wgi.Add(1)
					go c.TryConnect(stop, connection, &wgi)
				}
			}
		}
	}
}

func (c *сlientModbus) startSender(tagId int, quit chan struct{}, wg *sync.WaitGroup, connect chan bool) {
	// c.log.WriteWithTag(logger.INFO, c.state, c.tags[tagId].Name(), "Запущен опрос тега!")

	ticker := time.NewTicker(time.Duration(c.tags[tagId].ScanPeriod()) * time.Second)

	defer wg.Done()

	for {
		select {
		case <-quit:
			{
				// c.log.WriteWithTag(logger.INFO, c.state, c.tags[tagId].Name(), "Завершен опрос тега!")
				return
			}
		case <-ticker.C:
			{
				err := c.Send(tagId)
				if err != nil {
					c.log.WriteWithTag(logger.ERROR, c.state, c.tags[tagId].Name(), err.Error())
					connect <- false
				}
			}
		}
	}
}

func (c *сlientModbus) TryConnect(quit chan struct{}, connection chan bool, wg *sync.WaitGroup) { /// connection day out
	defer wg.Done()
	ticker := time.NewTicker(time.Second)

	for {
		select {
		case <-quit:
			{
				return
			}
		case <-ticker.C:
			{
				ticker.Stop()
				for i := 1; i <= c.connectionAttempts; i++ {
					select {
					case <-quit:
						{
							return
						}
					default:
						{
							err := c.Connect()
							if err == nil {
								if isChanged := c.setState(GOOD); isChanged {
									c.log.Write(logger.INFO, c.state, "Подключенно!")
								}
								connection <- true
								return
							}
						}
					}
				}
				if isChanged := c.setState(BAD); isChanged {
					c.log.Write(logger.WARNING, c.state, "Не удалось подключиться!")
				}

				ticker.Reset(c.connectionTimeout)
			}
		}
	}
}
