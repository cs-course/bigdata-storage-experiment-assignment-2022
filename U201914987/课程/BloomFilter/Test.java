package BloomFilter;

import java.io.BufferedReader;
import java.io.FileReader;
import java.io.IOException;
import java.util.ArrayList;
import java.util.List;

public class Test {
    public static void main(String []args){
        BufferedReader br1=null;
        BufferedReader br2=null;
        List input=new ArrayList<String>();
        List test=new ArrayList<String>();
        try{
            br1=new BufferedReader(new FileReader("D:\\java\\BloomFilter\\src\\BloomFilter\\input.txt"));
            br2=new BufferedReader(new FileReader("D:\\java\\BloomFilter\\src\\BloomFilter\\test.txt"));
            String line=br1.readLine();
            while(line!=null){
                input.add(line);
                line= br1.readLine();
            }
            String cur= br2.readLine();
            while(cur!=null){
                test.add(cur);
                cur=br2.readLine();
            }
        }catch (IOException e) {
            e.printStackTrace();
        }
        finally {
            try {
                if (br1 != null) {
                    br1.close();
                }
                if (br2 != null) {
                    br2.close();
                }
            }catch (IOException f){
                f.printStackTrace();
            }
        }
        int n=input.size();
        MultiBloomFilter bloomFilter=new MultiBloomFilter();
        String[] temp=new String[2];
        for(int i=0;i<n;i+=2){
            temp[0]=input.get(i).toString();
            temp[1]=input.get(i+1).toString();
            bloomFilter.add(temp);
        }

        n=test.size();
        boolean judge;
        String ans;
        int wrong=0;
        float total=n/3;
        for(int i=0;i<n;i+=3){
            temp[0]=test.get(i).toString();
            temp[1]=test.get(i+1).toString();
            ans=test.get(i+2).toString();
            judge=bloomFilter.contains(temp);
            if(Boolean.parseBoolean(ans)==judge){
                System.out.printf("input: [[%-16s],[%-16s]]     expect: %-5s      ans: %-5s\n",temp[0],temp[1],ans,judge);
            }
            else{
                System.out.printf("input: [[%-16s],[%-16s]]     expect: %-5s      ans: %-5s\n",temp[0],temp[1],ans,judge);
                ++wrong;
            }
        }
        System.out.printf("error: %f\n",wrong/total);
        return;
    }
}
