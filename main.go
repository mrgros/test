package main

import (
	"github.com/google/uuid"
	"github.com/tarantool/go-tarantool"
	"gopkg.in/vmihailenco/msgpack.v2"
	"log"
	"reflect"
)

func init() {
	msgpack.Register(reflect.TypeOf((*uuid.UUID)(nil)).Elem(),
		func(e *msgpack.Encoder, v reflect.Value) error {
			id := v.Interface().(uuid.UUID)
			bytes, err := id.MarshalBinary()
			if err != nil {
				return err
			}
			_, err = e.Writer().Write(bytes)
			if err != nil {
				return err
			}

			return nil
		},
		func(d *msgpack.Decoder, v reflect.Value) error {
			bytes := make([]byte, 16)
			_, err := d.Buffered().Read(bytes)
			if err != nil {
				return err
			}

			id, err := uuid.FromBytes(bytes)
			if err != nil {
				return err
			}

			v.Set(reflect.ValueOf(id))
			return nil
		},
	)
	msgpack.RegisterExt(2, (*uuid.UUID)(nil))
}

func main() {
	tntClient, err := tarantool.Connect("localhost:3301", tarantool.Opts{User: "tarantool", Pass: "tarantool"})
	if err != nil {
		log.Fatalf("can't connect to tarantool: %s", err)
	}

	// работает
	id1 := uuid.New()
	log.Printf("uuid as string is '%s'", id1.String())
	resp1, err := tntClient.Call17("api.test", []interface{}{id1})
	if err != nil {
		log.Fatalf("can't call api.test: %s", err)
	}
	log.Printf("id returned from tnt as string is '%s'", resp1.String())

	// не работает
	id2, _ := uuid.Parse("6dd87a3b-3d34-29bd-f80c-ab12ba8c696e")
	log.Printf("uuid as string is '%s'", id2.String())
	resp2, err := tntClient.Call17("api.test", []interface{}{id2})
	if err != nil {
		log.Fatalf("can't call api.test: %s", err)
	}
	log.Printf("id returned from tnt as string is '%s'", resp2.String())
}
