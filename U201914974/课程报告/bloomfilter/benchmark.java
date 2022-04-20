package bloomfilter;

import java.io.*;
import java.nio.charset.StandardCharsets;
import java.util.HashSet;
import java.util.Random;

public class benchmark {
    private static int m = 2000;    //位数组长度
    private static int k = 8;       //hash函数数量
    private static int loop = 6000; //查询次数

    public static void main(String[] args) {
        int[] seeds = new int[k];
        Random random = new Random();
        for(int i = 0; i < k; i++){
            seeds[i] = random.nextInt(100);
        }
        bloomfilter bf = new bloomfilter(m,seeds);
        HashSet<String> hashset = new HashSet<String>();
        BufferedReader br0,br1;
        try {
            br0 = new BufferedReader(new FileReader("student0.txt", StandardCharsets.UTF_8));
            br1 = new BufferedReader(new FileReader("student1.txt",StandardCharsets.UTF_8));
            String str;
            while((str = br0.readLine()) != null){
                bf.add(str);
                hashset.add(str);
            }
            int total = 0, wrong = 0;
            while((str = br1.readLine()) != null){
                if(!hashset.contains(str)){
                    total++;
                    if(bf.contains(str)) {
                        wrong++;
                    }
                }
            }
            System.out.println("bloom filter错误率：" + (double)wrong/total);
            long starttime = System.currentTimeMillis();
            for(int i = 0; i < loop; i++){
                br0 = new BufferedReader(new FileReader("student0.txt", StandardCharsets.UTF_8));
                while((str= br0.readLine()) != null){
                    bf.contains(str);
                }
            }
            long endtime = System.currentTimeMillis();
            System.out.println("bloom filter查询耗时：" + (endtime - starttime) + "ms");
            starttime = System.currentTimeMillis();
            for(int i = 0; i < loop; i++){
                br0 = new BufferedReader(new FileReader("student0.txt", StandardCharsets.UTF_8));
                while((str= br0.readLine()) != null){
                    hashset.contains(str);
                }
            }
            endtime = System.currentTimeMillis();
            System.out.println("hashset查询耗时：" + (endtime - starttime) + "ms");
        } catch (IOException e) {
            e.printStackTrace();
        }
    }
}
