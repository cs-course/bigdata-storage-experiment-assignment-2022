// const Minio = require('minio');
// const minioClient = new Minio.Client({
//     endPoint: '127.0.0.1',
//     port: 9000,
//     useSSL: false,
//     accessKey: 'hust',
//     secretKey: 'hust_obs'
// });

const { execSync } = require("child_process");
const fs = require('fs');

const INDEXES = [
	"objectSize",
	"numClients",
	"numSamples",
	"Total Transferred",
	"Total Throughput",
	"Total Duration",
	"Number of Errors",
];

const BASE_CONFIG = `s3bench.exe     -accessKey=hust     -accessSecret=hust_obs     -bucket=loadgen     -endpoint=http://127.0.0.1:9000     -numSamples=1024     -objectNamePrefix=loadgen`;

let data = [];
let OBJECT_SIZE = 1024;
let CLIENT_NUM = 8;


while (OBJECT_SIZE <= 1024000) {
	const command = `${BASE_CONFIG}    -objectSize=${OBJECT_SIZE}    -numClients=${CLIENT_NUM}`;
	const res = execSync(command).toString();
	let validData = res
		.split("\n")
		.filter(Boolean)
		.filter((s) => {
			for (let i = 0; i < INDEXES.length; i++) {
				if (s.startsWith(INDEXES[i])) {
					return true;
				}
			}
			return false;
		})
		.map((s) =>
			s
				.split("")
				.filter((s) => s !== " ")
				.join("")
				.split(":")
		).slice(3);
  const base = validData.slice(0,3);
  const write = validData.slice(3,7);
  const read = validData.slice(7);
	const obj = {};
	[base, write, read].forEach((items, i) => {
    let key;
    switch(i){
      case 0: key = 'base';break;
      case 1: key = 'write'; break;
      case 2: key = 'read'; break;
    }
    obj[key] = {};
    items.forEach(index=> {
      const val = index[1].split("").filter((s) => s.match(/\d|\./)).join('');
      obj[key][index[0]] = val;
    })
		
	});

	data.push(obj)
	OBJECT_SIZE *= 2;
}

console.log(data);
fs.writeFileSync('objectSize.json', JSON.stringify(data));


OBJECT_SIZE = 1024;
CLIENT_NUM = 1;
data = [];

while (CLIENT_NUM <= 1024) {
	const command = `${BASE_CONFIG}    -objectSize=${OBJECT_SIZE}    -numClients=${CLIENT_NUM}`;
	const res = execSync(command).toString();
	let validData = res
		.split("\n")
		.filter(Boolean)
		.filter((s) => {
			for (let i = 0; i < INDEXES.length; i++) {
				if (s.startsWith(INDEXES[i])) {
					return true;
				}
			}
			return false;
		})
		.map((s) =>
			s
				.split("")
				.filter((s) => s !== " ")
				.join("")
				.split(":")
		).slice(3);
  const base = validData.slice(0,3);
  const write = validData.slice(3,7);
  const read = validData.slice(7);
	const obj = {};
	[base, write, read].forEach((items, i) => {
    let key;
    switch(i){
      case 0: key = 'base';break;
      case 1: key = 'write'; break;
      case 2: key = 'read'; break;
    }
    obj[key] = {};
    items.forEach(index=> {
      const val = index[1].split("").filter((s) => s.match(/\d|\./)).join('');
      obj[key][index[0]] = val;
    })
		
	});

	data.push(obj)
	CLIENT_NUM *= 2;
}

console.log(data);
fs.writeFileSync('clientNum.json', JSON.stringify(data));