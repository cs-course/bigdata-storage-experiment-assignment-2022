const BloomFilter = require("./BloomFilter");
const fs = require("fs");

// res 是一个20w条不重复的数据
const res = JSON.parse(fs.readFileSync("data.json", { encoding: "utf-8" }));

// 75%假设存入存储，25%数据不存入
const SPLITINDEX = 150000;
const data_stored = res.slice(0, SPLITINDEX);
const data_unstored = res.slice(SPLITINDEX);

// 初始化布隆过滤器
const bloomFilterInstance = new BloomFilter(SPLITINDEX, 0.01);
bloomFilterInstance.getFilterInfo();

let containMissJudged = 0;
let addSuccess = 0;
// 模拟存入数据
console.log("====================开始加入数据====================\n");
data_stored.forEach((d) => {
	try {
		if (bloomFilterInstance.add(d)) {
			addSuccess++;
		} else {
			containMissJudged++;
		}
	} catch (error) {
		console.log("加入数据库失败...\n");
	}
});
console.log("====================加入数据结束====================\n");
console.log(
	`在加入${
		data_stored.length
	}条数据中，成功加入${addSuccess}条，误判率的个数有${containMissJudged}, 误判率为${(
		containMissJudged / data_stored.length
	).toFixed(5)}`
);

containMissJudged = 0;

// 模拟查询数据
console.log("====================开始查询数据====================\n");
data_unstored.forEach((d) => {
	try {
		if (bloomFilterInstance.contain(d)) {
			containMissJudged++;
		}
	} catch (error) {
		console.log("查询失败...\n");
	}
});
console.log("====================查询数据结束====================\n");
console.log(
	`在查询${
		data_unstored.length
	}条数据中，误判率的个数有${containMissJudged}, 误判率为${(
		containMissJudged / data_unstored.length
	).toFixed(5)}`
);
