var Minio = require('minio')

var minioClient = new Minio.Client({
    endPoint: '127.0.0.1',
    port: 9000,
    useSSL: false,
    accessKey: 'hust',
    secretKey: 'hust_obs'
});

var file = 'clientNum.json'

minioClient.makeBucket('europetrip', 'us-east-1', function(err) {
    if (err) return console.log(err)
    console.log('Bucket created successfully in "us-east-1".')
    var metaData = {
        'Content-Type': 'application/octet-stream',
        'X-Amz-Meta-Testing': 1234,
        'example': 5678
    }

    minioClient.fPutObject('europetrip', 'clientNumTest', file, metaData, function(err, etag) {
      if (err) return console.log(err)
      console.log('File uploaded successfully.')
    });
});