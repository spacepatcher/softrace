package softrace

import (
	"github.com/boltdb/bolt"
	"io"
	"strings"
	"time"
)

// RDSFile CSV file path
const RDSFilesPath = "data/nsrl_rds/rds_modern/NSRLFile_1000000.txt"

// RDSProd CSV file path
const RDSProductsPath = "data/nsrl_rds/rds_modern/NSRLProd.txt"

// RDSProd CSV file path
const RDSOpSystemsPath = "data/nsrl_rds/rds_modern/NSRLOS.txt"

// BoltDB path
	const BoltPath = "data/bolt/bolt_1000000.db"

// Gin-gonic configuration
const GinConn = "0.0.0.0:8001"

//
type RDSFiles struct {
	SHA1         string
	MD5          string
	CRC32        string
	FileName     string
	FileSize     string
	ProductCode  string
	OpSystemCode string
	SpecialCode  string
}

//
type RDSProducts struct {
	ProductCode     string
	ProductName     string
	ProductVersion  string
	OpSystemCode    string
	MfgCode         string
	Language        string
	ApplicationType string
}

//
type RDSOpSystems struct {
	OpSystemCode    string
	OpSystemName    string
	OpSystemVersion string
	MfgCode         string
}

//
type RDS struct {
	SHA1            string
	MD5             string
	CRC32           string
	FileName        string
	FileSize        string
	ProductName     string
	ProductVersion  string
	ApplicationType string
	OpSystemName    string
	OpSystemVersion string
}

//
type Index struct {
	KeyType []byte
	Key     []byte
	Encoded []byte
}

//
type Result struct {
	SHA1            string `json:"sha1"`
	MD5             string `json:"md5"`
	CRC32           string `json:"crc32"`
	FileName        string `json:"file_name"`
	FileSize        string `json:"file_size"`
	ProductName     string `json:"product_name"`
	ProductVersion  string `json:"product_version"`
	ApplicationType string `json:"application_type"`
	OpSystemName    string `json:"os_name"`
	OpSystemVersion string `json:"os_version"`
}

//
type LookupInput struct {
	HashUpper []byte
	Length    int
}

//
func (*RDS) MapToResponse(n *RDS) Result {
	return Result{
		SHA1:            strings.ToLower(n.SHA1),
		MD5:             strings.ToLower(n.MD5),
		CRC32:           strings.ToLower(n.CRC32),
		FileName:        n.FileName,
		FileSize:        n.FileSize,
		ProductName:     n.ProductName,
		ProductVersion:  n.ProductVersion,
		ApplicationType: n.ApplicationType,
		OpSystemName:    n.OpSystemName,
		OpSystemVersion: n.OpSystemVersion,
	}
}

//
func ConnectBolt() (*bolt.DB, error) {
	db, err := bolt.Open(BoltPath, 0600, nil)
	Check(err)

	db.MaxBatchSize = 10000
	db.MaxBatchDelay = 10 * time.Millisecond

	return db, err
}

//
func Check(e error) {
	if (e != nil) && (e != io.EOF) {
		panic(e)
	}
}
