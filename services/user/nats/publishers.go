package userNats

import (
	"microServiceBoilerplate/configs"
	natsPb "microServiceBoilerplate/proto/generated/nats"

	"github.com/mreza0100/golog"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

type publishers struct {
	lgr *golog.Core
	nc  *nats.Conn
}

func (p *publishers) DeleteUser(userId uint64) error {
	var (
		byteData []byte
		err      error
	)

	{
		byteData, err = proto.Marshal(&natsPb.UserDelete_EVENT{
			Id: userId,
		})
		if err != nil {
			p.lgr.Log("cant marshal pb (DeleteUser)")
			p.lgr.Log("error: ", err)
			return err
		}
	}

	{
		err = p.nc.Publish(configs.Nats.Subjects.DeleteUser, byteData)
		if err != nil {
			p.lgr.Log("cant publish pb (DeleteUser)")
			p.lgr.Log("error: ", err)
		}
	}

	return err
}
