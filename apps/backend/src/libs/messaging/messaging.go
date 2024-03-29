package messaging

import (
	"block-feed/src/libs/blockstore"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/redis/go-redis/v9"
)

// In this design, redis stream messages are structured as follows:
//
//	{ "Data": <JSON encoded data> }
//
// This makes XADD calls very simple:
//
//	XADD key "*" "Data" <JSON encoded data>
//
// As opposed to expanding the JSON out:
//
//	XADD key "*" <field1> <value1> <field2> <value2> ...
//
// This design also allows for storing nested data in the stream which
// is not as straightforward to do when using the expanded field-value
// form.
//
// One other important note is that this design makes it very easy
// to create a lua script which calls XADD in a loop to add elements
// in bulk. If we used the alternative approach of adding each JSON
// key-value pair individually in the XADD call, then we would need
// to send the key-value pairs as individual strings then perform
// extra string parsing in the Lua script itself since go-redis does
// not allow us to send slices/maps as arguments to Lua scripts.
// This extra string parsing will ultimately increase the script's
// execution time and block other Lua scripts from executing as a
// result. Parsing the input string is also very error prone (e.g.
// what if a key/value has a special character that conflicts with
// separator we're using). Sending the individual key-value pairs as
// strings also increases the chances that we'll run into memory
// errors since ARGV can only store a limited number of elements.
//
// One caveat with the JSON encoded approach is that we always need
// to ensure that when we add a new message to the stream, the JSON
// data is always added under the same field/key name, which in the
// example above is "Data". Instead of hardcoding this string everywhere
// we've defined it in a single place below. The StreamMessage[T] struct
// must have a field EXACTLY matching this name otherwise errors will
// occur when using JSON marshal/unmarshal on the messages we get back
// from the go-redis library.
var dataField string

type (
	StreamMessage[T any] struct {
		Data T `json:"Data"`
	}

	WebhookLoadBalancerStreamMsgData struct {
		WebhookID string
	}

	WebhookActivationStreamMsgData struct {
		WebhookID string
	}

	WebhookStreamMsgData struct {
		WebhookID   string
		BlockHeight uint64
		IsNew       bool
	}

	BlockStreamMsgData struct {
		ChainID string
		Blocks  []blockstore.BlockDocument
	}

	BlockFlushStreamMsgData struct {
		ChainID           string
		LatestBlockHeight uint64
		IsBlockStoreEmpty bool
	}
)

func init() {
	var empty StreamMessage[any]
	dataField = reflect.TypeOf(empty).Field(0).Tag.Get("json")
}

func (m StreamMessage[T]) MarshalBinary() ([]byte, error) {
	// https://github.com/redis/go-redis/issues/739#issuecomment-470634159
	return json.Marshal(m.Data)
}

func GetDataField() string {
	return dataField
}

func NewWebhookStreamMsg(blockHeight uint64, webhookID string, isNew bool) *StreamMessage[WebhookStreamMsgData] {
	return &StreamMessage[WebhookStreamMsgData]{
		Data: WebhookStreamMsgData{
			WebhookID:   webhookID,
			BlockHeight: blockHeight,
			IsNew:       isNew,
		},
	}
}

func NewWebhookLoadBalancerStreamMsg(webhookID string) *StreamMessage[WebhookLoadBalancerStreamMsgData] {
	return &StreamMessage[WebhookLoadBalancerStreamMsgData]{
		Data: WebhookLoadBalancerStreamMsgData{
			WebhookID: webhookID,
		},
	}
}

func NewWebhookActivationStreamMsg(webhookID string) *StreamMessage[WebhookActivationStreamMsgData] {
	return &StreamMessage[WebhookActivationStreamMsgData]{
		Data: WebhookActivationStreamMsgData{
			WebhookID: webhookID,
		},
	}
}

func NewBlockStreamMsg(chainID string, blocks []blockstore.BlockDocument) *StreamMessage[BlockStreamMsgData] {
	return &StreamMessage[BlockStreamMsgData]{
		Data: BlockStreamMsgData{
			ChainID: chainID,
			Blocks:  blocks,
		},
	}
}

func NewBlockFlushStreamMsg(chainID string, latestBlockHeight uint64, isBlockStoreEmpty bool) *StreamMessage[BlockFlushStreamMsgData] {
	return &StreamMessage[BlockFlushStreamMsgData]{
		Data: BlockFlushStreamMsgData{
			ChainID:           chainID,
			LatestBlockHeight: latestBlockHeight,
			IsBlockStoreEmpty: isBlockStoreEmpty,
		},
	}
}

func ParseMessage[T any](msg redis.XMessage) (*T, error) {
	values := msg.Values

	rawData, exists := values[dataField]
	if !exists {
		return nil, fmt.Errorf("key \"%s\" does not exist in map: %v", dataField, values)
	}

	var parsedMsg T
	if err := json.Unmarshal([]byte(fmt.Sprint(rawData)), &parsedMsg); err != nil {
		return nil, err
	}

	return &parsedMsg, nil
}
