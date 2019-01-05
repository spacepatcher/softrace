package main

import (
	"bytes"
	"encoding/json"
	"github.com/boltdb/bolt"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/proto"
	"github.com/spacepatcher/softrace"
	"net/http"
	"regexp"
)

var (
	regexSHA1 = regexp.MustCompile(`^[0-9A-F]{40}$`)
	regexMD5  = regexp.MustCompile(`^[0-9A-F]{32}$`)
)

func deserialize(encoded []byte) (data *softrace.RDS) {
	rds := &softrace.RDSProto{}
	_ = proto.Unmarshal(encoded, rds)
	data = &softrace.RDS{
		SHA1: 			 rds.SHA1,
		MD5:             rds.MD5,
		CRC32:           rds.CRC32,
		FileName:        rds.FileName,
		FileSize:        rds.FileSize,
		ProductName:     rds.ProductName,
		ProductVersion:  rds.ProductVersion,
		ApplicationType: rds.ApplicationType,
		OpSystemName:    rds.OpSystemName,
		OpSystemVersion: rds.OpSystemVersion,
	}

	return data
}

func newLookup(clientInput string) (validInput *softrace.LookupInput, valid bool) {
	input := []byte(clientInput)
	hash := bytes.ToUpper(input)
	if !(regexSHA1.Match(hash) || regexMD5.Match(hash)) {
		return &softrace.LookupInput{}, false
	}

	return &softrace.LookupInput{
		HashUpper: hash,
		Length:    len(hash),
	}, true
}

func main() {
	var res *softrace.RDS
	var bucket []byte

	db, _ := softrace.ConnectBolt()

	router := gin.Default()
	router.GET("/lookup/:hash", func(context *gin.Context) {
		lookup, valid := newLookup(context.Param("hash"))
		if valid {
			if lookup.Length == 32 {
				bucket = []byte("md5")
			} else {
				bucket = []byte("sha1")
			}

			err := db.View(func(tx *bolt.Tx) error {
				encoded := tx.Bucket(bucket).Get(lookup.HashUpper)
				res = deserialize(encoded)

				return nil
			})

			if err != nil {
				context.String(http.StatusInternalServerError, err.Error())
			} else {
				if res.SHA1 != "" {
					dat, _ := json.Marshal(res.MapToResponse(res))
					context.String(http.StatusOK, string(dat))
				} else {
					context.String(http.StatusNotFound, "Not found")
				}
			}
		} else {
			context.String(http.StatusBadRequest, "Bad request")
		}
	})
	router.Run(softrace.GinConn)
}
