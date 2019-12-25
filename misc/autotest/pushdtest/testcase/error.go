package testcase

import (
	"bytes"
	"fmt"

	"github.com/meitu/bifrost/grpc/conn"
	pub "github.com/meitu/bifrost/grpc/publish"
	pb "github.com/meitu/bifrost/grpc/push"
)

var (
	ErrConnectRecords        = "expect the num of '%d' but actual the num of '%d' in connect records"
	ErrConnectSessionPresent = "expect  '%v' but actual '%v' in connect sessionpresent"
	ErrConnectMessageID      = "expect  '%v' but actual '%v' in connect MessageID"
	ErrPullMessages          = "expect the num of '%d' but actual the num of '%d' in pull message"
	ErrPullTopic             = "expect the topic '%s' but actual the topic '%s' int pull"
	ErrPullRetain            = "expect the retain '%v' but actual the retain '%v' int pull"
	ErrSubRetains            = "expect the num of '%d' but actual the num of '%d' in sub retain"
	ErrSubIndex              = "expect '%d' but '%d' in sub index"
	ErrSubPayload            = "expect '%v' but '%v' in sub payload"
	ErrNotifyTopic           = "expect the topic '%s' but actual the topic '%s' int notify"
	ErrNotifyIndex           = "expect '%d' but '%d' in notify index"
	ErrRangeMessage          = "expect the num of '%d' but actual the num of '%d' in range message"
	ErrRangeComplete         = "expect '%v' but actual  '%v' in range complete"
	ErrRangeOffset           = "expect '%v' but actual  '%v' in range offset"
	ErrRangeTopic            = "expect '%v' but actual  '%v' in range topic"
	ErrRecordLast            = "expect '%v' but actual  '%v' in record lastindex"
	ErrRecordCurrent         = "expect '%v' but actual  '%v' in record currentindex"
	ErrRecordTopic           = "expect '%v' but actual  '%v' in record topic"
	ErrPublishCode           = "expect '%v' but actual  '%v' in publish code"
	ErrPublishResult         = "expect '%v' but actual  '%v' in publish result"
	ErrPublishCount          = "expect '%v' but actual  '%v' in publish result count"
)

func EqualConnect(resp *pb.ConnectResp, count int, session bool, mid int64) error {
	if len(resp.Records) != count {
		return fmt.Errorf(ErrConnectRecords, count, len(resp.Records))
	}

	if resp.SessionPresent != session {
		return fmt.Errorf(ErrConnectSessionPresent, session, resp.SessionPresent)
	}
	//TODO MessageID
	if resp.MessageID != mid {
		return fmt.Errorf(ErrConnectMessageID, mid, resp.MessageID)
	}
	return nil
}

func EqualSubscribe(resp *pb.SubscribeResp, count int, payload []byte) error {
	if len(resp.RetainMessage) != count {
		return fmt.Errorf(ErrSubRetains, count, len(resp.RetainMessage))
	}
	for _, msg := range resp.RetainMessage {
		if bytes.Equal(msg.Payload, payload) {
			return fmt.Errorf(ErrSubPayload, payload, msg.Payload)
		}
	}
	return nil
}

func EqualPull(resp *pb.PullResp, count int, topic string, retain bool) error {
	if len(resp.Messages) != count {
		return fmt.Errorf(ErrPullMessages, count, len(resp.Messages))
	}
	for _, msg := range resp.Messages {
		//topic pull retain nil
		/*
			if msg.Topic != topic {
				return fmt.Errorf(ErrPullTopic, topic, msg.Topic)
			}
		*/
		if msg.Retain != retain {
			return fmt.Errorf(ErrPullRetain, retain, msg.Retain)
		}
	}
	return nil
}

func EqualNotify(req *conn.NotifyReq, topic string) error {
	if req.Topic != topic {
		return fmt.Errorf(ErrNotifyTopic, topic, req.Topic)
	}
	return nil
}

func EqualRange(req *pb.RangeUnackResp, count int, complete bool, offset []byte, topic string) error {
	//	if bytes.Compare(req.Offset, offset) != 0 {
	//		return fmt.Errorf(ErrRangeOffset, offset, req.Offset)
	//	}
	if req.Complete != complete {
		return fmt.Errorf(ErrRangeComplete, complete, req.Complete)
	}
	if len(req.Messages) != count {
		return fmt.Errorf(ErrRangeMessage, count, len(req.Messages))
	}
	for _, m := range req.Messages {
		if m.Topic != topic {
			return fmt.Errorf(ErrRangeTopic, topic, m.Topic)
		}
	}

	return nil
}

func EqualRecords(req *pb.Record, cindex, lindex []byte, topic string) error {
	if bytes.Compare(req.LastestIndex, lindex) != 1 {
		return fmt.Errorf(ErrRecordLast, lindex, req.LastestIndex)
	}
	if bytes.Equal(req.CurrentIndex, cindex) {
		return fmt.Errorf(ErrRecordCurrent, cindex, req.CurrentIndex)
	}
	if req.Topic != topic {
		return fmt.Errorf(ErrRecordTopic, topic, req.Topic)
	}
	return nil
}

func EqualPublish(resp *pub.PublishReply, code pub.ErrCode, count int, result pub.ErrCode) error {
	if resp.ReturnCode != code {
		return fmt.Errorf(ErrPublishCode, code, resp.ReturnCode)
	}
	if len(resp.Results) != count {
		return fmt.Errorf(ErrPublishCount, count, resp.Results)
	}
	for _, r := range resp.Results {
		if r != result {
			return fmt.Errorf(ErrPublishResult, result, r)
		}
	}
	return nil
}
