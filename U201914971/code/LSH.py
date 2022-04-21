import pandas as pd
from datasketch import MinHash, MinHashLSH


# 取出HeartDisease标签
def getFeatures(data):
    newData = data.loc[:, [
        'BMI', 'Smoking', 'AlcoholDrinking', 'Stroke', 'PhysicalHealth',
        'MentalHealth', 'DiffWalking', 'Sex', 'AgeCategory', 'Race',
        'Diabetic', 'PhysicalActivity', 'GenHealth', 'SleepTime', 'Asthma',
        'KidneyDisease', 'SkinCancer'
    ]]
    return newData


# 数据编码
def seriesEncode(s):
    newS = pd.Series([
        s[0], s[1].encode('utf8'), s[2].encode('utf8'), s[3].encode('utf8'),
        s[4], s[5], s[6].encode('utf8'), s[7].encode('utf8'),
        s[8].encode('utf8'), s[9].encode('utf8'), s[10].encode('utf8'),
        s[11].encode('utf8'), s[12].encode('utf8'), s[13],
        s[14].encode('utf8'), s[15].encode('utf8'), s[16].encode('utf8')
    ],
                     index=[
                         'BMI', 'Smoking', 'AlcoholDrinking', 'Stroke',
                         'PhysicalHealth', 'MentalHealth', 'DiffWalking',
                         'Sex', 'AgeCategory', 'Race', 'Diabetic',
                         'PhysicalActivity', 'GenHealth', 'SleepTime',
                         'Asthma', 'KidneyDisease', 'SkinCancer'
                     ])
    return newS


# 根据相似结果完成预测
def getFinalResult(dataSet, dataIndex):
    num_Yes = 0
    num_No = 0
    for i in dataIndex:
        if (dataSet.loc[i, 'HeartDisease'] == 'Yes'):
            num_Yes += 1
        else:
            num_No += 1
    if num_Yes >= num_No:
        return 'Yes'
    else:
        return 'No'


# 计算准确率
def accuracyRate(real_data, test_data):
    num_true = 0
    num_falsePositive = 0
    num_falseNegative = 0
    length = real_data.shape[0]
    for i in range(length):
        if (real_data.loc[i, 'HeartDisease'] == 'Yes'):
            if (test_data[i] == 'Yes'):
                num_true += 1
            else:
                num_falsePositive += 1
        else:
            if (test_data[i] == 'No'):
                num_true += 1
            else:
                num_falseNegative += 1
    result_list = [
        num_true / length, num_falsePositive / length,
        num_falseNegative / length
    ]
    return result_list


# 读取数据集
file = "../data/heart.csv"
original_data = pd.read_csv(file)
print("Size of original data:", original_data.shape[0], original_data.shape[1])

# 数据集划分训练集、验证集
train_data = original_data.sample(frac=0.95)
print("Size of train data:", train_data.shape[0], train_data.shape[1])
test_data = original_data[~original_data.index.isin(train_data.index)]
print("Size of test data:", test_data.shape[0], test_data.shape[1])

# 整理数据
train_data = train_data.reset_index(drop=True)
train_feature = getFeatures(train_data)
print("Size of new train data:", train_feature.shape[0],
      train_feature.shape[1])
test_data = test_data.reset_index(drop=True)
test_feature = getFeatures(test_data)
print("Size of new test data:", test_feature.shape[0], test_feature.shape[1])

# 创建lsh对象
lsh = MinHashLSH(threshold=0.8, num_perm=128)

# 初始化训练集，创建lsh索引
for i in range(train_feature.shape[0]):
    minhash = MinHash(num_perm=128)
    for element in seriesEncode(train_feature.loc[i]):
        minhash.update(element)
    lsh.insert(i, minhash)

# 初始化验证集
minhashes = {}
for i in range(test_feature.shape[0]):
    minhash = MinHash(num_perm=128)
    for element in seriesEncode(test_feature.loc[i]):
        minhash.update(element)
    minhashes[i] = minhash

# 预测结果
result_data = []
for i in range(len(minhashes)):
    result = lsh.query(minhashes[i])
    result_data.append(getFinalResult(train_data, result))
print("The result is:", result_data)

# 计算准确率
accuracy_rate = accuracyRate(test_data, result_data)
print("The true rate is:", accuracy_rate[0])
print("The falsePositive rate is:", accuracy_rate[1])
print("The falseNegative rate is:", accuracy_rate[2])
