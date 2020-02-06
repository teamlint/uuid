package main

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"encoding/base64"
	"encoding/hex"

	"github.com/chilts/sid"
	guuid "github.com/google/uuid"
	"github.com/kjk/betterguid"
	"github.com/lithammer/shortuuid"
	"github.com/oklog/ulid"
	"github.com/rs/xid"
	uuid "github.com/satori/go.uuid"
	"github.com/segmentio/ksuid"
	"github.com/sony/sonyflake"

	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
)

var engine *xorm.Engine
var conn = "root:123456@tcp(localhost:3306)/word?charset=utf8mb4&parseTime=True&loc=Local"

func base() {

	hello := base64.StdEncoding.EncodeToString([]byte("hello"))
	fmt.Printf("base64,hello:%v", hello)

}

func genShortUUID() {
	id := shortuuid.New()
	fmt.Printf("github.com/lithammer/shortuuid: %s\n", id)
}

func genUUID() {
	id := guuid.New()
	fmt.Printf("github.com/google/uuid:         %s\n", id.String())
}

func genXid() {
	id := xid.New()
	fmt.Printf("github.com/rs/xid:              %s\n", id.String())
	// fmt.Printf("    timestamp:%v,time:%v,raw:%v,pid:%v,counter:%v\r\n", id.Time().Unix(), id.Time(), hex.EncodeToString(id.Bytes()), id.Pid(), id.Counter())
	// Output: 9m4e2mr0ui3e8a215n4g
}

func genKsuid() {
	id := ksuid.New()
	fmt.Printf("github.com/segmentio/ksuid:     %s\n", id.String())
	// fmt.Printf("    timestamp:%v,raw:%s,payload:%v\r\n", id.Timestamp(), strings.ToUpper(hex.EncodeToString(id.Bytes())), strings.ToUpper(hex.EncodeToString(id.Payload())))
}

func genBetterGUID() {
	id := betterguid.New()
	fmt.Printf("github.com/kjk/betterguid:      %s\n", id)
}

func genUlid() {
	t := time.Now().UTC()
	entropy := rand.New(rand.NewSource(t.UnixNano()))
	id := ulid.MustNew(ulid.Timestamp(t), entropy)
	fmt.Printf("github.com/oklog/ulid:          %s\n", id.String())
}

func genSonyflake() {
	flake := sonyflake.NewSonyflake(sonyflake.Settings{})
	id, err := flake.NextID()
	if err != nil {
		log.Fatalf("flake.NextID() failed with %s\n", err)
	}
	// Note: this is base16, could shorten by encoding as base62 string
	fmt.Printf("github.com/sony/sonyflake:      %x\n", id)
}

func genSid() {
	id := sid.Id()
	fmt.Printf("github.com/chilts/sid:          %s\n", id)
}

func genUUIDv4() {
	// id, err := uuid.NewV4()
	// if err != nil {
	// 	log.Fatalf("uuid.NewV4() failed with %s\n", err)
	// }
	id := uuid.NewV4()
	fmt.Printf("github.com/satori/go.uuid:      %s\n", id)
}

func init() {
	var err error
	engine, err = xorm.NewEngine("mysql", conn)
	engine.ShowSQL(true)
	if err != nil {
		log.Fatal(err)
	}
}
func genData() {
	wg := sync.WaitGroup{}
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func(i int) {
			for j := 0; j < 5; j++ {
				execInsert()
			}
			wg.Done()
			log.Printf("insert rutine[%v] done.\n", i)
		}(i)
	}
	wg.Wait()
}
func genULID() ulid.ULID {
	t := time.Now().UTC()
	entropy := rand.New(rand.NewSource(t.UnixNano()))
	id := ulid.MustNew(ulid.Timestamp(t), entropy)
	return id
}
func execInsert() {
	sql := `INSERT INTO photos (ksuid, ksuid_raw, ksuid_time, xid, xid_raw, xid_time, ulid,ulid_raw,ulid_time,created_at)
VALUES
	(?,?,?, ?,?,?, ?,?,?,  ?);`
	kID := ksuid.New()
	xID := xid.New()
	uID := genULID()
	result, err := engine.Exec(sql,
		kID.String(), raw(kID.Bytes()), kID.Time(),
		xID.String(), raw(xID.Bytes()), xID.Time(),
		uID.String(), raw(uID[:]), ulid.Time(uID.Time()),
		time.Now())
	if err != nil {
		log.Fatalf("sql err: %v", err)
	}
	lastID, _ := result.LastInsertId()
	log.Printf("insert result: %v", lastID)
}
func raw(data []byte) string {
	return hex.EncodeToString(data)
}
func hex2dec(val string) uint64 {
	n, err := strconv.ParseUint(val, 16, 64)
	if err != nil {
		log.Fatal(err)
	}
	return n
}

func main() {
	genXid()
	genKsuid()
	genBetterGUID()
	genUlid()
	genSonyflake()
	genSid()
	genShortUUID()
	genUUIDv4()
	genUUID()
	// genData()
	execInsert()
}
