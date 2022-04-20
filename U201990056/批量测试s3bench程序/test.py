import csv
from matplotlib import pyplot as plt


def enumerate_count(file_name):
    with open(file_name) as f0:
        for count, _ in enumerate(f0, 1):
            pass
    return count


path = 'D:/s3bench/minio.txt'
lines = enumerate_count(path)

f = open('D:/s3bench/data.csv', 'w', encoding='utf-8', newline="")

writer = csv.writer(f)

trans = []
through = []
durationL = []
err = []
ileL = []

trans1 = []
through1 = []
durationL1 = []
error1 = []
ileL1 = []

NumClient = [1, 1, 1, 1,  1, 1, 1, 1, 1, 1, 1, 1, 2, 4, 8, 16, 32, 64, 70, 80, 90, 100, 2, 4, 8, 16, 32, 64, 1, 1, 1, 1,
             1, 1, 1, 1, 1]

ObjectSize = [2048, 4096, 10240, 20480, 40960, 102400, 204800, 409600, 1048576,  4194304, 4096, 4096, 4096, 4096, 4096,
              4096, 4096, 4096, 4096, 4096, 4096,  102400, 102400, 102400, 102400, 102400, 102400,   4096,  4096, 4096,
              4096, 4096, 4096, 4096, 4096, 4096]
with open(path, 'rb') as f1:
    # length = len(line)
    # print(length)
    # print(lines)
    for i in range(0, 37):
        position = f1.readline()
        # print(content)
        # print(content1)
        # array = content.split()
        # trans-test = array
        # position = f1.read(40)
        # 1 represents the current pos
        transferred = f1.readline()
        transferred = (transferred.split())[2]
        print("The trans is:", transferred)
        transferred = float(transferred)
        trans.append(transferred)
        #
        throughput = f1.readline()
        throughput = (throughput.split())[2]
        print("The throughput is:", throughput)
        through.append(float(throughput))
        #
        duration = f1.readline()
        duration = (duration.split())[2]
        print("The duration is:", duration)
        durationL.append(float(duration))
        #
        error = f1.readline()
        error = (error.split())[3]
        print("The error is:", error)
        err.append(int(error))
        #
        ile = f1.readline()
        ile = (ile.split())[4]
        print("The ile is:", ile)
        ileL.append(float(ile))

        #
        position1 = f1.readline()
        # 1 represents the current pos
        transferred = f1.readline()
        transferred = (transferred.split())[2]
        print("The trans is:", transferred)
        transferred = float(transferred)
        trans1.append(transferred)
        #
        throughput = f1.readline()
        throughput = (throughput.split())[2]
        print("The throughput is:", throughput)
        through1.append(float(throughput))
        #
        duration = f1.readline()
        duration = (duration.split())[2]
        print("The duration is:", duration)
        durationL1.append(float(duration))
        #
        error = f1.readline()
        error = (error.split())[3]
        print("The error is:", error)
        error1.append(int(error))
        #
        ile = f1.readline()
        ile = (ile.split())[4]
        print("The ile is:", ile)
        ileL1.append(float(ile))
        #
        pos = f1.readline()
        # i += 13
        print(i)

f1.close()

print(through)
print(ObjectSize)


ObjectSize, through = (list(t) for t in zip(*sorted(zip(ObjectSize, through))))

print(through)
print(ObjectSize)

fig = plt.figure(dpi=128, figsize=(10, 6))
plt.plot(ObjectSize, through, c='red')
plt.title("speed-size", fontsize=24)
plt.xlabel("ObjectSize", fontsize=16)
fig.autofmt_xdate()
plt.ylabel("Total Throughput: (MB/s)", fontsize=16)
plt.tick_params(axis='both', which='major', labelsize=16)

plt.show()
