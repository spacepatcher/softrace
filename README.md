# Softrace 

Softrace is a simple Golang application for storing NIST National Software Reference Library Reference Data Set (<a href="https://www.nist.gov/software-quality-group/national-software-reference-library-nsrl" target="_blank">NSRL RDS</a>). Softrace is based on Bolt database, so it is fast and tiny.

The application is able to process md5 and sha1 hash lookup searches.

Hash lookup example:
```
curl -XGET http://localhost:8001/lookup/6d62c33de5e65af81ae52a529ed2f385
```

Softrace response:
```
{
    "sha1":"00012b50bc346b4b3cecf13cdbf734ad03258322",
    "md5":"6d62c33de5e65af81ae52a529ed2f385",
    "crc32":"a681c4f7",
    "file_name":"APPLYStrings.xml",
    "file_size":"6820",
    "product_name":"Microsoft Office 2000 Standard",
    "product_version":"2000",
    "application_type":"Operating System",
    "os_name":"Windows Media Center 2005",
    "os_version":"2005"
}
```

### Create database file

To create the required database file you first need to download archive `Modern RDS Minimal` from Current RDS Hash Sets page https://www.nist.gov/itl/ssd/software-quality-group/national-software-reference-library-nsrl/nsrl-download/current-rds. 

You should unpack the downloaded archive into `data/nsrl_rds/rds_modern/` with names of the extracted files:
```
NSRLFile.txt
NSRLMfg.txt
NSRLOS.txt
NSRLProd.txt
```

Build docker image `insert_bolt` for creating Bolt database with NSRL data set:
```
docker build -t insert_bolt -f docker/Dockerfile.insert .
```

To avoid the problem with access rights to Docker volume, set the correct access rights:
```
chown -R 1000:1000 data/
```

Create Bolt database file with `insert_bolt` container:
```
docker run -ti --name insert_bolt -v `pwd`/data:/go/src/github.com/spacepatcher/softrace/data:delegated insert_bolt
```

Creating a database file which is located in `data/bolt/bolt.db` takes several hours (in some cases, with low R/W it will take more than a day, so be patient). Total size of the database file by the end of the process is about 38 gigabytes.

### Start API service

To start API service run:
```
docker-compose up
```

By default the API service is available on `127.0.0.1:8001`.

### Update data

To update the database file data, you need to put new exemplars of files downloaded from Current RDS Hash Sets page in `data/nsrl_rds/rds_modern/`. After that, you should remove the existing database file and generate a new one.
