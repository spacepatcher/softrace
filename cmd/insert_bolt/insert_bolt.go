package main

import (
	"bufio"
	"encoding/csv"
	"github.com/boltdb/bolt"
	"github.com/golang/protobuf/proto"
	"github.com/spacepatcher/softrace"
	"io"
	"os"
	"sync"
)

func serialize(decoded *softrace.RDS) *[]byte {
	rds := &softrace.RDSProto{
		SHA1: decoded.SHA1,
		MD5:             decoded.MD5,
		CRC32:           decoded.CRC32,
		FileName:        decoded.FileName,
		FileSize:        decoded.FileSize,
		ProductName:     decoded.ProductName,
		ProductVersion:  decoded.ProductVersion,
		ApplicationType: decoded.ApplicationType,
		OpSystemName:    decoded.OpSystemName,
		OpSystemVersion: decoded.OpSystemVersion,
	}
	out, _ := proto.Marshal(rds)

	return &out
}

func newIndex(csvInput *softrace.RDS) (toIndex []softrace.Index, err error) {
	if !(len(csvInput.SHA1) == 40 && len(csvInput.MD5) == 32) {
		return []softrace.Index{}, nil
	}

	encoded := serialize(csvInput)
	toIndexSHA1 := softrace.Index{
		KeyType: []byte("sha1"),
		Key:     []byte(csvInput.SHA1),
		Encoded: *encoded,
	}
	toIndexMD5 := softrace.Index{
		KeyType: []byte("md5"),
		Key:     []byte(csvInput.MD5),
		Encoded: *encoded,
	}
	toIndex = append(toIndex, toIndexSHA1)
	toIndex = append(toIndex, toIndexMD5)

	return toIndex, nil
}

func readCSV(filesPath string, productsPath string, osPath string, tasks chan<- softrace.Index) {
	f, err := os.Open(filesPath)
	softrace.Check(err)
	defer f.Close()

	p, err := os.Open(productsPath)
	softrace.Check(err)
	defer p.Close()

	o, err := os.Open(osPath)
	softrace.Check(err)
	defer o.Close()

	//
	productsLines, err := csv.NewReader(p).ReadAll()
	if err != nil {
		panic(err)
	}

	products := make(map[string]softrace.RDSProducts)
	for _, line := range productsLines {
		products[line[0]] = softrace.RDSProducts{
			ProductCode: 	 line[0],
			ProductName: 	 line[1],
			ProductVersion:  line[2],
			OpSystemCode: 	 line[3],
			MfgCode: 		 line[4],
			Language: 		 line[5],
			ApplicationType: line[6],
		}
	}

	//
	opSystemsLines, err := csv.NewReader(o).ReadAll()
	if err != nil {
		panic(err)
	}

	opSystems := make(map[string]softrace.RDSOpSystems)
	for _, line := range opSystemsLines {
		opSystems[line[0]] = softrace.RDSOpSystems{
			OpSystemCode: 	 line[0],
			OpSystemName: 	 line[1],
			OpSystemVersion: line[2],
			MfgCode: 		 line[3],
		}
	}

	//
	reader := csv.NewReader(bufio.NewReader(f))
	for {
		csvLine, err := reader.Read()
		softrace.Check(err)
		if err == io.EOF {
			close(tasks)

			return
		}

		rdsData := softrace.RDS{
			SHA1:         	 csvLine[0],
			MD5:          	 csvLine[1],
			CRC32:        	 csvLine[2],
			FileName:     	 csvLine[3],
			FileSize:     	 csvLine[4],
			ProductName: 	 products[csvLine[4]].ProductName,
			ProductVersion:	 products[csvLine[4]].ProductVersion,
			ApplicationType: products[csvLine[4]].ApplicationType,
			OpSystemName:	 opSystems[csvLine[6]].OpSystemName,
			OpSystemVersion: opSystems[csvLine[6]].OpSystemVersion,
		}

		indexSlice, _ := newIndex(&rdsData)
		for _, index := range indexSlice {
			tasks <- index
		}
	}
}

func boltPut(db *bolt.DB, tasks <-chan softrace.Index, wg *sync.WaitGroup) {
	defer wg.Done()

	for forIndex := range tasks {
		db.Batch(func(tx *bolt.Tx) error {
			b, err := tx.CreateBucketIfNotExists(forIndex.KeyType)
			softrace.Check(err)

			if err := b.Put(forIndex.Key, forIndex.Encoded); err != nil {
				return err
			}

			return nil
		})
	}
}

func main() {
	forIndex := make(chan softrace.Index)
	wgIndexers := sync.WaitGroup{}

	db, _ := softrace.ConnectBolt()
	for w := 1; w <= 10000; w++ {
		wgIndexers.Add(1)
		go boltPut(db, forIndex, &wgIndexers)
	}
	readCSV(softrace.RDSFilesPath, softrace.RDSProductsPath, softrace.RDSOpSystemsPath, forIndex)

	wgIndexers.Wait()
	db.Close()
}
