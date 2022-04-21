package BloomFilter;

import java.util.BitSet;

public class SingleBloomFilter {

    //设置容量
    private static final int SIZE=32;
    //各个hash函数的应当有不同的seed，选取素数
    private static final int[] seeds={3,5,7,11,13};

    private BitSet bits=new BitSet(SIZE);
    private SingleSimpleHash[] func=new SingleSimpleHash[seeds.length];

    public SingleBloomFilter(){
        //建立各个hash函数
        for(int i=0;i<seeds.length;++i){
            func[i]=new SingleSimpleHash(SIZE,seeds[i]);
        }
    }

    //添加
    public void add(String val){
        for(SingleSimpleHash f:func){
            bits.set(f.hash(val),true);
        }
    }

    //判断是否存在
    public boolean contains(String val){
        boolean ans=true;
        for(SingleSimpleHash f:func){
            ans=ans&&bits.get(f.hash(val));
        }
        return ans;
    }

    public static void main(String []args){
        String test="this is a test";
        SingleBloomFilter bloomFilter=new SingleBloomFilter();
        System.out.println(bloomFilter.contains(test));
        bloomFilter.add(test);
        System.out.println(bloomFilter.contains(test));
        System.out.println(bloomFilter.contains("this is test2"));
        return;
    }
}

class SingleSimpleHash{

    //过滤器容量
    private int space;
    private int seed;

    public SingleSimpleHash(int space,int seed){
        this.space=space;
        this.seed=seed;
    }

    public int hash(String val){
        int result=0;
        int len=val.length();
        for(int i=0;i<len;++i){
            result=seed*result+val.charAt(i);
        }

        //与运算两位同时为1，结果才为1。确保计算出来的结果不会超过size
        return (space-1)&result;
    }
}
