// Copyright (C) 2018 Storj Labs, Inc.
// See LICENSE for copying information.

package overlay

import (
	"context"
	"log"
	"crypto/rand"

	"github.com/gogo/protobuf/proto"
	"github.com/zeebo/errs"
	"go.uber.org/zap"

	"storj.io/storj/pkg/dht"
	"storj.io/storj/pkg/kademlia"
	"storj.io/storj/protos/overlay"
	"storj.io/storj/storage"
	"storj.io/storj/storage/boltdb"
	"storj.io/storj/storage/redis"
)

// ErrNodeNotFound error standardization
var ErrNodeNotFound = errs.New("Node not found")

// OverlayError creates class of errors for stack traces
var OverlayError = errs.Class("Overlay Error")

// Cache is used to store overlay data in Redis
type Cache struct {
	DB  storage.KeyValueStore
	DHT dht.DHT
}

// NewRedisOverlayCache returns a pointer to a new Cache instance with an initalized connection to Redis.
func NewRedisOverlayCache(address, password string, db int, DHT dht.DHT) (*Cache, error) {
	rc, err := redis.NewClient(address, password, db)
	if err != nil {
		return nil, err
	}

	return &Cache{
		DB:  rc,
		DHT: DHT,
	}, nil
}

// NewBoltOverlayCache returns a pointer to a new Cache instance with an initalized connection to a Bolt db.
func NewBoltOverlayCache(dbPath string, DHT dht.DHT) (*Cache, error) {
	bc, err := boltdb.NewClient(zap.L(), dbPath, boltdb.OverlayBucket)
	if err != nil {
		return nil, err
	}

	return &Cache{
		DB:  bc,
		DHT: DHT,
	}, nil
}

// Get looks up the provided nodeID from the redis cache
func (o *Cache) Get(ctx context.Context, key string) (*overlay.Node, error) {
	b, err := o.DB.Get([]byte(key))
	if err != nil {
		return nil, err
	}
	if b.IsZero() {
		// TODO: log? return an error?
		return nil, nil
	}

	na := &overlay.Node{}
	if err := proto.Unmarshal(b, na); err != nil {
		return nil, err
	}

	return na, nil
}

// Put adds a nodeID to the redis cache with a binary representation of proto defined Node
func (o *Cache) Put(nodeID string, value overlay.Node) error {
	data, err := proto.Marshal(&value)
	if err != nil {
		return err
	}

	return o.DB.Put(kademlia.StringToNodeID(nodeID).Bytes(), []byte(data))
}

// Bootstrap walks the initialized network and populates the cache
func (o *Cache) Bootstrap(ctx context.Context) error {
	nodes, err := o.DHT.GetNodes(ctx, "0", 1280)

	if err != nil {
		zap.Error(OverlayError.New("Error getting nodes from DHT", err))
	}

	for _, v := range nodes {
		found, err := o.DHT.FindNode(ctx, kademlia.StringToNodeID(v.Id))
		if err != nil {
			zap.Error(ErrNodeNotFound)
		}

		node, err := proto.Marshal(&found)
		if err != nil {
			return err
		}

		if err := o.DB.Put(kademlia.StringToNodeID(found.Id).Bytes(), node); err != nil {
			return err
		}
	}

	return err
}

// Refresh updates the cache db with the current DHT.
// We currently do not penalize nodes that are unresponsive,
// but should in the future.
func (o *Cache) Refresh(ctx context.Context) error {
	log.Print("starting cache refresh")
	r, err := randomID()
	if err != nil {
		return err
	}

	rid := kademlia.NodeID(r)
	near, err := o.DHT.GetNodes(ctx, rid.String(), 128)
	if err != nil {
		return err
	}

	for _, node := range near {
		pinged, err := o.DHT.Ping(ctx, *node)
		if err != nil {
			return err
		}
		err = o.DB.Put([]byte(pinged.Id), []byte(pinged.Address.Address))
		if err != nil {
			return err
		}
	}
	
	// TODO: Kademlia hooks to do this automatically rather than at interval
	nodes, err := o.DHT.GetNodes(ctx, "", 128)
	if err != nil {
		return err
	}

	for _, node := range nodes {
		pinged, err := o.DHT.Ping(ctx, *node)
		if err != nil {
			zap.Error(ErrNodeNotFound)
			return err
		} else {
			err := o.DB.Put([]byte(pinged.Id), []byte(pinged.Address.Address))
			if err != nil {
				return err
			}
		}
	}

	return err
}

// Walk iterates over each node in each bucket to traverse the network
func (o *Cache) Walk(ctx context.Context) error {
	// TODO: This should walk the cache, rather than be a duplicate of refresh
	return nil
}

func randomID() ([]byte, error) {
	result := make([]byte, 64)
	_, err := rand.Read(result)
	return result, err
}
