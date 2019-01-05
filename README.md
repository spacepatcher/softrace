# Softrace 

Softrace is a simple Golang application for storing NIST National Software Reference Library Reference Data Set (NSRL RDS). Softrace based on Bolt database, so it should be fast.

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

### Run

Build docker image `insert_bolt` for creating Bolt database with NSRL data set:
```
docker build -t insert_bolt -f docker/Dockerfile.insert .
```

Create Bolt database file with `insert_bolt` container:
```
docker run -ti --name insert_bolt --mount src=`pwd`/data,target=/go/src/github.com/spacepatcher/softrace/data,type=bind insert_bolt
```

To start API service run:
```
docker-compose up
```

By default the API service is available on `localhost:8001`.


### Update data

