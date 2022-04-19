import mmh3
import time

class CountingBloomFilter:
    def __init__(self, size, hash_num):
        self.size = size
        self.hash_num = hash_num
        self.byte_array = bytearray(size)

    def add(self, s):
        for seed in range(self.hash_num):
            result = mmh3.hash(s, seed) % self.size
            if self.byte_array[result] < 256:
                self.byte_array[result] += 1

    def lookup(self, s):
        for seed in range(self.hash_num):
            result = mmh3.hash(s, seed) % self.size
            if self.byte_array[result] == 0:
                return "Nope"
        return "Probably"

    def remove(self, s):
        for seed in range(self.hash_num):
            result = mmh3.hash(s, seed) % self.size
            if self.byte_array[result] > 0:
                self.byte_array[result] -= 1


for i in [1000,2000,3000,4000,5000]:
    count = 0
    cbf = CountingBloomFilter(i, 8)
    start1 = time.perf_counter()
    with open("1.txt", "r") as f1:
        for line1 in f1.readlines():
            line1 = line1.strip('\n')
            cbf.add(line1)
    end1 = time.perf_counter()
    time1 = end1 - start1

    start2 = time.perf_counter()
    with open("2.txt", "r") as f2:
        for line2 in f2.readlines():
            line2 = line2.strip('\n')
            res = cbf.lookup(line2)
            if (res == "Nope"):
                count += 1
    end2 = time.perf_counter()
    time2 = end2 - start2

    start3 = time.perf_counter()
    with open("1.txt", "r") as f1:
        for line3 in f1.readlines():
            line3 = line3.strip('\n')
            cbf.remove(line3)
    end3 = time.perf_counter()
    time3 = end3 - start3

    print("\nsize =", i,
          "\nAdd Time:", time1 * 100000,
          "\nSearch Time", time2 * 100000,
          "\nRemove Time", time3 * 100000,
          "\nFalse Rate:", (count - 1) / 50)
