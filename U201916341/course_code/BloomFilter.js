const rd = require("random-js");
const hashFn = require("murmurhash3js").x86.hash32;

/**
 * 布隆过滤器类
 * 构造函数接收两个参数
 *    n: the max number of elements
 *    maxError: the max error rate
 */
class BloomFilter {
	bitMap; // bit array. Type: Boolean[]
	bitSize; // the length of bitMap. Type: number
	hashNum; // the total number of hash functions, which depends on the number of elements and maxError. Type: number
	hashSeed;
	elCount;
	elMax;
	constructor(n, maxError) {
		this.elMax = n;
		this.bitSize = Math.ceil(
			n * (-Math.log(maxError) / (Math.log(2) * Math.log(2)))
		);
		this.hashNum = Math.ceil(Math.log(2) * (this.bitSize / n));
		this.elCount = 0;
		this.initBitMap();
		this.initHashSeeds();
	}

	initBitMap() {
		this.bitMap = new Array(this.bitSize).fill(false);
	}

	initHashSeeds() {
		const engine = rd.MersenneTwister19937.autoSeed();
		this.hashSeed = new Set();
		while (this.hashSeed.size < this.hashNum) {
			this.hashSeed.add(rd.integer(0, 9999999999999)(engine));
		}
	}

	getBit(index) {
		return this.bitMap[index];
	}

	setBit(index) {
		this.bitMap[index] = true;
	}

	add(el) {
		if (this.contain(el) || this.elCount === this.elMax) {
			return false;
		}

		for (const seed of this.hashSeed.values()) {
			let index = hashFn(el, seed) % this.bitSize;
			this.setBit(index);
		}
		this.elCount++;
		return true;
	}

	contain(el) {
		for (let seed of this.hashSeed.values()) {
			const index = hashFn(el, seed) % this.bitSize;
			if (!this.getBit(index)) {
				return false;
			}
		}
		return true;
	}

	clear() {
		this.bitMap.fill(false); // 重新置为false
	}

	getFilterInfo() {
		console.log(`bitMap位数:${this.bitSize},hashFn个数:${this.hashNum}`);
	}

	getVerificationValue(el){
		let sum = 0;
		for (let seed of this.hashSeed.values()) {
			const index = hashFn(el, seed) % this.bitSize;
			sum += index;
		}
		return sum;
	}
}

module.exports = BloomFilter;
