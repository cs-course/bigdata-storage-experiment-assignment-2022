package BloomFilter;

import java.util.BitSet;

public class MultiBloomFilter {
    //设置容量
    private static final int SIZE=32;
    //各个hash函数的应当有不同的seed，选取素数2,3,5,7,11,13,17,19,23,29,31,37,41,43,47,53,59,61,67,71,73,79,83,89,91,101,103
    private static final int[][] seeds={{3,7,13},
            {5,11,11},{1,2,3}};

    private BitSet[] bits=new BitSet[seeds.length];
    private MultiSimpleHash[][] func=new MultiSimpleHash[seeds.length][seeds[0].length];

    public MultiBloomFilter(){
        for(int i=0;i<seeds.length;++i){
            bits[i]=new BitSet();
        }

        //建立各个hash函数
        for(int i=0;i<seeds.length;++i){
            for(int j=0;j<seeds[0].length;++j){
                func[i][j]=new MultiSimpleHash(SIZE,seeds[i][j]);
            }
        }
    }

    //添加
    public void add(String[] val){
        int n=val.length;
        for(int i=0;i<n;++i){
            //第i维数据hash存储到第i个BitSet
            for(MultiSimpleHash f:func[i]){
                bits[i].set(f.hash(val[i]),true);
            }
        }
    }

    //判断是否存在
    public boolean contains(String[] val){
        boolean ans=true;
        int n=val.length;
        for(int i=0;i<n;++i) {
            for (MultiSimpleHash f : func[i]) {
                ans = ans && bits[i].get(f.hash(val[i]));
            }
        }
        return ans;
    }

    public static void main(String []args){
        String[] test={"123","456"};
        MultiBloomFilter bloomFilter=new MultiBloomFilter();
        System.out.println(bloomFilter.contains(test));
        bloomFilter.add(test);
        System.out.println(bloomFilter.contains(test));
        String[] notIn={"abc","def"};
        System.out.println(bloomFilter.contains(notIn));
        return;
    }
}

class MultiSimpleHash{
    //过滤器容量
    private int space;
    private int seed;

    public MultiSimpleHash(int space,int seed){
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

docker build -t openstack-swift-docker .