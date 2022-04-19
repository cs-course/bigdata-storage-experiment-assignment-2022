const fs = require('fs');

function createUniqueId(n) {
	var random = function () {
		// 生成10-12位不等的字符串
		return Number(Math.random().toString().substr(2)).toString(36); // 转换成十六进制
	};
	var arr = [];
	function createId() {
		var num = random();
		var _bool = false;
		arr.forEach((v) => {
			if (v === num) _bool = true;
		});
		if (_bool) {
			createId();
		} else {
			arr.push(num);
		}
	}
	var i = 0;
	while (i < n) {
		createId();
		i++;
	}
	return arr;
}

// 随机生成200000个不重复的数据并写入文件
fs.writeFileSync('data.json', JSON.stringify(createUniqueId(200000)))