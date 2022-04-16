const fs = require("fs");
const xlsx = require("node-xlsx");

const clientData = JSON.parse(fs.readFileSync("clientNum.json"));
const d = clientData.map((item) => {
	return [
		item.base.numClients,
		item.write.TotalThroughput,
		item.write.TotalDuration,
		item.read.TotalThroughput,
		item.read.TotalDuration,
	];
});
console.log(d);

const objectSizeData = JSON.parse(fs.readFileSync("objectSize.json"));
const d2 = objectSizeData.map((item) => {
	return [
		item.base.objectSize,
		item.write.TotalTransferred,
		item.write.TotalThroughput,
		item.write.TotalDuration,
		item.read.TotalTransferred,
		item.read.TotalThroughput,
		item.read.TotalDuration,
	];
});

const sheetData = [
	{
		name: "client_num",
		data: [
			[
				"client_num",
				"W_TotalThroughput",
				"W_TotalDuration",
				"R_TotalThroughput",
				"R_TotalDuration",
			],
			...d
		],
	},

	{
		name: "object_size",
		data: [
			[
				"object_size",
				"W_TotalTransferred",
				"W_TotalThroughput",
				"W_TotalDuration",
				"R_TotalTransferred",
				"R_TotalThroughput",
				"R_TotalDuration",
			],
			...d2
		],
	},
];

const buffer = xlsx.build(sheetData);

fs.writeFile('finalData.xlsx', buffer, function(err) {
  if (err) {
      console.log("Write failed: " + err);
      return;
  }

  console.log("Write completed.");
});
