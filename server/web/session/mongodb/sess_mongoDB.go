package mongodb

import (
	"bytes"
	"context"
	"encoding/gob"
	"net/http"
	"sync"
	"time"

	"github.com/beego/beego/v2/server/web/session"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	// Register the specific map type used by Beego sessions
	gob.Register(map[interface{}]interface{}{})
	// Register common types that might be stored inside the interface{} values
	gob.Register("")
	gob.Register(int(0))
	gob.Register(int32(0))
	gob.Register(int64(0))
	gob.Register(float32(0))
	gob.Register(float64(0))
	gob.Register(bool(false))
	gob.Register(time.Time{})

	// Register the provider with the name "mongodb"
	session.Register("mongodb", mongoProvider)
}

var mongoProvider = &MongoProvider{}

type MongoSessionStore struct {
	collection *mongo.Collection
	sid        string
	lock       sync.RWMutex
	values     map[interface{}]interface{}
}

// Set sets a value in the session store
func (st *MongoSessionStore) Set(ctx context.Context, key, value interface{}) error {
	st.lock.Lock()
	defer st.lock.Unlock()
	if st.values == nil {
		st.values = make(map[interface{}]interface{})
	}
	st.values[key] = value
	return nil
}

// Get retrieves a value from the session store
func (st *MongoSessionStore) Get(ctx context.Context, key interface{}) interface{} {
	st.lock.RLock()
	defer st.lock.RUnlock()
	if st.values == nil {
		return nil
	}
	return st.values[key]
}

// Delete removes a specific key from the session store
func (st *MongoSessionStore) Delete(ctx context.Context, key interface{}) error {
	st.lock.Lock()
	defer st.lock.Unlock()
	if st.values != nil {
		delete(st.values, key)
	}
	return nil
}

// Flush clears all values in the session store
func (st *MongoSessionStore) Flush(ctx context.Context) error {
	st.lock.Lock()
	defer st.lock.Unlock()
	st.values = make(map[interface{}]interface{})
	return nil
}

// SessionID returns the current session identifier
func (st *MongoSessionStore) SessionID(ctx context.Context) string {
	return st.sid
}

// SessionRelease saves the current session values back to MongoDB
func (st *MongoSessionStore) SessionRelease(ctx context.Context, w http.ResponseWriter) {
	st.lock.RLock()
	data, err := encodeMongoData(st.values)
	st.lock.RUnlock()

	if err != nil {
		return
	}

	filter := bson.M{"_id": st.sid}
	update := bson.M{
		"$set": bson.M{
			"data":        data,
			"last_access": time.Now(),
		},
	}
	opts := options.Update().SetUpsert(true)

	_, _ = st.collection.UpdateOne(ctx, filter, update, opts)
}

func (st *MongoSessionStore) SessionReleaseIfPresent(ctx context.Context, w http.ResponseWriter) {
	st.SessionRelease(ctx, w)
}

type MongoProvider struct {
	client     *mongo.Client
	collection *mongo.Collection
}

// SessionInit initializes the MongoDB connection and TTL indexes
func (mp *MongoProvider) SessionInit(ctx context.Context, gclifetime int64, config string) error {
	clientOptions := options.Client().ApplyURI(config)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}

	if err := client.Ping(ctx, nil); err != nil {
		return err
	}

	mp.client = client
	mp.collection = client.Database("beego_sessions").Collection("sessions")

	// Setup TTL index for auto-cleanup based on last_access
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "last_access", Value: 1}},
		Options: options.Index().SetExpireAfterSeconds(int32(gclifetime)),
	}
	_, err = mp.collection.Indexes().CreateOne(ctx, indexModel)
	return err
}

// SessionRead reads a session from MongoDB by its ID
func (mp *MongoProvider) SessionRead(ctx context.Context, sid string) (session.Store, error) {
	var doc struct {
		Data []byte `bson:"data"`
	}

	values := make(map[interface{}]interface{})

	err := mp.collection.FindOne(ctx, bson.M{"_id": sid}).Decode(&doc)
	if err == nil && len(doc.Data) > 0 {
		if decodedValues, decodeErr := decodeMongoData(doc.Data); decodeErr == nil && decodedValues != nil {
			values = decodedValues
		}
	}

	return &MongoSessionStore{
		collection: mp.collection,
		sid:        sid,
		values:     values,
	}, nil
}

// SessionExist checks if a session exists in the database
func (mp *MongoProvider) SessionExist(ctx context.Context, sid string) (bool, error) {
	count, err := mp.collection.CountDocuments(ctx, bson.M{"_id": sid})
	return count > 0, err
}

// SessionRegenerate migrates data from an old session ID to a new one
func (mp *MongoProvider) SessionRegenerate(ctx context.Context, oldsid, sid string) (session.Store, error) {
	var doc struct {
		Data []byte `bson:"data"`
	}
	err := mp.collection.FindOne(ctx, bson.M{"_id": oldsid}).Decode(&doc)
	if err == nil {
		filter := bson.M{"_id": sid}
		update := bson.M{"$set": bson.M{"data": doc.Data, "last_access": time.Now()}}
		_, _ = mp.collection.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
		_ = mp.SessionDestroy(ctx, oldsid)
	}
	return mp.SessionRead(ctx, sid)
}

// SessionDestroy removes a session from MongoDB
func (mp *MongoProvider) SessionDestroy(ctx context.Context, sid string) error {
	_, err := mp.collection.DeleteOne(ctx, bson.M{"_id": sid})
	return err
}

// SessionAll returns the total count of active sessions
func (mp *MongoProvider) SessionAll(ctx context.Context) int {
	count, _ := mp.collection.CountDocuments(ctx, bson.D{})
	return int(count)
}

// SessionGC is a no-op because MongoDB handles cleanup via TTL indexes
func (mp *MongoProvider) SessionGC(ctx context.Context) {
	/*
		Why is this empty?
		
		1. Native Cleanup: In the SessionInit function, a TTL (Time-To-Live) index 
		   was set on the "last_access" field. MongoDB automatically deletes session 
		   documents in the background once they exceed the configured 'gclifetime'.
		   
		2. Interface Requirement: Beego's session provider interface strictly 
		   requires a SessionGC method to exist so it can trigger garbage collection.
		
		Leaving this empty acts as a "no-op" (no operation). It safely fulfills Beego's 
		structural interface requirement without wasting resources duplicating the 
		cleanup work that MongoDB is already handling natively.
	*/
}

func encodeMongoData(data map[interface{}]interface{}) ([]byte, error) {
	if data == nil || len(data) == 0 {
		return []byte{}, nil
	}
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func decodeMongoData(b []byte) (map[interface{}]interface{}, error) {
	if len(b) == 0 {
		return make(map[interface{}]interface{}), nil
	}
	buf := bytes.NewBuffer(b)
	dec := gob.NewDecoder(buf)
	var data map[interface{}]interface{}
	err := dec.Decode(&data)
	if err != nil {
		return nil, err
	}
	return data, nil
}
