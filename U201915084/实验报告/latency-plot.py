# -*- coding: utf-8 -*-
import numpy as np
import pandas as pd
import matplotlib.pyplot as plt
from matplotlib.ticker import FuncFormatter
from matplotlib.pyplot import MultipleLocator, plot
rootpath='D:\\bigdata-storage-lab\\'
wait_time='0.612'
request_num='2'
filename='h_'+wait_time+'_'+request_num+'_put_4096_5_200.csv'
latency = pd.read_csv(rootpath+filename)['latency'].apply(pd.to_numeric).values
print(sum(latency)/len(latency))
plt.rcParams['figure.dpi'] = 300
plt.rcParams['figure.figsize'] = (5,5)
plt.suptitle(f'wait_time = {wait_time}, request_num = {request_num}',y=0.92)
plt.subplots_adjust(wspace=0.4)
plt.subplot(211)
plt.plot(latency)
plt.subplot(223)
plt.plot(sorted(latency, reverse=True))
# plt.show()
plt.subplot(224)

# 百分比换算
def to_percent(y, position):
    return str(100 * round(y, 2)) + "%"

# 设置纵轴为百分比
formatter = FuncFormatter(to_percent)
ax = plt.gca()
# ax.xaxis.set_major_locator(MultipleLocator(5))
ax.yaxis.set_major_formatter(formatter)
# 避免横轴数据起始位置与纵轴重合，调整合适座标范围
x_min = max(min(latency) * 0.8, min(latency) - 5)
x_max = max(latency)
plt.xlim(x_min, x_max)
# 绘制实际百分位延迟。bins即直方图条带数，bins越大，绘制的累积曲线越平滑
plt.hist(latency, cumulative=True, histtype='step', weights=[1./ len(latency)] * len(latency), bins=100000) 

# 排队论模型
# F(t)=1-e^(-1*a*t)
alpha = 0.3
X_qt = np.arange(min(latency), max(latency), 1.)
Y_qt = 1 - np.exp(alpha * (min(latency) - X_qt))
# 绘制排队论模型拟合
plt.plot(X_qt, Y_qt)

plt.grid()
plt.show()