package relationNats

import (
	"microServiceBoilerplate/configs"
	natsPb "microServiceBoilerplate/proto/generated/nats"
	"microServiceBoilerplate/services/relation/instances"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/mreza0100/golog"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

type subscribers struct {
	srv instances.Sevice
	lgr *golog.Core
	nc  *nats.Conn
}

func (s *subscribers) GetFollowers() {
	subject := configs.Nats.Subjects.GetFollowers
	dbug, sussecc := s.lgr.DebugPKG("GetFollowers_REQUEST", false)

	{
		s.nc.Subscribe(subject, func(msg *nats.Msg) {
			var (
				userId    uint64
				response  []byte
				followers []uint64
				err       error
			)

			{
				request := &natsPb.GetFollowers_REQUESTRequest{}
				err = proto.Unmarshal(msg.Data, request)
				if dbug("cant Unmarshal request")(err) != nil {
					return
				}
				userId = request.GetUserId()
			}

			{
				followers = s.srv.GetFollowers(userId)
			}

			{
				response, err = proto.Marshal(&natsPb.GetFollowers_REQUESTResponse{
					Followers: followers,
				})
				if dbug("cant Marshal response")(err) != nil {
					return
				}
			}
			{
				sussecc(response)
				msg.Respond(response)
			}
		})
	}
}

func (s *subscribers) DeleteUser() {
	subject := configs.Nats.Subjects.DeleteUser
	debug, sussecc := s.lgr.DebugPKG("DeleteUser", false)

	{
		s.nc.Subscribe(subject, func(msg *nats.Msg) {
			var (
				userId uint64
				err    error
			)

			{
				request := &natsPb.UserDelete_EVENT{}
				err = proto.Unmarshal(msg.Data, request)
				if debug("proto.Unmarshal")(err) != nil {
					return
				}
				userId = request.GetId()
			}

			{
				err = s.srv.DeleteAllRelations(userId)
				if debug("s.srv.DeleteUser")(err) != nil {
					return
				}
				err = s.srv.DeleteLikes(userId)
				if debug("s.srv.DeleteLikes")(err) != nil {
					return
				}
				sussecc(userId)
			}
		})

	}
}

func (s *subscribers) IsFollowingGroup() {
	subject := configs.Nats.Subjects.IsFollowingGroup
	dbug, success := s.lgr.DebugPKG("IsFollowingGroup", false)

	{
		s.nc.Subscribe(subject, func(msg *nats.Msg) {
			var (
				request  *natsPb.IsFollowingGroup_REQUESTRequest
				response *natsPb.IsFollowingGroup_REQUESTResponse

				err error
			)

			{
				request = &natsPb.IsFollowingGroup_REQUESTRequest{}
				response = &natsPb.IsFollowingGroup_REQUESTResponse{}
			}
			{
				err = proto.Unmarshal(msg.Data, request)
				if dbug("proto.Unmarshal")(err) != nil {
					return
				}
			}
			{
				result, err := s.srv.IsFollowing(request.Follower, request.Followings)
				if dbug("s.srv.IsFollowing")(err) != nil {
					return
				}
				response.Result = result
			}
			{
				byteResponse, err := proto.Marshal(response)
				if dbug("proto.Marshal")(err) != nil {
					return
				}

				success(response.Result)
				msg.Respond(byteResponse)
			}
		})
	}
}

func (s *subscribers) isLikedGroup() {
	subject := configs.Nats.Subjects.IsLikedGroup
	dbug, success := s.lgr.DebugPKG("isLikedGroup", false)

	{
		s.nc.Subscribe(subject, func(msg *nats.Msg) {
			var (
				request         *natsPb.IsLikedGroup_REQUESTRequest
				convertedResult map[uint64]*empty.Empty
			)

			{
				request = &natsPb.IsLikedGroup_REQUESTRequest{}
				err := proto.Unmarshal(msg.Data, request)
				if dbug("proto.Unmarshal")(err) != nil {
					return
				}
			}
			{
				result, err := s.srv.IsLikedGroup(request.Liker, request.PostIds)
				if dbug("s.srv.IsLikedGroup")(err) != nil {
					return
				}

				convertedResult = make(map[uint64]*empty.Empty)
				for _, id := range result {
					convertedResult[id] = &empty.Empty{}
				}
			}
			{
				response, err := proto.Marshal(&natsPb.IsLikedGroup_REQUESTResponse{
					Result: convertedResult,
				})
				if dbug("proto.Marshal")(err) != nil {
					return
				}

				success(convertedResult)
				msg.Respond(response)
			}
		})
	}
}

func (s *subscribers) countLikes() {
	subject := configs.Nats.Subjects.CountLikes
	dbug, success := s.lgr.DebugPKG("countLikes", false)

	{
		s.nc.Subscribe(subject, func(msg *nats.Msg) {
			var (
				request         *natsPb.CountLikes_REQUESTRequest
				counts          instances.CountResult
				convertedCounts map[uint64]uint64

				err error
			)

			{
				request = &natsPb.CountLikes_REQUESTRequest{}
				err = proto.Unmarshal(msg.Data, request)
				if dbug("proto.Unmarshal")(err) != nil {
					return
				}
			}
			{
				counts, err = s.srv.CountLikes(request.PostId)
				if dbug("s.srv.CountLikes")(err) != nil {
					return
				}
			}
			{
				convertedCounts = make(map[uint64]uint64)
				for _, c := range counts {
					convertedCounts[c.Post_id] = uint64(c.Count)
				}
			}
			{
				byteResult, err := proto.Marshal(&natsPb.CountLikes_REQUESTResponse{
					Counts: convertedCounts,
				})
				if dbug("proto.Marshal")(err) != nil {
					return
				}
				success(counts)
				msg.Respond(byteResult)
			}
		})
	}
}
