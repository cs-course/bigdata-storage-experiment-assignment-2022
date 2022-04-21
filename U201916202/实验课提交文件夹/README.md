# README

- my-s3bench

  该文件夹下存放了本次实验中我基于aws-sdk-go和go-echarts实现的测试程序，包括源代码和可执行程序。

  运行测试程序时需要传入一定的参数，以下给出一个例子。

  需要特别讲解的是**numOfClients, numOfSamples, objectSize可以接收数组**，程序会自动遍历数组元素的所有组合情况进行性能测试。注意，数组中的每个数字用英文逗号隔开，不能有其他字符。

  **requestMode参数用于设置测试时的请求模式，为0时表示常规请求，为1时表示对冲请求，为2时表示关联请求。**

  测试程序会将测试结果以文本报告的形式在终端输出，并生成两个可视化图表Write.html和Read.html，分别可视化展示了所有参数的写测试结果和读测试结果折线图。图中的图例"%d-%d-%d"代表这条折线所表示的测试参数{numOfClients, numOfSamples, objectSize}。例如，"8-256-1024"表示这条折线是参数为{numOfClients=8, numOfSamples=256, objectSize=1024}的测试的折线。**若读者觉得图表中折线太多影响观察，可以点击图例以隐藏或展示一条折线**。

  ```
  my-s3bench ^
      -accessKey=hust ^
      -secretKey=hust_obs ^
      -bucket=loadgen ^
      -endpoint=http://127.0.0.1:9090 ^
      -numOfClients=1,4,8 ^
      -numOfSamples=256,512,1024 ^
      -objectNamePrefix=loadgen ^
      -objectSize=1024,4096 ^
      -requestMode=0
  ```

- 常规请求测试结果

  作者在实验过程中某次“常规请求”测试(requestMode=0)的结果输出，保存在该文件夹中以供读者参考。

- 对冲请求测试结果

  作者在实验过程中某次“对冲请求”测试(requestMode=1)的结果输出，保存在该文件夹中以供读者参考。

- 关联请求测试结果

  作者在实验过程中某次“关联请求”测试(requestMode=2)的结果输出，保存在该文件夹中以供读者参考。

对本次实验有任何意见或建议者，欢迎联系作者，QQ573448239。